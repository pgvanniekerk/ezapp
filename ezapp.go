package ezapp

import (
	"context"
	"github.com/pgvanniekerk/ezapp/internal/app"
	"github.com/pgvanniekerk/ezapp/internal/config"
	"go.uber.org/zap"
)

// InitCtx provides the initialization context passed to an Initializer function.
// It contains all the resources needed during application startup including
// configuration, logging, and a startup context with timeout.
//
// The generic type parameter Config should be a struct with appropriate
// environment variable tags for automatic configuration loading.
type InitCtx[Config any] struct {

	// StartupCtx is a context with a configurable timeout (default 15 seconds)
	// that can be used during initialization to enforce startup time limits.
	// The timeout is controlled by the EZAPP_STARTUP_TIMEOUT environment variable.
	StartupCtx context.Context

	// Logger is a configured zap.Logger instance ready for use.
	// The log level is controlled by the EZAPP_LOG_LEVEL environment variable
	// (default: INFO). Supports DEBUG, INFO, WARN, ERROR, DPANIC, PANIC, FATAL.
	Logger *zap.Logger

	// Config contains the application configuration loaded from environment variables
	// using the Netflix go-env package. The Config type should be a struct with
	// appropriate `env` tags for field mapping.
	Config Config
}

// AppCtx represents the application context containing all the runners
// that will be executed concurrently by the application framework.
// This is typically constructed using the Construct function with functional options.
type AppCtx struct {
	runnerList  []app.Runner
	cleanupFunc func(shutdownCtx context.Context) error
}

// Initializer is a function type that takes an InitCtx and returns an AppCtx.
// This is the main entry point for application initialization logic where
// dependencies are wired together and runners are configured.
//
// The initializer function should:
//   - Use the provided logger for any initialization logging
//   - Use the provided config for application configuration
//   - Respect the StartupCtx timeout for any initialization operations
//   - Return an AppCtx containing all the runners to be executed
//
// Example:
//
//	func MyInitializer(ctx InitCtx[MyConfig]) (AppCtx, error) {
//	    server := NewServer(ctx.Config.Port, ctx.Logger)
//	    return Construct(WithRunners(server.Run))
//	}
type Initializer[Config any] func(InitCtx[Config]) (AppCtx, error)

// option represents a functional option for configuring an AppCtx.
// This type is not exported to ensure only predefined options can be used.
type option func(*AppCtx) error

// WithRunners is a functional option that adds a list of runners to the AppCtx.
// Runners are functions that implement the application's core business logic
// and will be executed concurrently by the framework.
//
// Each runner function receives a context.Context that will be cancelled
// when the application needs to shut down (e.g., on SIGINT/SIGTERM).
// Runners should monitor this context and perform graceful cleanup when cancelled.
//
// Example:
//
//	func serverRunner(ctx context.Context) error {
//	    server := &http.Server{Addr: ":8080"}
//	    go func() {
//	        <-ctx.Done()
//	        server.Shutdown(context.Background())
//	    }()
//	    return server.ListenAndServe()
//	}
//
//	appCtx, err := Construct(WithRunners(serverRunner, anotherRunner))
func WithRunners(runners ...app.Runner) option {
	return func(appCtx *AppCtx) error {
		appCtx.runnerList = append(appCtx.runnerList, runners...)
		return nil
	}
}

// WithCleanup is a functional option that sets a cleanup function for the AppCtx.
// The cleanup function is called after all runners have completed, allowing for
// graceful cleanup of resources like database connections, file handles, etc.
//
// The cleanup function receives a context with a timeout (controlled by the
// EZAPP_SHUTDOWN_TIMEOUT environment variable, default 15 seconds) to enforce
// cleanup time limits and prevent hanging during application shutdown.
//
// Example:
//
//	func cleanup(ctx context.Context) error {
//	    // Close database connections
//	    if err := db.Close(); err != nil {
//	        return fmt.Errorf("failed to close database: %w", err)
//	    }
//	    // Close other resources...
//	    return nil
//	}
//
//	appCtx, err := Construct(
//	    WithRunners(server.Run),
//	    WithCleanup(cleanup),
//	)
func WithCleanup(cleanupFunc func(shutdownCtx context.Context) error) option {
	return func(appCtx *AppCtx) error {
		appCtx.cleanupFunc = cleanupFunc
		return nil
	}
}

