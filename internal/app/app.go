package app

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"
)

// App is the core application type that implements the ezapp.EzApp interface.
// It manages the lifecycle of runnable components, handling their startup,
// monitoring, and graceful shutdown.
//
// The App type is created by the wire.App function and is typically not
// instantiated directly by users of the ezapp framework.
type App struct {
	// shutdownTimeout is the maximum time allowed for stopping all runnables
	shutdownTimeout time.Duration

	// runnables is the list of components managed by this application
	runnables []Runnable

	// shutdownSig is a channel that signals when the application should shut down
	shutdownSig <-chan error

	// logger is used for application-level logging
	logger *slog.Logger

	// critErrChan is a channel that receives critical errors from runnables
	critErrChan chan error

	// criticalErrHandler is a function that handles critical errors
	criticalErrHandler func(error)
}

// Run starts the application and blocks until the application exits.
// This method implements the ezapp.EzApp interface.
//
// The Run method:
// 1. Starts each runnable component in its own goroutine
// 2. Waits for either an error from a runnable or a shutdown signal
// 3. Initiates graceful shutdown of all runnables with a timeout
// 4. Exits with status code 1 if any runnable fails to stop properly
//
// This method blocks until all runnables have been stopped or the
// shutdown timeout is reached.
func (a *App) Run() {

	a.logger.Info("Starting application")

	// Create a channel to receive errors from runnables
	errChan := make(chan error, len(a.runnables))

	// Start a goroutine to handle critical errors
	go func() {
		for err := range a.critErrChan {
			if a.criticalErrHandler != nil {
				a.logger.Error("Critical error detected, calling handler", "error", err)
				a.criticalErrHandler(err)
			}
		}
	}()

	// Start each runnable in its own goroutine
	for _, runnable := range a.runnables {
		go func(r Runnable) {
			if err := r.Run(); err != nil {
				errChan <- err
				// Also notify about the error as a critical error
				r.NotifyCriticalError(err)
			}
		}(runnable)
	}

	// Wait for an error from a runnable or a shutdown signal
	var err error
	select {
	case err = <-errChan:
		a.logger.Error("Runnable error detected, initiating shutdown", "error", err)
	case err = <-a.shutdownSig:
		a.logger.Info("Shutdown signal received, initiating shutdown", "error", err)
	}

	// When shutting down, use the shutdownTimeout to create a context with a deadline
	ctx, cancel := context.WithTimeout(context.Background(), a.shutdownTimeout)
	defer cancel()

	// Stop all runnables with the timeout context
	for _, runnable := range a.runnables {
		if stopErr := runnable.Stop(ctx); stopErr != nil {
			a.logger.Error("Failed to stop runnable", "error", stopErr)
			if errors.Is(stopErr, context.DeadlineExceeded) {
				a.logger.Error("Stop timed out and app shutdown was forced")
			}
			os.Exit(1)
		}
	}
}
