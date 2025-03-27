// Package ezapp provides a simple framework for building and running applications.
//
// The ezapp package is designed to make it easy to wire together and execute applications
// by providing a clean and structured way to:
// 1. Load configuration from environment variables
// 2. Wire together application components
// 3. Run multiple components concurrently
// 4. Handle errors and graceful shutdown
//
// The main components of the package are:
// - EzApp: The main application struct that runs the application components
// - Build: A function that creates an EzApp instance
// - Runnable: An interface for components that can be run by the EzApp
// - WireFunc: A function that wires together application components
//
// Example usage:
//
//	type Config struct {
//		Port int `envconfig:"PORT" default:"8080"`
//	}
//
//	func main() {
//		app := ezapp.Build(wireApp)
//		app.Run()
//	}
//
//	func wireApp(cfg Config) (ezapp.WireBundle, error) {
//		server := NewServer(cfg.Port)
//		return ezapp.WireBundle{
//			Runnables: []ezapp.Runnable{server},
//			CleanupFunc: server.Cleanup,
//		}, nil
//	}
package ezapp

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// EzApp is the main application struct that runs the application components.
// It manages the lifecycle of multiple Runnable components, handles errors,
// and ensures proper cleanup when the application terminates.
type EzApp struct {
	runnables    []Runnable    // Components to be run concurrently
	errorHandler ErrHandler    // Function to handle errors from Runnables
	cleanupFunc  CleanupFunc   // Function to perform cleanup operations
}

// Run starts all the Runnable components concurrently and waits for them to complete.
//
// The Run method:
// 1. Sets up a deferred cleanup function to ensure resources are properly released
// 2. Creates a context with a cancel function for coordinating the Runnables
// 3. Sets up signal handling for graceful shutdown on Ctrl+C (SIGINT) or SIGTERM
// 4. Runs each Runnable in a separate goroutine
// 5. Handles errors from the Runnables using their HandleError methods
// 6. If errors persist, uses the app-level error handler
// 7. If errors still persist, cancels the context to initiate shutdown
// 8. Waits for all goroutines to finish before returning
//
// This method blocks until all Runnables have completed or been canceled.
// If the cleanup function returns an error, Run will panic.
func (e *EzApp) Run() {

	// Perform any cleanup operations after we're done
	defer func() {
		err := e.cleanupFunc()
		if err != nil {
			panic(err)
		}
	}()

	// Create a context with cancel function
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	// Create a wait group to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Run each Runnable in a separate goroutine
	for _, r := range e.runnables {
		wg.Add(1)
		go func(runnable Runnable) {
			defer wg.Done()

			// Run the Runnable
			err := runnable.Run(ctx)
			if err != nil && !errors.Is(err, context.Canceled) {

				// Handle the error with the Runnable's error handler
				err = runnable.HandleError(err)
				if err != nil && e.errorHandler != nil {

					// If the error persists and we have an app-level error handler, use it
					err = e.errorHandler(err)
					if err != nil {

						// If the error still persists, cancel the context
						cancel()
					}
				}
			}
		}(r)
	}

	// Wait for all goroutines to finish
	wg.Wait()

}
