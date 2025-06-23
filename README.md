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

        // Use a WaitGroup to wait for all cleanup operations to complete
        var wg sync.WaitGroup

        // Create a channel to collect errors
        errCh := make(chan error, 3) // Buffer for all possible errors

        // Close resources concurrently
        wg.Add(3)

        // Close HTTP server
        go func() {
            defer wg.Done()
            if err := server.Close(shutdownCtx); err != nil {
                errCh <- fmt.Errorf("failed to close HTTP server: %w", err)
            }
        }()

        // Close background worker
        go func() {
            defer wg.Done()
            if err := worker.Close(shutdownCtx); err != nil {
                errCh <- fmt.Errorf("failed to close background worker: %w", err)
            }
        }()

        // Close database
        go func() {
            defer wg.Done()
            if err := db.Close(); err != nil {
                errCh <- fmt.Errorf("failed to close database: %w", err)
            }
        }()

        // Wait for all goroutines to complete
        wg.Wait()
        close(errCh)

        // Collect all errors
        var errs []error
        for err := range errCh {
            errs = append(errs, err)
        }

        // Join all errors and return the combined error
        if len(errs) > 0 {
            return errors.Join(errs...)
        }
        return nil
    }

    // Construct and return the application context
    return ezapp.Construct(
        ezapp.WithRunners(server.Run, worker.Run),
        ezapp.WithCleanup(cleanup),
    )
}

// Runners are typically structs with Run and Close methods
type HTTPServer struct {
    server *http.Server
    logger *zap.Logger
}

func createHTTPServer(port int, db *sql.DB, logger *zap.Logger) *HTTPServer {
    return &HTTPServer{
        server: &http.Server{
            Addr: fmt.Sprintf(":%d", port),
            // ... configure your server
        },
        logger: logger,
    }
}

// Run method implements the Runner interface
func (s *HTTPServer) Run(ctx context.Context) error {
    s.logger.Info("Starting HTTP server", zap.String("addr", s.server.Addr))
    if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
        return err
    }
    return nil
}

// Close method for cleanup
func (s *HTTPServer) Close(ctx context.Context) error {
    s.logger.Info("Shutting down HTTP server...")
    return s.server.Shutdown(ctx)
}

type BackgroundWorker struct {
    db     *sql.DB
    logger *zap.Logger
    stopCh chan struct{}
}

func createBackgroundWorker(db *sql.DB, logger *zap.Logger) *BackgroundWorker {
    return &BackgroundWorker{
        db:     db,
        logger: logger,
        stopCh: make(chan struct{}),
    }
}

// Run method implements the Runner interface
func (w *BackgroundWorker) Run(ctx context.Context) error {
    w.logger.Info("Starting background worker...")

    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return nil
        case <-w.stopCh:
            return nil
        case <-ticker.C:
            // Do your work here
        }
    }
}

