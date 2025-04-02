package buildoption

import (
	"github.com/pgvanniekerk/ezapp/internal/app"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// BuildOptions is an interface that defines methods for configuring a build
type BuildOptions interface {
	// GetErrorHandler returns the error handler configured for the build
	GetErrorHandler() app.ErrorHandler
	// GetStartupTimeout returns the timeout for the startup context
	GetStartupTimeout() time.Duration
	// GetEnvVarPrefix returns the prefix for environment variables
	GetEnvVarPrefix() string
	// GetShutdownSignal returns the channel used for shutdown signaling
	GetShutdownSignal() <-chan struct{}
}

// Options holds configuration options for the Build function
type Options struct {
	ErrorHandler   app.ErrorHandler
	StartupTimeout time.Duration
	EnvVarPrefix   string
	ShutdownSignal <-chan struct{}
}

// DefaultStartupTimeout is the default timeout for the startup context
const DefaultStartupTimeout = 15 * time.Second

// DefaultEnvVarPrefix is the default prefix for environment variables
const DefaultEnvVarPrefix = ""

// WithoutOptions creates a new Options with default configuration
func WithoutOptions() BuildOptions {
	return &Options{
		ErrorHandler:   DefaultErrorHandler,
		StartupTimeout: DefaultStartupTimeout,
		EnvVarPrefix:   DefaultEnvVarPrefix,
		ShutdownSignal: nil, // Will use DefaultShutdownSignal when GetShutdownSignal is called
	}
}

// WithOptions creates a new Options with default configuration and applies the given options
func WithOptions(opts ...Option) BuildOptions {
	options := &Options{
		ErrorHandler:   DefaultErrorHandler,
		StartupTimeout: DefaultStartupTimeout,
		EnvVarPrefix:   DefaultEnvVarPrefix,
		ShutdownSignal: nil, // Will use DefaultShutdownSignal when GetShutdownSignal is called
	}

	// Apply options
	for _, opt := range opts {
		opt(options)
	}

	return options
}

// GetErrorHandler implements the BuildOptions interface
func (o *Options) GetErrorHandler() app.ErrorHandler {
	return o.ErrorHandler
}

// GetStartupTimeout implements the BuildOptions interface
func (o *Options) GetStartupTimeout() time.Duration {
	return o.StartupTimeout
}

// GetEnvVarPrefix implements the BuildOptions interface
func (o *Options) GetEnvVarPrefix() string {
	return o.EnvVarPrefix
}

// GetShutdownSignal implements the BuildOptions interface
func (o *Options) GetShutdownSignal() <-chan struct{} {
	// If no shutdown signal is provided, create a default one
	if o.ShutdownSignal == nil {
		return DefaultShutdownSignal()
	}
	return o.ShutdownSignal
}

// DefaultShutdownSignal creates a channel that closes when SIGTERM or SIGINT is received
func DefaultShutdownSignal() <-chan struct{} {
	// Create a channel for SIGTERM (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	// Create a channel for shutdown signal
	shutdownChan := make(chan struct{})

	// Convert os.Signal channel to struct{} channel
	go func() {
		<-sigChan
		close(shutdownChan)
	}()

	return shutdownChan
}

// Option is a function that configures Options
type Option func(*Options)

// WithErrorHandler sets a custom error handler for the application
func WithErrorHandler(handler app.ErrorHandler) Option {
	return func(options *Options) {
		options.ErrorHandler = handler
	}
}

// WithStartupTimeout sets a custom timeout for the startup context
func WithStartupTimeout(timeout time.Duration) Option {
	return func(options *Options) {
		options.StartupTimeout = timeout
	}
}

// WithEnvVarPrefix sets a custom prefix for environment variables
func WithEnvVarPrefix(prefix string) Option {
	return func(options *Options) {
		options.EnvVarPrefix = prefix
	}
}

// WithShutdownSignal sets a custom shutdown signal channel
func WithShutdownSignal(shutdownSignal <-chan struct{}) Option {
	return func(options *Options) {
		options.ShutdownSignal = shutdownSignal
	}
}

// DefaultErrorHandler is the default error handler that panics on errors
func DefaultErrorHandler(err error) error {
	panic(err)
}
