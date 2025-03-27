# EzApp

EzApp is a simple, lightweight framework for building Go applications. It provides a clean and structured way to:

1. Load configuration from environment variables
2. Wire together application components
3. Run multiple components concurrently
4. Handle errors and graceful shutdown

## Installation

```bash
go get github.com/pgvanniekerk/ezapp
```

## Usage

EzApp is designed to make it easy to wire together and execute applications. Here's a basic example of how to use it:

### Recommended Structure

We recommend organizing your code with the configuration struct and wire function in their own `.go` file, separate from your main application code. This makes your code more modular and easier to maintain.

#### main.go

```go
package main

import "github.com/pgvanniekerk/ezapp/pkg/ezapp"

func main() {
	ezapp.Build(
		WireRunnables,
	).Run()
}
```

#### wire.go

```go
package main

import "github.com/pgvanniekerk/ezapp/pkg/ezapp"

type Config struct {
	// Your configuration fields here
	Port int `envconfig:"PORT" default:"8080"`
}

func WireRunnables(config Config) (ezapp.WireBundle, error) {
	runnables := make([]ezapp.Runnable, 0)

	// Create your runnables here
	server := NewServer(config.Port)
	runnables = append(runnables, server)

	return ezapp.WireBundle{
		Runnables: runnables,
		CleanupFunc: func() error {
			// Cleanup resources here
			return nil
		},
	}, nil
}
```

### Creating a Runnable

A `Runnable` is an interface for components that can be run by EzApp. Here's an example of how to implement it:

```go
package main

import (
	"context"
	"fmt"
	"net/http"
)

type Server struct {
	server *http.Server
}

func NewServer(port int) *Server {
	return &Server{
		server: &http.Server{
			Addr: fmt.Sprintf(":%d", port),
		},
	}
}

func (s *Server) Run(ctx context.Context) error {
	// Start the server in a goroutine
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.server.ListenAndServe()
	}()

	// Wait for context cancellation or server error
	select {
	case <-ctx.Done():
		// Context was canceled, shut down gracefully
		return s.server.Shutdown(context.Background())
	case err := <-errCh:
		// Server encountered an error
		return err
	}
}

func (s *Server) HandleError(err error) error {
	// Log the error and return nil to indicate it was handled
	fmt.Printf("Server error: %v\n", err)
	return nil
}
```

### Advanced Usage

EzApp provides several options to customize its behavior:

```go
package main

import (
	"fmt"
	"github.com/pgvanniekerk/ezapp/pkg/ezapp"
)

func main() {
	app := ezapp.Build(
		WireRunnables,
		ezapp.WithConfigPrefix("APP"),
		ezapp.WithErrorHandler(func(err error) error {
			fmt.Printf("Error: %v\n", err)
			return nil
		}),
		ezapp.WithCleanupFunc(func() error {
			fmt.Println("Cleaning up resources...")
			return nil
		}),
	)

	app.Run()
}
```

## Features

- **Environment-based Configuration**: Automatically loads configuration from environment variables using the [envconfig](https://github.com/kelseyhightower/envconfig) package.
- **Graceful Shutdown**: Handles SIGINT and SIGTERM signals to gracefully shut down your application.
- **Concurrent Execution**: Runs multiple components concurrently in separate goroutines.
- **Error Handling**: Provides a structured way to handle errors at both the component and application levels.
- **Resource Cleanup**: Ensures proper cleanup of resources when the application terminates.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
