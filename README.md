# ezapp

ezapp is part of the ezGoing framework and streamlines application creation and wiring structure. It provides a lightweight framework for building Go applications with a focus on simplicity, testability, and maintainability. It also experimentally adds built-in support/context providing to the Junie AI tool - to allow Junie to know when and how to wire Runnables automatically.

## Overview

ezapp helps you build applications by:

1. Providing a consistent structure for your application components
2. Managing the lifecycle of your application components
3. Handling graceful shutdown of your application
4. Simplifying dependency injection and wiring

## Getting Started

### Installation

```bash
go get github.com/pgvanniekerk/ezapp
```

## Step 1: Recommended App Structure

The recommended structure for an ezapp application is:

```
myapp/
├── internal/
│   ├── cmd/
│   │   └── myapp.go       # Main entry point
│   ├── config/
│   │   └── config.go      # Configuration structures
│   ├── runnable/
│   │   └── runnable.go    # Example Runnable component 
│   └── wire/
│       └── wire.go        # Wiring code
└── main.go                # Calls into cmd/myapp.go
```

The `runnable` folder is just an example - you are welcome to name your service component packages anything you want to.

### Runnable Interface

The `Runnable` interface is the core building block of ezapp applications. It defines components that can be started and stopped as part of the application lifecycle.

Here's an example of a runnable component that prints the current time every second and stops when Stop is called:

```go
package runnable

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/pgvanniekerk/ezapp/pkg/ezapp"
)

// TimeRunnable is an example runnable component that prints the current time
// every second and queries the database for the system date.
type TimeRunnable struct {
	ezapp.Runnable // Embed the ezapp.Runnable struct
	db *sql.DB     // Database connection
	stopCh chan struct{} // Channel to signal stopping
}

// NewTimeRunnable creates a new TimeRunnable with the given database connection.
func NewTimeRunnable(db *sql.DB) *TimeRunnable {
	return &TimeRunnable{
		db:     db,
		stopCh: make(chan struct{}),
	}
}

// Run starts the runnable component. It prints the current time every second
// and queries the database for the system date.
// If Stop is called, Run will stop and return nil.
func (r *TimeRunnable) Run() error {
	r.Logger.Info("TimeRunnable started")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Print current time
			now := time.Now()
			r.Logger.Info("Current time", slog.Time("time", now))

			// Query database for system date
			var sysdate time.Time
			err := r.db.QueryRow("SELECT CURRENT_TIMESTAMP").Scan(&sysdate)
			if err != nil {
				r.Logger.Error("Failed to query database", slog.String("error", err.Error()))
				return err // Return error to trigger application shutdown
			}

			r.Logger.Info("Database system date", slog.Time("sysdate", sysdate))

		case <-r.stopCh:
			r.Logger.Info("TimeRunnable stopping due to stop signal")
			return nil
		}
	}
}

// Stop signals the runnable to stop and cleans up resources.
func (r *TimeRunnable) Stop(ctx context.Context) error {
	r.Logger.Info("TimeRunnable stopping")

	// Signal the Run method to stop
	close(r.stopCh)

	return nil
}
```

## Step 2: Create config/config.go

