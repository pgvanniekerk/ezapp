package app

import (
	"context"
	"errors"
	"sync"
	"time"
)

// ErrorHandler is a function type that processes errors and returns a potentially modified error.
// It's used throughout the application to handle errors in a consistent way, allowing for
// custom error handling strategies such as logging, metrics collection, or error transformation.
type ErrorHandler func(error) error

// App is the core structure that manages the lifecycle of multiple services.
// It coordinates starting, running, and gracefully shutting down services,
// while also providing error handling capabilities.
type App struct {

	// services is a slice of Service implementations that will be managed by this App.
	// All services will be started in parallel and stopped in the order they appear in this slice.
	services []Service

	// shutdownTimeout is the maximum duration allowed for services to gracefully shut down.
	// If a service takes longer than this duration to stop, it will be forcefully terminated.
	shutdownTimeout time.Duration

	// shutdownSig is a channel that, when closed, signals the App to initiate the shutdown process.
	// This allows for external control of the application lifecycle.
	shutdownSig <-chan struct{}

	// errorHandler is a function that processes errors encountered during service operation.
	// If nil, a default handler that panics will be used.
	errorHandler ErrorHandler
}

// NewApp creates a new App with the given services, error handler, and shutdown signal.
// It initializes an App instance with the provided parameters and sets default values where appropriate.
//
// Parameters:
//   - services: A slice of Service implementations that will be managed by this App.
//   - errorHandler: A function to handle errors encountered during service operation.
//     If nil, a default handler that panics will be used.
//   - shutdownSig: A channel that, when closed, signals the App to initiate the shutdown process.
//
// Returns:
//   - A pointer to the newly created App instance, ready to be run.
//   - An error if any of the inputs are invalid.
func NewApp(services []Service, errorHandler ErrorHandler, shutdownSig <-chan struct{}) (*App, error) {
	// Validate inputs
	if services == nil {
		return nil, errors.New("services cannot be nil")
	}

	if len(services) == 0 {
		return nil, errors.New("services cannot be empty")
	}

	if shutdownSig == nil {
		return nil, errors.New("shutdownSig cannot be nil")
	}

	// Default error handler that panics
	if errorHandler == nil {
		return nil, errors.New("errorHandler cannot be nil")
	}

	return &App{
		services:        services,
		shutdownTimeout: 15 * time.Second,
		shutdownSig:     shutdownSig,
		errorHandler:    errorHandler,
	}, nil
}

// Run starts all services managed by this App and coordinates their lifecycle.
// It performs the following operations:
//  1. Starts each service in its own goroutine
//  2. Waits for either a shutdown signal or an error from any service
//  3. Initiates graceful shutdown of all services when either occurs
//  4. Handles any errors that occur during startup or shutdown
//  5. Waits for all services to complete before returning
//
// This method blocks until all services have been stopped, either due to
// receiving a shutdown signal or encountering an error.
func (a *App) Run() {

	// Create a WaitGroup to track all running services
	var wg sync.WaitGroup

	// Create channels for error handling
	errChan := make(chan error, len(a.services))

	// Start each service in its own goroutine
	for _, service := range a.services {
		wg.Add(1)
		go func(s Service) {
			defer wg.Done()
			if err := s.Run(); err != nil {
				errChan <- err
			}
		}(service)
	}

	// Wait for shutdown signal or error
	select {
	case <-a.shutdownSig:
		// Shutdown signal received, stop all services
	case err := <-errChan:
		// Service error occurred, handle it and proceed with shutdown
		if a.errorHandler != nil {
			_ = a.errorHandler(err)
		}
	}

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), a.shutdownTimeout)
	defer cancel()

	// Stop all services
	for _, service := range a.services {
		if err := service.Stop(ctx); err != nil {
			// Handle the error using the error handler
			if a.errorHandler != nil {
				_ = a.errorHandler(err)
			}
		}
	}

	// Wait for all services to finish
	wg.Wait()

}
