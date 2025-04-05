package runnable

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/pgvanniekerk/ezapp/pkg/ezapp"
)

// DBRunnable is an example runnable component that depends on a database connection.
// This struct demonstrates how to create a runnable component that will be managed
// by the ezapp framework.
//
// A runnable component is any struct that embeds ezapp.Runnable and overrides
// the Run and Stop methods. The ezapp framework will call these methods to
// start and stop the component as part of the application lifecycle.
//
// This example shows how to:
// 1. Embed the ezapp.Runnable struct to inherit its behavior
// 2. Add additional fields for dependencies (like a database connection)
// 3. Implement the Run method to start the component
// 4. Implement the Stop method to gracefully shut down the component
type DBRunnable struct {
	ezapp.Runnable         // Embed the ezapp.Runnable struct to inherit its behavior
	db             *sql.DB // Database connection dependency
}

// NewDBRunnable creates a new DBRunnable with the given database connection.
// This is a constructor function that follows the common Go pattern of
// providing a New function to create instances of a type.
//
// Parameters:
//   - db: A pointer to a sql.DB instance representing the database connection
//
// Returns:
//   - A pointer to a new DBRunnable instance
func NewDBRunnable(db *sql.DB) *DBRunnable {
	return &DBRunnable{
		db: db,
	}
}

// Run implements the Runnable interface and is called when the application starts.
// This method is executed in a separate goroutine by the ezapp framework.
//
// In this example, the Run method:
// 1. Logs that the component has started
// 2. Queries the database to verify the connection
// 3. Logs the database version if successful
//
// IMPORTANT: Returning an error from this method will trigger application
// shutdown, as it indicates a critical failure that prevents the component
// from operating correctly.
//
// Returns:
//   - nil if the component started successfully
//   - an error if the component failed to start (which will trigger application shutdown)
func (r *DBRunnable) Run() error {
	r.Logger.Info("DBRunnable started")

	// Example: Query the database to verify connection
	var version string
	err := r.db.QueryRow("SELECT version()").Scan(&version)
	if err != nil {
		r.Logger.Error("Failed to query database", slog.String("error", err.Error()))
		return err // Returning an error will trigger application shutdown
	}

	r.Logger.Info("Database connection verified", slog.String("version", version))
	return nil
}

// Stop implements the Runnable interface and is called when the application is shutting down.
// This method is responsible for gracefully stopping the component and cleaning up any resources.
//
// In this example, the Stop method:
// 1. Logs that the component is stopping
// 2. Notes that the database connection will be closed by the app
//
// The provided context may include a deadline after which the shutdown
// process will be aborted, so implementations should respect context cancellation.
//
// Parameters:
//   - ctx: A context that may include a deadline for the shutdown operation
//
// Returns:
//   - nil if the component stopped successfully
//   - an error if the component failed to stop
func (r *DBRunnable) Stop(_ context.Context) error {
	r.Logger.Info("DBRunnable stopping")

	// Example: Close any resources that need to be cleaned up
	// The database connection will be closed by the app, so we don't need to close it here

	return nil
}
