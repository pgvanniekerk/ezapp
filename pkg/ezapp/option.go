package ezapp

// option is a function type that configures an options struct.
//
// An option is used to configure the behavior of the EzApp when calling Build.
// It follows the functional options pattern, which allows for a clean and
// extensible API for configuring the EzApp.
//
// Options are provided to the Build function as variadic arguments:
//
//	app := ezapp.Build(
//		wireApp,
//		ezapp.WithErrorHandler(handleError),
//		ezapp.WithCleanupFunc(cleanup),
//		ezapp.WithConfigPrefix("APP"),
//	)
//
// The available options are:
// - WithErrorHandler: Sets the error handler for the EzApp
// - WithCleanupFunc: Sets the cleanup function for the EzApp
// - WithConfigPrefix: Sets the prefix for environment variables when loading configuration
//
// Custom options can be created by implementing a function that takes an *options
// struct and modifies it.
type option func(*options)
