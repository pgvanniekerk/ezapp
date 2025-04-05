package ezapp

import (
	"context"
	"github.com/kelseyhightower/envconfig"
	"github.com/pgvanniekerk/ezapp/internal/app"
	"time"
)

// BuildOptions is an interface for configuring the Build function
type BuildOptions interface {
	GetErrorHandler() app.ErrorHandler
	GetStartupTimeout() time.Duration
	GetEnvVarPrefix() string
	GetShutdownSignal() <-chan struct{}
}

// Build creates an EzApp from a WireFunc with configuration
func Build[C any](wireFunc WireFunc[C], options BuildOptions) EzApp {

	// Panic if options is nil
	if options == nil {
		panic("options cannot be nil")
	}

	serviceSet := invokeWireFunc(wireFunc, options.GetStartupTimeout(), options.GetEnvVarPrefix())

	// Create a new App with the ServiceSet's services, error handler, and shutdown signal
	ezApp, err := app.NewApp(serviceSet.Services, options.GetErrorHandler(), options.GetShutdownSignal())
	if err != nil {
		panic(err)
	}
	return ezApp
}

func invokeWireFunc[C any](wireFunc WireFunc[C], startupTimeout time.Duration, envVarPrefix string) ServiceSet {

	// Create a context with the configured timeout
	ctx, cancel := context.WithTimeout(context.Background(), startupTimeout)
	defer cancel()

	// Create a new instance of the config struct
	var config C

	// Load environment variables into the config struct with the configured prefix
	if err := envconfig.Process(envVarPrefix, &config); err != nil {
		panic(err)
	}

	// Call the wire function to get a ServiceSet
	serviceSet, err := wireFunc(ctx, config)
	if err != nil {
		panic(err)
	}

	return serviceSet
}

// WireFunc is a function that returns a ServiceSet and an error
type WireFunc[C any] func(context.Context, C) (ServiceSet, error)
