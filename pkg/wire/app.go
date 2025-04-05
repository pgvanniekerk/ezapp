package wire

import (
	"fmt"
	"github.com/pgvanniekerk/ezapp/internal/app"
)

// App creates a new application instance with the provided runnables and options.
// This is the main function for wiring up an application in the ezapp framework.
//
// Parameters:
//   - runnablesFunc: A function that returns a slice of Runnable components.
//     This is typically created using the wire.Runnables function.
//   - opts: Optional AppOption functions for configuring the application.
//     Common options include WithAppShutdownTimeout, WithAppStartupTimeout,
//     WithLogger, WithLogAttrs, and WithShutdownSignal.
//
// Returns:
//   - A pointer to an app.App instance that implements the ezapp.EzApp interface
//   - An error if the application could not be created
//
// Example:
//
//	app, err := wire.App(
//	    wire.Runnables(myRunnable1, myRunnable2),
//	    wire.WithAppShutdownTimeout(15*time.Second),
//	    wire.WithAppStartupTimeout(10*time.Second),
//	)
//
// Potential errors:
//   - runnablesFunc is nil
//   - Failed to retrieve default options
//   - Invalid runnable components (not embedding ezapp.Runnable)
//   - Other errors from app.New
func App(runnablesFunc func() []app.Runnable, opts ...AppOption) (*app.App, error) {

	// Check if runnablesFunc is nil
	if runnablesFunc == nil {
		return nil, fmt.Errorf("runnablesFunc cannot be nil")
	}

	// Apply default options
	options, err := defaultOptions()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve default options for app: %w", err)
	}

	// Apply user-provided options
	for _, opt := range opts {
		opt(options)
	}

	// Apply log attributes to the logger if they are not empty
	if len(options.logAttrs) > 0 {

		// Create a new logger with the log attributes
		for _, attr := range options.logAttrs {
			options.logger = options.logger.With(attr.Key, attr.Value.Any())
		}
	}

	// Get the runnables
	runnables := runnablesFunc()

	// Create a new app with the configured parameters
	params := app.Params{
		ShutdownTimeout: options.appConf.ShutdownTimeout,
		Runnables:       runnables,
		ShutdownSig:     options.shutdownSig,
		Logger:          options.logger,
		LogAttrs:        options.logAttrs,
	}

	appInstance, err := app.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create app: %w", err)
	}

	return appInstance, nil
}
