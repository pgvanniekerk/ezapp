package app

// New creates a new App instance with the provided parameters.
// This function is used internally by the wire.App function and is not
// meant to be called directly by users of the ezapp framework.
//
// Parameters:
//   - params: A Params struct containing the configuration for the App.
//
// Returns:
//   - A pointer to a new App instance
//   - An error if the App could not be created
//
// The New function performs the following steps:
//  1. Applies log attributes to the logger
//  2. Creates a new App instance with the provided parameters
//  3. Validates that each runnable embeds the ezapp.Runnable struct
//  4. Sets the logger for each runnable
//
// Potential errors:
//   - Invalid runnable components (not embedding ezapp.Runnable)
func New(params Params) (*App, error) {

	// Apply log attributes to the logger if they are not empty
	logger := params.Logger
	if len(params.LogAttrs) > 0 {
		// Create a new logger with the log attributes
		for _, attr := range params.LogAttrs {
			logger = logger.With(attr.Key, attr.Value.Any())
		}
	}

	app := &App{
		shutdownTimeout: params.ShutdownTimeout,
		runnables:       params.Runnables,
		shutdownSig:     params.ShutdownSig,
		logger:          logger,
	}

	// Validate and set the logger for each runnable
	for _, runnable := range params.Runnables {

		// Validate that the runnable embeds the ezapp.Runnable struct
		if err := EnsureEmbedsRunnableStruct(runnable); err != nil {
			logger.Error("Invalid runnable", "error", err)
			return nil, err
		}

		// Set the logger for each runnable that has the toggle:"useEzAppLogger" tag
		setRunnableLogger(runnable, logger)
	}

	return app, nil
}
