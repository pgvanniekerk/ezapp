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

	// critErrChan is a write-only channel for reporting critical errors to the app.
	// This channel is automatically injected by the application framework.
	critErrChan chan<- error
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
func (r Runnable) Stop(_ context.Context) error {
	return nil
}

// NotifyCriticalError notifies the application of a critical error.
// This method sends the provided error to the critical error channel,
// which is handled by the application's critical error handler.
func (r Runnable) NotifyCriticalError(err error) {
	r.critErrChan <- err
}
