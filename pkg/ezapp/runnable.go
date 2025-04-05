package ezapp

import (
	"context"
	"log/slog"
)

// Runnable is a base struct that should be embedded in application components
// that need to be started and stopped as part of the application lifecycle.
// The ezapp framework will manage the lifecycle of all Runnables, calling Run() on startup
// and Stop() during shutdown.
//
// Structs that embed Runnable should override the Run and Stop methods
// to implement their specific startup and shutdown logic.
type Runnable struct {
	// Logger is automatically injected by the application framework
	// and can be used by the Runnable implementation for logging.
	Logger *slog.Logger
}

// Run is called by ezapp when the application starts and is executed in a separate goroutine.
// This method should be overridden by structs that embed Runnable.
//
// IMPORTANT: Returning an error from this method will trigger application
// shutdown, as it indicates a critical failure that prevents the component
// from operating correctly.
func (r Runnable) Run() error {
	return nil
}

// Stop is called by ezapp during application shutdown.
// This method should be overridden by structs that embed Runnable to
// properly clean up resources and perform graceful shutdown.
//
// The provided context may include a deadline after which the shutdown
// process will be aborted, so implementations should respect context cancellation.
func (r Runnable) Stop(ctx context.Context) error {
	return nil
}

// Sentinel is a marker method used internally by the framework to identify
// types that embed the Runnable struct.
func (r Runnable) Sentinel() {}