// Close method for cleanup
func (w *BackgroundWorker) Close(ctx context.Context) error {
    w.logger.Info("Background worker shutting down...")
    close(w.stopCh)
    return nil
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
    // Create service instances
    httpServer := createHTTPServer(ctx.Config.Port, ctx.Logger)
    grpcServer := createGRPCServer(ctx.Config.GRPCPort, ctx.Logger)
    metricsServer := createMetricsServer(ctx.Config.MetricsPort, ctx.Logger)
    worker := createBackgroundWorker(ctx.Config, ctx.Logger)

    // Define cleanup function for resource cleanup
    cleanup := func(shutdownCtx context.Context) error {
        ctx.Logger.Info("Starting concurrent cleanup of services...")

        // Use a WaitGroup to wait for all cleanup operations to complete
        var wg sync.WaitGroup

        // Create a channel to collect errors
        errCh := make(chan error, 4) // Buffer for all possible errors

        // Close all services concurrently
        wg.Add(4)

        // Close worker
        go func() {
            defer wg.Done()
            if err := worker.Close(shutdownCtx); err != nil {
                errCh <- fmt.Errorf("failed to close worker: %w", err)
            }
        }()

        // Close metrics server
        go func() {
            defer wg.Done()
            if err := metricsServer.Close(shutdownCtx); err != nil {
                errCh <- fmt.Errorf("failed to close metrics server: %w", err)
            }
        }()

        // Close gRPC server
        go func() {
            defer wg.Done()
            if err := grpcServer.Close(shutdownCtx); err != nil {
                errCh <- fmt.Errorf("failed to close gRPC server: %w", err)
            }
        }()

        // Close HTTP server
        go func() {
            defer wg.Done()
            if err := httpServer.Close(shutdownCtx); err != nil {
                errCh <- fmt.Errorf("failed to close HTTP server: %w", err)
            }
        }()

        // Wait for all goroutines to complete
        wg.Wait()
        close(errCh)

        // Collect all errors
        var errs []error
        for err := range errCh {
            errs = append(errs, err)
        }

        // Join all errors and return the combined error
        if len(errs) > 0 {
            return errors.Join(errs...)
        }
        return nil
    }

    return ezapp.Construct(
        ezapp.WithRunners(
            httpServer.Run,
            grpcServer.Run,
            metricsServer.Run,
            worker.Run,
        ),
        ezapp.WithCleanup(cleanup),
    )
}
```

### Complex Cleanup

```go
func Initialize(ctx ezapp.InitCtx[Config]) (ezapp.AppCtx, error) {
    // Setup resources
    db := setupDatabase(ctx.Config.DatabaseURL)
    cache := setupRedis(ctx.Config.RedisURL)
    messageQueue := setupMessageQueue(ctx.Config.MQUrl)
    fileStorage := setupFileStorage(ctx.Config.StorageURL)

    // Create server instance
    server := createServer(db, cache, messageQueue, fileStorage, ctx.Logger)

    // Define cleanup function for resource cleanup
    cleanup := func(shutdownCtx context.Context) error {
        ctx.Logger.Info("Starting concurrent cleanup of resources...")

        // Use a WaitGroup to wait for all cleanup operations to complete
        var wg sync.WaitGroup

        // Create a slice to collect errors
        errCh := make(chan error, 5) // Buffer for all possible errors

        // First close the server (this must be done first and sequentially)
        if err := server.Close(shutdownCtx); err != nil {
            return fmt.Errorf("failed to close server: %w", err)
        }

        // Close remaining resources concurrently
        wg.Add(4)

        // Close file storage
        go func() {
            defer wg.Done()
            if err := fileStorage.Close(); err != nil {
                errCh <- fmt.Errorf("failed to close file storage: %w", err)
            }
        }()

        // Close message queue
        go func() {
            defer wg.Done()
            if err := messageQueue.Close(); err != nil {
                errCh <- fmt.Errorf("failed to close message queue: %w", err)
            }
        }()

        // Close cache
        go func() {
            defer wg.Done()
            if err := cache.Close(); err != nil {
                errCh <- fmt.Errorf("failed to close cache: %w", err)
            }
        }()

        // Close database (last resource to be initialized, closed last)
        go func() {
            defer wg.Done()
            if err := db.Close(); err != nil {
                errCh <- fmt.Errorf("failed to close database: %w", err)
            }
        }()

        // Wait for all goroutines to complete
        wg.Wait()
        close(errCh)

        // Collect all errors
        var errs []error
        for err := range errCh {
            errs = append(errs, err)
        }

        // Join all errors and return the combined error
        if len(errs) > 0 {
            return errors.Join(errs...)
        }
        return nil
    }

    return ezapp.Construct(
        ezapp.WithRunners(server.Run),
        ezapp.WithCleanup(cleanup),
    )
}
```

### Proper Cleanup in Initializer

```go
func Initialize(ctx ezapp.InitCtx[Config]) (ezapp.AppCtx, error) {
    // Setup resources
    db := setupDatabase(ctx.Config.DatabaseURL)
    server := createHTTPServer(ctx.Config.Port, db, ctx.Logger)

    // Define cleanup function for resource cleanup
    cleanup := func(shutdownCtx context.Context) error {
        ctx.Logger.Info("Cleaning up resources...")

        // Use a WaitGroup to wait for all cleanup operations to complete
        var wg sync.WaitGroup

        // Create a channel to collect errors
        errCh := make(chan error, 2) // Buffer for all possible errors

        // Close resources concurrently
        wg.Add(2)

        // Close HTTP server
        go func() {
            defer wg.Done()
            if err := server.Close(shutdownCtx); err != nil {
                errCh <- fmt.Errorf("failed to close HTTP server: %w", err)
            }
        }()

        // Close database
        go func() {
            defer wg.Done()
            if err := db.Close(); err != nil {
                errCh <- fmt.Errorf("failed to close database: %w", err)
            }
        }()

        // Wait for all goroutines to complete
        wg.Wait()
        close(errCh)

        // Collect all errors
        var errs []error
        for err := range errCh {
            errs = append(errs, err)
        }

        // Join all errors and return the combined error
        if len(errs) > 0 {
            return errors.Join(errs...)
        }
        return nil
    }

    // Construct and return the application context
    return ezapp.Construct(
        ezapp.WithRunners(server.Run),
        ezapp.WithCleanup(cleanup),
    )
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
