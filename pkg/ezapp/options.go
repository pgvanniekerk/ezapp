package ezapp

// options is a struct that holds configuration options for the EzApp.
//
// The options struct is used internally by the Build function to configure
// the behavior of the EzApp. It is not exported and should not be used directly.
// Instead, use the provided option functions to configure the EzApp.
type options struct {
	errHandler   ErrHandler  // Function to handle errors from Runnables
	configPrefix string      // Prefix for environment variables when loading configuration
}

// getDefaultOptions returns a new options struct with default values.
//
// The default values are:
// - errHandler: nil (no error handler)
// - configPrefix: "" (no prefix for environment variables)
//
// This function is used internally by the Build function to create a new
// options struct before applying the provided options.
func getDefaultOptions() *options {
	return &options{
		errHandler: nil,
		configPrefix: "",
	}
}

// WithErrorHandler returns an option that sets the error handler for the EzApp.
//
// The error handler is called when a Runnable returns an error that it couldn't handle.
// It should attempt to handle the error and return nil if it was handled successfully,
// or return an error if it couldn't be handled.
//
// If the error handler returns an error, the EzApp will cancel the context to initiate
// shutdown of all Runnables.
//
// Example:
//
//	app := ezapp.Build(
//		wireApp,
//		ezapp.WithErrorHandler(func(err error) error {
//			log.Printf("Error: %v", err)
//			return nil
//		}),
//	)
func WithErrorHandler(errHandler ErrHandler) option {
	return func(o *options) {
		o.errHandler = errHandler
	}
}


// WithConfigPrefix returns an option that sets the prefix for environment variables when loading configuration.
//
// The prefix is used when loading configuration from environment variables.
// For example, if the prefix is "APP", then environment variables will be prefixed with "APP_".
// This is useful when you have multiple applications running on the same host and want to
// avoid naming conflicts in environment variables.
//
// Example:
//
//	app := ezapp.Build(
//		wireApp,
//		ezapp.WithConfigPrefix("APP"),
//	)
//
// With this configuration, a struct field like:
//
//	type Config struct {
//		Port int `envconfig:"PORT" default:"8080"`
//	}
//
// Would be populated from the environment variable APP_PORT instead of PORT.
func WithConfigPrefix(prefix string) option {
	return func(o *options) {
		o.configPrefix = prefix
	}
}