// Construct builds an AppCtx using the provided functional options.
// This is the primary way to configure an application context with runners
// and other configuration options.
//
// The function applies each option in order, allowing for flexible
// composition of application components. If any option returns an error,
// construction is aborted and the error is returned.
//
// Example:
//
//	appCtx, err := Construct(
//	    WithRunners(server.Run, worker.Run),
//	    // Future options like WithMiddleware, WithHealthCheck, etc.
//	)
//	if err != nil {
//	    return err
//	}
func Construct(options ...option) (AppCtx, error) {

	appCtx := AppCtx{
		runnerList:  make([]app.Runner, 0, 8),
		cleanupFunc: nil,
	}

	for _, opt := range options {
		if err := opt(&appCtx); err != nil {
			return AppCtx{}, err
		}
	}

	return appCtx, nil
}

// Run is the main entry point for starting an EzApp application.
// It orchestrates the complete application lifecycle and takes full control
// of the application execution:
//
// 1. Loads configuration from environment variables using the provided Config type
// 2. Initializes a structured logger with configurable log levels
// 3. Creates a startup context with configurable timeout
// 4. Invokes the provided initializer function to build the application
// 5. Runs all configured runners concurrently with graceful shutdown
// 6. Performs cleanup operations after all runners complete
//
// This function does not return - it handles all error cases by logging
// and calling logger.Fatal() to terminate the application. It will block
// until all runners complete successfully or an error occurs.
//
// Environment Variables:
//   - EZAPP_LOG_LEVEL: Controls logging verbosity (DEBUG, INFO, WARN, ERROR, etc.)
//   - EZAPP_STARTUP_TIMEOUT: Timeout in seconds for initialization (default: 15)
//   - EZAPP_SHUTDOWN_TIMEOUT: Timeout in seconds for graceful shutdown (default: 15)
//   - Plus any variables defined in your Config struct
//
// Example:
//
//	type MyConfig struct {
//	    Port int `env:"PORT" default:"8080"`
//	    DatabaseURL string `env:"DATABASE_URL" required:"true"`
//	}
//
//	func main() {
//	    ezapp.Run(func(ctx ezapp.InitCtx[MyConfig]) (ezapp.AppCtx, error) {
//	        server := NewServer(ctx.Config.Port, ctx.Logger)
//	        return ezapp.Construct(ezapp.WithRunners(server.Run))
//	    })
//	    // This point is never reached - Run() handles application lifecycle
//	}
func Run[Config any](initializer Initializer[Config]) {

	// Load logger
	logger := config.LoadLogger()

	// Load configuration from environment variables
	cfg, err := config.LoadVar[Config]()
	if err != nil {
		logger.Fatal("failed to load configuration", zap.Error(err))
	}

	// Create startup context with timeout
	startupCtx, err := config.StartupCtx()
	if err != nil {
		logger.Fatal("failed to create startup context", zap.Error(err))
	}

	// Create initialization context
	initCtx := InitCtx[Config]{
		StartupCtx: startupCtx,
		Logger:     logger,
		Config:     cfg,
	}

	// Invoke the initializer to get the app context
	appCtx, err := initializer(initCtx)
	if err != nil {
		logger.Fatal("initialization failed", zap.Error(err))
	}

	// Create and run the app
	application := app.New(appCtx.runnerList, logger)
	appErr := application.Run()

	// After app completes, run cleanup if provided
	if appCtx.cleanupFunc != nil {

		// Create a shutdown context with the configured timeout
		shutdownCtx, err := config.ShutdownCtx()
		if err != nil {
			logger.Fatal("failed to create shutdown context", zap.Error(err))
		}

		// Run cleanup function
		if cleanupErr := appCtx.cleanupFunc(shutdownCtx); cleanupErr != nil {
			logger.Error("cleanup failed", zap.Error(cleanupErr))
			// If the app ran successfully but cleanup failed, fatal exit
			if appErr == nil {
				logger.Fatal("application cleanup failed", zap.Error(cleanupErr))
			}
			// If both app and cleanup failed, fatal exit with app error (more critical)
		}
	}

	// If the app failed, fatal exit
	if appErr != nil {
		logger.Fatal("application failed", zap.Error(appErr))
	}

	// Application completed successfully
	logger.Info("application completed successfully")
}
