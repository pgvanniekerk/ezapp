# EzApp

EzApp is a simple, opinionated Go framework for building applications with zero configuration overhead. It provides a structured way to quickly bootstrap applications that handle configuration loading, logging, concurrent service execution, graceful shutdown, and resource cleanup - all with sensible defaults and minimal boilerplate.

## Why EzApp?

EzApp eliminates the repetitive setup code found in most Go applications by providing:

- **Zero-config startup**: Automatic environment variable configuration loading
- **Built-in logging**: Structured logging with configurable levels
- **Lifecycle management**: Handles startup, execution, and shutdown phases
- **Concurrent execution**: Run multiple services safely with coordinated shutdown
- **Resource cleanup**: Systematic cleanup with timeout control
- **Error handling**: Centralized error handling with proper logging and exit codes

Perfect for microservices, CLI tools, web applications, and background workers that need reliable startup and shutdown behavior.

## Installation

```bash
go get github.com/pgvanniekerk/ezapp
```

## Quick Start

### 1. Define Your Configuration

Create a configuration struct using [Netflix go-env](https://github.com/Netflix/go-env) tags:

```go
package main

type Config struct {
    Port        int    `env:"PORT" default:"8080"`
    DatabaseURL string `env:"DATABASE_URL" required:"true"`
    LogLevel    string `env:"LOG_LEVEL" default:"INFO"`
    Workers     int    `env:"WORKER_COUNT" default:"5"`
}
```

### 2. Create Your Initializer Function

Create an initializer function (recommended in a separate file):

```go
// initializer.go
package main

import (
    "context"
    "database/sql"
    "net/http"
    
    "github.com/pgvanniekerk/ezapp"
    _ "github.com/lib/pq"
)

func Initialize(ctx ezapp.InitCtx[Config]) (ezapp.AppCtx, error) {
    // Access your configuration
    config := ctx.Config
    logger := ctx.Logger
    
    // Setup dependencies
    db, err := sql.Open("postgres", config.DatabaseURL)
    if err != nil {
        return ezapp.AppCtx{}, err
    }
    
    // Create your services as runner functions
    server := createHTTPServer(config.Port, db, logger)
    worker := createBackgroundWorker(db, logger)
    
    // Define cleanup function for resource cleanup
    cleanup := func(shutdownCtx context.Context) error {
        logger.Info("Cleaning up resources...")
        return db.Close()
    }
    
    // Construct and return the application context
    return ezapp.Construct(
        ezapp.WithRunners(server, worker),
        ezapp.WithCleanup(cleanup),
    )
}

// Runner functions must match: func(context.Context) error
func createHTTPServer(port int, db *sql.DB, logger *zap.Logger) func(context.Context) error {
    return func(ctx context.Context) error {
        server := &http.Server{
            Addr: fmt.Sprintf(":%d", port),
            // ... configure your server
        }
        
        // Start graceful shutdown listener
        go func() {
            <-ctx.Done()
            logger.Info("Shutting down HTTP server...")
            server.Shutdown(context.Background())
        }()
        
        logger.Info("Starting HTTP server", zap.Int("port", port))
        if err := server.ListenAndServe(); err != http.ErrServerClosed {
            return err
        }
        return nil
    }
}

func createBackgroundWorker(db *sql.DB, logger *zap.Logger) func(context.Context) error {
    return func(ctx context.Context) error {
        logger.Info("Starting background worker...")
        
        for {
            select {
            case <-ctx.Done():
                logger.Info("Background worker shutting down...")
                return nil
            default:
                // Do your work here
                time.Sleep(time.Second)
            }
        }
    }
}
```

### 3. Wire Everything Together

Your main function becomes incredibly simple:

```go
// main.go
package main

import "github.com/pgvanniekerk/ezapp"

func main() {
    ezapp.Run(Initialize)
    // That's it! EzApp handles everything else
}
```

## Application Lifecycle

EzApp manages the complete application lifecycle in this order:

### 1. **Configuration Loading**
- Loads your configuration struct from environment variables
- Validates required fields and applies defaults
- Fails fast with clear error messages if configuration is invalid

### 2. **Logger Initialization**
- Creates a structured zap.Logger with configurable log level
- Controlled by `EZAPP_LOG_LEVEL` environment variable

### 3. **Startup Context Creation**
- Creates a context with configurable startup timeout
- Controlled by `EZAPP_STARTUP_TIMEOUT` (default: 15 seconds)
- Contains shutdown timeout information for later use

### 4. **Application Initialization**
- Calls your initializer function with populated `InitCtx`
- Provides access to config, logger, and startup context
- Your initializer wires dependencies and creates runners

### 5. **Concurrent Execution**
- Runs all runners concurrently in separate goroutines
- Monitors for SIGINT/SIGTERM signals for graceful shutdown
- Uses error groups for coordinated error handling

### 6. **Graceful Shutdown**
- Cancels context to signal all runners to stop
- Waits for all runners to complete gracefully

### 7. **Resource Cleanup**
- Calls cleanup function (if provided) with shutdown timeout
- Controlled by `EZAPP_SHUTDOWN_TIMEOUT` (default: 15 seconds)
- Ensures resources are properly released

### 8. **Exit**
- Logs completion status and exits
- Uses appropriate exit codes for different scenarios

## Environment Variables

### EzApp Framework Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `EZAPP_LOG_LEVEL` | `INFO` | Log level: `DEBUG`, `INFO`, `WARN`, `ERROR`, `DPANIC`, `PANIC`, `FATAL` |
| `EZAPP_STARTUP_TIMEOUT` | `15` | Startup timeout in seconds |
| `EZAPP_SHUTDOWN_TIMEOUT` | `15` | Cleanup timeout in seconds |

### Your Application Variables

Define your own configuration using go-env tags:

```go
type Config struct {
    // Required field
    DatabaseURL string `env:"DATABASE_URL" required:"true"`
    
    // Optional with default
    Port int `env:"PORT" default:"8080"`
    
    // String with default
    Environment string `env:"ENV" default:"development"`
    
    // Boolean values
    EnableMetrics bool `env:"ENABLE_METRICS" default:"true"`
}
```

## Advanced Usage

### Multiple Services

```go
func Initialize(ctx ezapp.InitCtx[Config]) (ezapp.AppCtx, error) {
    return ezapp.Construct(
        ezapp.WithRunners(
            createHTTPServer(ctx.Config, ctx.Logger),
            createGRPCServer(ctx.Config, ctx.Logger),
            createMetricsServer(ctx.Config, ctx.Logger),
            createBackgroundWorker(ctx.Config, ctx.Logger),
        ),
        ezapp.WithCleanup(cleanupResources),
    )
}
```

### Complex Cleanup

```go
func Initialize(ctx ezapp.InitCtx[Config]) (ezapp.AppCtx, error) {
    db := setupDatabase(ctx.Config.DatabaseURL)
    cache := setupRedis(ctx.Config.RedisURL)
    
    cleanup := func(shutdownCtx context.Context) error {
        // Cleanup in reverse order of initialization
        cacheErr := cache.Close()
        dbErr := db.Close()
        
        // Join all errors and return the combined error
        return errors.Join(cacheErr, dbErr)
    }
    
    return ezapp.Construct(
        ezapp.WithRunners(createServer(db, cache)),
        ezapp.WithCleanup(cleanup),
    )
}
```

### Accessing Shutdown Timeout in Runners

```go
func createServer(config Config, logger *zap.Logger) func(context.Context) error {
    return func(ctx context.Context) error {
        server := &http.Server{Addr: fmt.Sprintf(":%d", config.Port)}
        
        go func() {
            <-ctx.Done()
            
            // Get shutdown timeout from the startup context
            // (passed through the context chain)
            shutdownTimeout := config.GetShutdownTimeout(ctx)
            shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
            defer cancel()
            
            server.Shutdown(shutdownCtx)
        }()
        
        return server.ListenAndServe()
    }
}
```

## Error Handling

EzApp uses `logger.Fatal()` for all error conditions, which logs the error and exits with an appropriate code:

- **Configuration errors**: Invalid environment variables or struct validation
- **Initialization errors**: Failures in your initializer function
- **Runtime errors**: Failures from any runner function
- **Cleanup errors**: Failures in cleanup function (after successful run)

All errors are logged with context before termination.

## Best Practices

1. **Keep initializer separate**: Put your initializer function in a separate file (e.g., `initializer.go`)

2. **Design for graceful shutdown**: All runners should respect context cancellation

3. **Use structured logging**: Leverage the provided zap.Logger for consistent logging

4. **Handle cleanup properly**: Release resources in reverse order of acquisition

5. **Fail fast**: Validate dependencies early in the initializer

6. **Configure via environment**: Use environment variables for all configuration

## License

This project is licensed under the MIT License - see the LICENSE file for details.