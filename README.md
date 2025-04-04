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

We recommend organizing your code with the configuration struct and wire function in their own `.go` file, separate from your main application code. This makes your code more modular and easier to maintain.

#### main.go

```go
package main

import (
	"github.com/pgvanniekerk/ezapp/pkg/buildoption"
	"github.com/pgvanniekerk/ezapp/pkg/ezapp"
)

func main() {
  ezapp.Build(WireFunc, buildoption.WithoutOptions()).Run()
}
```

#### wire.go

```go
package main

import (
	"context"
	"database/sql"
	"github.com/pgvanniekerk/ezapp/pkg/ezapp"
	_ "github.com/lib/pq" // PostgreSQL driver
)

type Config struct {
	// Your configuration fields here
	Port        int    `envconfig:"PORT" default:"8080"`
	DatabaseURL string `envconfig:"DATABASE_URL" required:"true"`
}

// WireFunc is the function used for dependency injection and wiring
// It connects services and dependencies together
func WireFunc(startupCtx context.Context, config Config) (ezapp.ServiceSet, error) {
	// Connect to the database using the startup context
	db, err := sql.Open("postgres", config.DatabaseURL)
	if err != nil {
		return ezapp.ServiceSet{}, err
	}

	// Use the context for database operations
	if err := db.PingContext(startupCtx); err != nil {
		return ezapp.ServiceSet{}, err
	}

	// Create your services here
	server := NewServer(config.Port, db)

	return ezapp.NewServiceSet(
		ezapp.WithServices(server),
		ezapp.WithCleanupFunc(func() error {
			// Cleanup resources here, including closing the database
			return db.Close()
		}),
	), nil
}
```

### Creating a Service

A `Service` is an interface for components that can be run by EzApp. The interface is defined as follows:

```go
type Service interface {
	Run() error
	Stop(context.Context) error
}
```

The `Run()` method starts the service and should only return an error in exceptional circumstances such as dependency failures or timeouts (application-impacting errors). The `Stop(context.Context)` method stops the service, taking a context as a parameter. If it returns an error, it will be reported during shutdown. If the context timeout is reached, the application will force close.

Here's an example of how to implement it:

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

func (s *Server) Run() error {
	// Start the server
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	// Gracefully shut down the server
	return s.server.Shutdown(ctx)
}
```

### Dependency Injection Example

Here's a more comprehensive example of dependency injection with multiple services:

```go
// WireFunc is the function used for dependency injection and wiring
// It connects services and dependencies together
func WireFunc(startupCtx context.Context, config Config) (ezapp.ServiceSet, error) {
    // Connect to the database using the startup context
    db, err := sql.Open("postgres", config.DatabaseURL)
    if err != nil {
        return ezapp.ServiceSet{}, err
    }

    // Use the context for database operations
    if err := db.PingContext(startupCtx); err != nil {
        return ezapp.ServiceSet{}, err
    }

    // Create additional dependencies
    cache := cache.New(config.CacheURL)

    // Create services with dependencies
    userService := services.NewUserService(db, cache)
    authService := services.NewAuthService(db, userService)

    return ezapp.NewServiceSet(
        ezapp.WithServices(userService, authService),
        ezapp.WithCleanupFunc(func() error {
            // Ensure the database is properly closed during cleanup
            return db.Close()
        }),
    ), nil
}
```

In this example:
1. We create shared dependencies (database and cache)
2. We create services that depend on these shared resources
3. We also create services that depend on other services (authService depends on userService)
4. We provide a cleanup function to close the database connection when the application shuts down

This pattern allows for clean separation of concerns and makes testing easier by allowing dependencies to be mocked.

### Advanced Usage: Functional Options

EzApp provides functional options to customize its behavior:

#### Custom Error Handler

```go
package main

import (
	"log"
	"github.com/pgvanniekerk/ezapp/pkg/buildoption"
	"github.com/pgvanniekerk/ezapp/pkg/ezapp"
)

// CustomErrorHandler logs errors instead of panicking
func CustomErrorHandler(err error) error {
	log.Printf("Error occurred: %v", err)
	return err
}

func main() {
	ezapp.Build(
		WireFunc,
		buildoption.WithOptions(buildoption.WithErrorHandler(CustomErrorHandler)),
	).Run()
}
```

#### Custom Startup Timeout

```go
package main

import (
	"time"
	"github.com/pgvanniekerk/ezapp/pkg/buildoption"
	"github.com/pgvanniekerk/ezapp/pkg/ezapp"
)

func main() {
	ezapp.Build(
		WireFunc,
		buildoption.WithOptions(buildoption.WithStartupTimeout(30 * time.Second)), // Default is 15 seconds
	).Run()
}
```

#### Custom Environment Variable Prefix

```go
package main

import (
	"github.com/pgvanniekerk/ezapp/pkg/buildoption"
	"github.com/pgvanniekerk/ezapp/pkg/ezapp"
)

func main() {
	ezapp.Build(
		WireFunc,
		buildoption.WithOptions(buildoption.WithEnvVarPrefix("MYAPP")), // Default is empty string
	).Run()
}
```

#### Using Multiple Options

```go
package main

import (
	"github.com/pgvanniekerk/ezapp/pkg/buildoption"
	"github.com/pgvanniekerk/ezapp/pkg/ezapp"
	"time"
)

func main() {
	ezapp.Build(
		WireFunc,
		buildoption.WithOptions(
			buildoption.WithErrorHandler(CustomErrorHandler),
			buildoption.WithStartupTimeout(30 * time.Second),
			buildoption.WithEnvVarPrefix("MYAPP"),
		),
	).Run()
}
```

#### Using Default Options

You can also use `buildoption.WithoutOptions()` to get the default options without any customization:

```go
package main

import (
	"github.com/pgvanniekerk/ezapp/pkg/buildoption"
	"github.com/pgvanniekerk/ezapp/pkg/ezapp"
)

func main() {
	ezapp.Build(
		WireFunc,
		buildoption.WithoutOptions(),
	).Run()
}
```

## Features

- **Environment-based Configuration**: Automatically loads configuration from environment variables using the [envconfig](https://github.com/kelseyhightower/envconfig) package. Currently, this is the only supported method for configuration at app startup.
- **Graceful Shutdown**: Handles SIGINT and SIGTERM signals to gracefully shut down your application.
- **Concurrent Execution**: Runs multiple services concurrently in separate goroutines.
- **Error Handling**: Provides a structured way to handle errors at both the service and application levels.
- **Resource Cleanup**: Ensures proper cleanup of resources when the application terminates.
- **Functional Options**: Customize application behavior with functional options, including:
  - Custom error handlers
  - Custom startup timeout
  - Custom environment variable prefix

## Examples

Check out the examples directory for more detailed examples:

- [Example App](examples/exampleapp/README.md): A complete application example that demonstrates all features of ezapp, including dependency injection, service implementation, and resource cleanup.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
