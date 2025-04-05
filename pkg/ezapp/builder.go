package ezapp

import (
	"context"
)

// Builder is a generic function type that defines the signature for application builder functions.
// It takes a context and a configuration object, and returns an application instance and an error.
//
// The type parameters are:
//   - APP: A type that implements the EzApp interface
//   - CONF: Any type that represents the application configuration
//
// Builder functions are used with ezapp.Run to create and run applications:
//
// Example:
//
//	type Config struct {
//	    LogLevel string `envvar:"LOG_LEVEL" default:"info"`
//	    Port int `envvar:"PORT" default:"8080"`
//	}
//
//	func BuildFunc(ctx context.Context, cfg Config) (ezapp.EzApp, error) {
//	    // Create and configure your application
//	    return wire.App(
//	        wire.Runnables(myRunnable),
//	        wire.WithAppShutdownTimeout(15*time.Second),
//	    )
//	}
//
//	func main() {
//	    ezapp.Run(BuildFunc)
//	}
//
// Potential errors:
//   - Configuration validation errors
//   - Resource initialization errors (database connections, etc.)
//   - Invalid runnable components
type Builder[APP EzApp, CONF any] func(context.Context, CONF) (APP, error)
