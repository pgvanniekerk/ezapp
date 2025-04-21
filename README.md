# EzApp

EzApp is a simple, lightweight framework for building Go applications. It provides a clean and structured way to:

1. Load configuration from environment variables
2. Wire together application components
3. Run multiple services concurrently
4. Handle errors and graceful shutdown

## Installation

```bash
go get github.com/pgvanniekerk/ezapp
```

## Usage

EzApp is designed to make it easy to wire together and execute applications. Here's a basic example of how to use it:

### Recommended Structure

We recommend organizing your code with the configuration struct and builder function in their own `.go` file in an app package, separate from your main application code. This makes your code more modular and easier to maintain. It's especially recommended to have the builder function located in the app package rather than the main package.

#### main.go

```go
package main

import (
	"github.com/pgvanniekerk/ezapp/pkg/ezapp"
	"myapp/app"
)

func main() {
  ezapp.Run(app.Builder)
}
```

#### app/builder.go

```go
package app

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pgvanniekerk/ezapp/pkg/ezapp"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// Logger is a simple logging interface
type Logger interface {
	Info(format string, args ...interface{})
	Error(format string, args ...interface{})
}

// SimpleLogger implements the Logger interface
type SimpleLogger struct{}

// NewLogger creates a new SimpleLogger
func NewLogger() *SimpleLogger {
	return &SimpleLogger{}
}

// Info logs an informational message
func (l *SimpleLogger) Info(format string, args ...interface{}) {
	fmt.Printf("[INFO] "+format+"\n", args...)
}

// Error logs an error message
func (l *SimpleLogger) Error(format string, args ...interface{}) {
	fmt.Printf("[ERROR] "+format+"\n", args...)
}

type Config struct {
	// Your configuration fields here
	Port        int    `env:"PORT" default:"8080"`
	DatabaseURL string `env:"DATABASE_URL" required:"true"`
}

// Builder is the function used for dependency injection and wiring
// It connects services and dependencies together
func Builder(config Config) ([]ezapp.Runnable, error) {
	// Connect to the database
	db, err := sql.Open("postgres", config.DatabaseURL)
	if err != nil {
		return nil, err
	}

	// Create a logger (not shown in this basic example)
	logger := NewLogger()

	// Create a server with the database connection
	server := NewServer(config.Port, db, logger, func() error {
		return db.Close()
	})

	// Return the list of runnables
	return []ezapp.Runnable{server}, nil
}
```

### Runnable Interface

To be used with EzApp, your services must implement the `Runnable` interface:

```go
type Runnable interface {
	Run(context.Context) error
}
```

The `Run(context.Context)` method starts the service and should only return an error in exceptional circumstances such as dependency failures or timeouts (application-impacting errors). The context parameter can be used to detect when the application is shutting down, allowing for graceful termination of long-running operations.

**Note:** The `EzApp.Run()` method itself does not return an error. Instead, it prints appropriate messages to the console about why the app is shutting down and handles graceful exit in all cases.

Here's an example of how to implement the Runnable interface:

```go
package app

import (
	"context"
	"fmt"
	"net/http"
	"database/sql"
)

type Server struct {
	server  *http.Server
	logger  Logger
	cleanup func() error
}

func NewServer(port int, db *sql.DB, logger Logger, cleanup func() error) *Server {
	return &Server{
		server: &http.Server{
			Addr: fmt.Sprintf(":%d", port),
		},
		logger: logger,
		cleanup: cleanup,
	}
}

func (s *Server) Run(ctx context.Context) error {
	s.logger.Info("Starting server on %s", s.server.Addr)

	// Start a goroutine to listen for context cancellation
	go func() {
		<-ctx.Done()
		// Context was cancelled, shut down the server
		s.logger.Info("Shutting down server")
		s.server.Shutdown(context.Background())
		// Clean up resources
		if s.cleanup != nil {
			s.logger.Info("Cleaning up resources")
			s.cleanup()
		}
	}()

	// Start the server
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		s.logger.Error("Server error: %v", err)
		return err
	}
	return nil
}
```

## Features

- **Environment-based Configuration**: Automatically loads configuration from environment variables using the [go-env](https://github.com/Netflix/go-env) package.
- **Graceful Shutdown**: Handles SIGINT and SIGTERM signals to gracefully shut down your application.
- **Concurrent Execution**: Runs multiple services concurrently in separate goroutines.
- **Error Handling**: Provides a structured way to handle errors at both the service and application levels.
- **Type Safety with Generics**: Uses Go generics to provide type safety for your configuration.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
