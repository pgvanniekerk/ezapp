package app

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func New(
	runnables []Runnable,
	logger *slog.Logger,
) *App {

	// Set up a stop context that will be cancelled when an unhandled error occurs
	// or when the os signals SIGTERM.
	stopCtx, stopCtxCancel := context.WithCancel(context.Background())

	// Initialize signal channel to handle Ctrl+C (SIGINT) and SIGTERM.
	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGINT, syscall.SIGTERM)

	return &App{
		runnables:     runnables,
		sigTerm:       sigTerm,
		stopCtx:       stopCtx,
		stopCtxCancel: stopCtxCancel,
		logger:        logger,
	}
}

type App struct {
	runnables     []Runnable
	sigTerm       chan os.Signal
	stopCtx       context.Context
	stopCtxCancel context.CancelFunc
	logger        *slog.Logger
}

func (a *App) Run() error {
	a.logger.Debug("Starting application")

	// Create an error group with the stop context
	g, ctx := errgroup.WithContext(a.stopCtx)

	// Start each runnable in its own goroutine
	for _, r := range a.runnables {
		runnable := r // Create a new variable to avoid closure issues
		a.logger.Debug("Starting runnable", "runnable", runnable)

		g.Go(func() error {
			// Run the runnable with the stop context
			err := runnable.Run(ctx)
			if err != nil {
				a.logger.Debug("Runnable returned error", "error", err)
				// Handle the error
				handledErr := runnable.HandleError(err)
				if handledErr != nil {
					a.logger.Debug("Error handling failed", "error", handledErr)
					// Cancel the context to notify other runnables to stop
					a.stopCtxCancel()
					return handledErr
				}
			}
			return nil
		})
	}

	// Create a channel to signal when the error group is done
	done := make(chan error, 1)

	// Start a goroutine to wait for the error group to complete
	go func() {
		done <- g.Wait()
	}()

	// Wait for either a sigTerm signal, the context to be done, or the error group to complete
	var err error
	select {
	case <-a.sigTerm:
		a.logger.Debug("Received termination signal")
		a.stopCtxCancel() // Cancel the context to notify runnables to stop
		err = <-done // Wait for all goroutines to complete
	case <-ctx.Done():
		a.logger.Debug("Context done", "error", ctx.Err())
		// Context was cancelled, either by an error or by the sigTerm handler
		err = <-done // Wait for all goroutines to complete
	case err = <-done:
		a.logger.Debug("All runnables completed")
	}

	// Check if the error is context.Canceled, which we can ignore
	if err != nil && !errors.Is(err, context.Canceled) {
		a.logger.Debug("Error group returned error", "error", err)
		return err
	}

	a.logger.Debug("Application stopped")
	return nil
}
