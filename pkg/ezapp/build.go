package ezapp

import (
	"fmt"

	"github.com/pgvanniekerk/ezapp/internal/config"
)

// Build creates a new EzApp instance using the provided wire function and options.
//
// The Build function:
// 1. Applies default and provided options
// 2. Loads configuration from environment variables into a struct of type C
// 3. Calls the wire function with the configuration to get a wire bundle
// 4. Creates and returns an EzApp with the runnables, error handler, and cleanup function
//
// The generic type parameter C represents the configuration struct type.
// Environment variables are loaded into this struct using the envconfig package.
// The struct should have `envconfig` tags to specify which environment variables to load.
//
// Example:
//
//	type Config struct {
//		Port int `envconfig:"PORT" default:"8080"`
//	}
//
//	app := ezapp.Build(
//		func(cfg Config) (ezapp.WireBundle, error) {
//			server := NewServer(cfg.Port)
//			return ezapp.WireBundle{
//				Runnables: []ezapp.Runnable{server},
//				CleanupFunc: server.Cleanup,
//			}, nil
//		},
//		ezapp.WithConfigPrefix("APP"),
//		ezapp.WithErrorHandler(customErrorHandler),
//	)
//
// If there's an error loading the configuration or calling the wire function,
// Build will attempt to handle it with the provided error handler.
// If the error persists or no error handler is provided, Build will panic.
func Build[C any](wireFunc WireFunc[C], opts ...option) *EzApp {

	// Get default options
	o := getDefaultOptions()

	// Apply all options
	for _, opt := range opts {
		opt(o)
	}

	// Create a zero value of type C
	var c C

	// Load configuration from environment variables
	err := config.Load(o.configPrefix, &c)
	if err != nil {
		// If there's an error and we have an error handler, use it
		if o.errHandler != nil {
			err = o.errHandler(err)
		}

		if err != nil {
			panic(fmt.Errorf("error loading config: %w", err))
		}
	}

	// Call the wire function to get the wire bundle
	bundle, err := wireFunc(c)
	if err != nil {

		// If there's an error and we have an error handler, use it
		if o.errHandler != nil {
			err = o.errHandler(err)
		}

		if err != nil {
			panic(fmt.Errorf("error wiring app: %w", err))
		}
	}

	// Determine which cleanup function to use
	// If the options cleanup function is nil, use the wire bundle's cleanup function
	if o.cleanupFunc == nil && bundle.CleanupFunc != nil {
		o.cleanupFunc = bundle.CleanupFunc
	}

	// Create and return the EzApp
	return &EzApp{
		runnables:    bundle.Runnables,
		errorHandler: o.errHandler,
		cleanupFunc:  o.cleanupFunc,
	}
}
