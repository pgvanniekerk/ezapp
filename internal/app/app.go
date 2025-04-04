package app

import (
	"context"
	"sync"
	"time"
)

// ErrorHandler is a function that handles errors
type ErrorHandler func(error) error

type App struct {
	services        []Service
	shutdownTimeout time.Duration
	shutdownSig     <-chan struct{}
	errorHandler    ErrorHandler
}

// NewApp creates a new App with the given services, error handler, and shutdown signal
func NewApp(services []Service, errorHandler ErrorHandler, shutdownSig <-chan struct{}) *App {
	// Default error handler that panics
	if errorHandler == nil {
		errorHandler = func(err error) error {
			panic(err)
		}
	}

	return &App{
		services:        services,
		shutdownTimeout: 15 * time.Second,
		shutdownSig:     shutdownSig,
		errorHandler:    errorHandler,
	}
}

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