Create a configuration package with a struct compatible with the [envvar library](https://github.com/kelseyhightower/envconfig). This struct will hold configuration for your application, including database connection details.

```go
package config

import (
	"strconv"
)

// Config contains configuration for the application.
// This struct is compatible with the envconfig library
// (github.com/kelseyhightower/envconfig) for loading
// configuration from environment variables.
type Config struct {
	// Database configuration
	DB DBConf

	// LogLevel sets the application's logging level
	// Environment variable: LOG_LEVEL
	// Default: info
	LogLevel string `envvar:"LOG_LEVEL" default:"info"`

	// AppName is the name of the application
	// Environment variable: APP_NAME
	// Required: true
	AppName string `envvar:"APP_NAME" required:"true"`
}

// DBConf contains configuration for connecting to a database.
type DBConf struct {
	// Host is the database server hostname or IP address
	// Environment variable: DB_HOST
	// Default: localhost
	Host string `envvar:"DB_HOST" default:"localhost"`

	// Port is the database server port
	// Environment variable: DB_PORT
	// Default: 5432
	Port int `envvar:"DB_PORT" default:"5432"`

	// User is the database username
	// Environment variable: DB_USER
	// Default: postgres
	User string `envvar:"DB_USER" default:"postgres"`

	// Password is the database password
	// Environment variable: DB_PASSWORD
	// Default: postgres
	Password string `envvar:"DB_PASSWORD" default:"postgres"`

	// DBName is the name of the database to connect to
	// Environment variable: DB_NAME
	// Default: postgres
	DBName string `envvar:"DB_NAME" default:"postgres"`

	// SSLMode is the SSL mode to use for the connection
	// Environment variable: DB_SSL_MODE
	// Default: disable
	SSLMode string `envvar:"DB_SSL_MODE" default:"disable"`
}

// GetConnectionString returns a formatted connection string for the database.
func (c DBConf) GetConnectionString() string {
	return "host=" + c.Host + 
		" port=" + strconv.Itoa(c.Port) + 
		" user=" + c.User + 
		" password=" + c.Password + 
		" dbname=" + c.DBName + 
		" sslmode=" + c.SSLMode
}
```

## Step 3: Create wire/wire.go

Create a wire package with a Build function that matches ezapp.Builder and uses the elements/options of the wire package to create objects and provide runnable to wire.Runnables.

```go
package wire

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/yourusername/myapp/internal/config"
	"github.com/yourusername/myapp/internal/runnable"
	"github.com/pgvanniekerk/ezapp/internal/app"
	appwire "github.com/pgvanniekerk/ezapp/pkg/wire"
)

// Build creates an application with all dependencies wired up.
// This function matches the ezapp.Builder signature and is used
// with ezapp.Run to create and run the application.
func Build(startupCtx context.Context, cfg config.Config) (*app.App, error) {
	// Create database connection
	db, err := createDBConnection(startupCtx, cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	// Create runnable component
	timeRunnable := runnable.NewTimeRunnable(db)

	// Create and return the app using the wire package
	return appwire.App(
		// Provide the runnable component to the application
		appwire.Runnables(timeRunnable),

		// Configure application timeouts
		appwire.WithAppStartupTimeout(15*time.Second),
		appwire.WithAppShutdownTimeout(15*time.Second),
	)
}

// createDBConnection establishes a connection to the database.
func createDBConnection(ctx context.Context, dbConf config.DBConf) (*sql.DB, error) {
	// Open database connection
	db, err := sql.Open("postgres", dbConf.GetConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Ping the database to verify connectivity
	if err := db.PingContext(ctx); err != nil {
		db.Close() // Close the connection if ping fails
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
```

## Step 4: Create /cmd/myapp/main.go

Finally, create the main entry point for your application. This file will use the wire.Build function to create and run the application.

```go
package main

import (
	"github.com/yourusername/myapp/internal/wire"
	"github.com/pgvanniekerk/ezapp/pkg/ezapp"
)

func main() {
	// Run the application with the wire.Build function
	ezapp.Run(wire.Build)
}
```

This simple main function uses ezapp.Run to start the application with the wire.Build function. The ezapp.Run function:

1. Creates a background context
2. Creates an empty Config (in a real application, you would load this from environment variables)
3. Calls the Build function to create the application
4. Runs the application

## Wire Options

The wire.App function accepts various options for configuring your application. Here are the available options:

### WithAppShutdownTimeout

Sets the maximum time allowed for stopping all runnables during application shutdown.

```go
// Set shutdown timeout to 15 seconds
appwire.WithAppShutdownTimeout(15*time.Second)
```

### WithAppStartupTimeout

Sets the maximum time allowed for starting all runnables during application startup.

```go
// Set startup timeout to 10 seconds
appwire.WithAppStartupTimeout(10*time.Second)
```

### WithLogger

Sets the logger for the application.

```go
// Use a custom logger
appwire.WithLogger(slog.New(slog.NewTextHandler(os.Stdout, nil)))
```

### WithLogAttrs

Adds attributes to log entries.

```go
// Add application name and version to log entries
appwire.WithLogAttrs(slog.String("app", "myapp"), slog.Int("version", 1))
```

### WithShutdownSignal

Sets the channel for receiving shutdown signals.

```go
// Use a custom shutdown signal channel
appwire.WithShutdownSignal(customShutdownChan)
```

## Example Usage

Here's an example of how to use these options together:

```go
func Build(startupCtx context.Context, cfg config.Config) (*app.App, error) {
	// Create database connection
	db, err := createDBConnection(startupCtx, cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	// Create runnable components
	timeRunnable := runnable.NewTimeRunnable(db)
	httpRunnable := runnable.NewHTTPRunnable(cfg.Port)

	// Create and return the app using the wire package
	return appwire.App(
		// Provide multiple runnable components
		appwire.Runnables(timeRunnable, httpRunnable),

		// Configure application timeouts
		appwire.WithAppStartupTimeout(15*time.Second),
		appwire.WithAppShutdownTimeout(15*time.Second),

		// Configure logging
		appwire.WithLogAttrs(
			slog.String("app", cfg.AppName),
			slog.String("env", "production"),
		),
	)
}
```

## Best Practices

1. **Keep Runnables Focused**: Each runnable component should have a single responsibility.

2. **Handle Errors Properly**: Return errors from `Run` only for critical failures that should trigger application shutdown.

3. **Respect Context Cancellation**: In the `Stop` method, respect the context deadline to ensure graceful shutdown.

4. **Use Dependency Injection**: Pass dependencies to runnable components through constructor functions.

5. **Configure Timeouts Appropriately**: Set appropriate shutdown and startup timeouts based on your application's needs.

## Examples

See the [examples/exampleapp](examples/exampleapp) directory for a complete example of an ezapp application.

## License

[MIT](LICENSE)
