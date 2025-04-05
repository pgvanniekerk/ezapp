package wire

import (
	"context"
	"database/sql"
	"fmt"
	// _ "github.com/lib/pq" // PostgreSQL driver - uncomment and run go get github.com/lib/pq in a real application
	"github.com/pgvanniekerk/ezapp/examples/exampleapp/internal/config"
	"github.com/pgvanniekerk/ezapp/examples/exampleapp/internal/runnable"
	"github.com/pgvanniekerk/ezapp/internal/app"
	appwire "github.com/pgvanniekerk/ezapp/pkg/wire"
	"time"
)

// Wire creates an application with all dependencies wired up.
// This function demonstrates how to use the wire package to create an application
// with its dependencies properly initialized and configured.
//
// The Wire function is responsible for:
// 1. Creating and initializing all application components (database connections, services, etc.)
// 2. Creating runnable components that will be managed by the application
// 3. Configuring the application with options like timeouts
//
// Parameters:
//   - startupCtx: A context that can be used for initialization operations
//   - dbConf: Database configuration for connecting to PostgreSQL
//
// Returns:
//   - A pointer to an app.App instance that implements the ezapp.EzApp interface
//   - An error if the application could not be created
//
// Potential errors:
//   - Failed to create database connection
//   - Failed to create runnable components
//   - Failed to create application
func Wire(startupCtx context.Context, dbConf config.DBConf) (*app.App, error) {

	// Create database connection
	db, err := createDBConnection(startupCtx, dbConf)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection: %w", err)
	}

	// RunnableDB is a mock object to showcase the use of the Runnable interface and supplying it
	// to ezapp for management. In a real application, you would create runnable components
	// for your application services, like HTTP servers, message consumers, etc.
	runnableDB := runnable.NewDBRunnable(db)

	// Create and return the app using the wire package.
	// The App function takes a function that returns a slice of runnables and
	// optional configuration options.
	return appwire.App(
		// Runnables creates a function that returns the provided runnables.
		// You can provide multiple runnables as variadic arguments.
		appwire.Runnables(runnableDB),

		// WithAppStartupTimeout sets the timeout for application startup.
		// If any runnable takes longer than this to start, the application will fail.
		appwire.WithAppStartupTimeout(15*time.Second),

		// WithAppShutdownTimeout sets the timeout for application shutdown.
		// If any runnable takes longer than this to stop, it will be forcibly terminated.
		appwire.WithAppShutdownTimeout(15*time.Second),
	)
}

// createDBConnection establishes a connection to the database and pings it to verify connectivity.
// This function demonstrates best practices for creating and configuring database connections:
// 1. Open the connection with the appropriate driver and connection string
// 2. Ping the database to verify connectivity
// 3. Configure connection pool settings for optimal performance
//
// Parameters:
//   - ctx: A context that can be used to cancel the ping operation
//   - dbConf: Database configuration for connecting to PostgreSQL
//
// Returns:
//   - A pointer to a sql.DB instance representing the database connection
//   - An error if the connection could not be established
//
// Potential errors:
//   - Failed to open database connection (invalid connection string, driver not found)
//   - Failed to ping database (database not reachable, invalid credentials)
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
	// These settings are important for production applications to ensure
	// efficient use of database connections and prevent connection leaks.
	db.SetMaxOpenConns(25)                 // Maximum number of open connections to the database
	db.SetMaxIdleConns(5)                  // Maximum number of idle connections in the pool
	db.SetConnMaxLifetime(5 * time.Minute) // Maximum amount of time a connection may be reused

	return db, nil
}
