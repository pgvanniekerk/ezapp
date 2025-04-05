package ezapp

import (
	"context"
	"fmt"
)

// Run is the main entry point for starting an application built with the ezapp framework.
// It takes a Builder function, creates a default configuration, builds the application,
// and then runs it.
//
// The type parameters are:
//   - APP: A type that implements the EzApp interface
//   - CONF: Any type that represents the application configuration
//
// This function creates a background context and an empty configuration object,
// then calls the provided builder function to create the application. If the builder
// function returns an error, Run will panic with that error. Otherwise, it calls
// the Run method on the application, which blocks until the application exits.
//
// Example:
//
//	func main() {
//	    ezapp.Run(BuildFunc)
//	}
//
// Note: This function creates an empty configuration object (zero value). If your
// application requires configuration from environment variables or other sources,
// you should handle that in your builder function.
//
// Potential errors (resulting in panic):
//   - Builder function returns an error (e.g., configuration validation errors,
//     resource initialization errors, invalid runnable components)
func Run[APP EzApp, CONF any](builder Builder[APP, CONF]) {
	ctx := context.Background()

	var conf CONF

	app, err := builder(ctx, conf)
	if err != nil {
		panic(fmt.Errorf("failed to build app: %w", err))
	}

	app.Run()
}
