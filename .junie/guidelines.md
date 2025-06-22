# EzApp Framework Guidelines

## Project Overview

EzApp is a lightweight, opinionated Go framework designed to simplify the process of building robust applications with minimal boilerplate code. It provides a structured approach to handling common application concerns such as configuration management, logging, concurrent service execution, graceful shutdown, and resource cleanup.

## Purpose

The primary goal of EzApp is to eliminate repetitive setup code found in most Go applications by providing sensible defaults and a clear structure. This allows developers to focus on implementing business logic rather than infrastructure concerns.

## Tech Stack

- **Language**: Go (Golang)
- **Dependencies**:
  - [go-env](https://github.com/Netflix/go-env) - For environment variable configuration
  - [zap](https://github.com/uber-go/zap) - For structured logging
  - [errgroup](https://golang.org/x/sync/errgroup) - For concurrent execution with error handling

## Project Structure

```
ezapp/
├── ezapp.go                  # Main package API
├── ezapp_test.go             # Tests for the main package
├── go.mod                    # Go module definition
├── go.sum                    # Go module checksums
└── internal/                 # Internal implementation details
    ├── app/                  # Application runtime components
    │   ├── app.go            # Core application functionality
    │   ├── app_test.go       # Tests for app functionality
    │   └── runner.go         # Runner type definition
    └── config/               # Configuration components
        ├── loadlogger.go     # Logger initialization
        ├── loadlogger_test.go # Tests for logger initialization
        ├── loadvar.go        # Environment variable loading
        ├── loadvar_test.go   # Tests for env var loading
        ├── startupctx.go     # Startup context creation
        └── startupctx_test.go # Tests for startup context
```

## Core Components

### 1. Public API (ezapp.go)

The main package exposes the following key components:

- **InitCtx[Config]**: Generic type that provides initialization context including configuration, logging, and startup context.
- **AppCtx**: Represents the application context containing all runners to be executed.
- **Initializer[Config]**: Function type for application initialization logic.
- **WithRunners**: Functional option to add runners to the application.
- **WithCleanup**: Functional option to set a cleanup function.
- **Construct**: Builds an AppCtx using provided functional options.
- **Run**: Main entry point for starting an EzApp application.

### 2. Application Runtime (internal/app)

- **Runner** (runner.go): Type definition for functions that can be executed by the application.
- **App** (app.go): Core implementation that handles concurrent execution of runners, signal handling, and graceful shutdown.

### 3. Configuration Management (internal/config)

- **LoadVar** (loadvar.go): Generic function to load configuration from environment variables.
- **LoadLogger** (loadlogger.go): Creates and configures a structured logger.
- **StartupCtx** (startupctx.go): Creates contexts with appropriate timeouts for application startup and shutdown.

## Application Lifecycle

EzApp manages the complete application lifecycle in this order:

1. **Configuration Loading**: Loads configuration from environment variables.
2. **Logger Initialization**: Creates a structured logger with configurable log level.
3. **Startup Context Creation**: Creates a context with configurable timeout.
4. **Application Initialization**: Calls the initializer function to wire dependencies.
5. **Concurrent Execution**: Runs all runners concurrently with signal handling.
6. **Graceful Shutdown**: Cancels context to signal all runners to stop.
7. **Resource Cleanup**: Calls cleanup function with shutdown timeout.
8. **Exit**: Logs completion status and exits with appropriate code.

## Key Design Patterns

1. **Functional Options**: Used for flexible and composable configuration of the application context.
2. **Dependency Injection**: Configuration, logger, and context are injected into the initializer function.
3. **Context Propagation**: Contexts are used for timeout management and cancellation signals.
4. **Error Groups**: Used for concurrent execution with coordinated error handling.
5. **Graceful Shutdown**: Signal handling and context cancellation for clean application termination.

## Environment Variables

### Framework Variables

- `EZAPP_LOG_LEVEL`: Controls logging verbosity (DEBUG, INFO, WARN, ERROR, etc.)
- `EZAPP_STARTUP_TIMEOUT`: Timeout in seconds for initialization (default: 15)
- `EZAPP_SHUTDOWN_TIMEOUT`: Timeout in seconds for graceful shutdown (default: 15)

### Application Variables

Applications can define their own configuration variables using struct tags:

```go
type Config struct {
    Port        int    `env:"PORT" default:"8080"`
    DatabaseURL string `env:"DATABASE_URL" required:"true"`
    LogLevel    string `env:"LOG_LEVEL" default:"INFO"`
}
```

## Best Practices for Using EzApp

1. **Keep initializer separate**: Put your initializer function in a separate file.
2. **Design for graceful shutdown**: All runners should respect context cancellation.
3. **Use structured logging**: Leverage the provided zap.Logger for consistent logging.
4. **Handle cleanup properly**: Release resources in reverse order of acquisition.
5. **Fail fast**: Validate dependencies early in the initializer.
6. **Configure via environment**: Use environment variables for all configuration.

## Common Use Cases

EzApp is well-suited for:

- **Microservices**: Provides a consistent structure for service implementation.
- **CLI Tools**: Simplifies configuration and error handling.
- **Web Applications**: Handles HTTP server lifecycle and graceful shutdown.
- **Background Workers**: Manages worker lifecycle and signal handling.

## Integration Points

When extending or integrating with EzApp, focus on these key points:

1. **Runner Functions**: Implement business logic as functions that respect context cancellation.
2. **Configuration Structs**: Define configuration using appropriate struct tags.
3. **Initializer Function**: Wire dependencies and create runners in a single place.
4. **Cleanup Function**: Ensure proper resource release during shutdown.