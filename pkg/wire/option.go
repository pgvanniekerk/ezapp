package wire

import (
	"log/slog"

	"github.com/pgvanniekerk/ezapp/internal/conf"
)

// appOptions holds the configuration options for the App function.
// This struct is used internally by the wire package to store configuration
// that is applied to the application during creation.
type appOptions struct {
	// appConf contains application configuration like timeouts
	appConf conf.AppConf

	// shutdownSig is a channel that signals when the application should shut down
	shutdownSig <-chan error

	// logger is the logger used by the application
	logger *slog.Logger

	// logAttrs are additional attributes to add to log entries
	logAttrs []slog.Attr
}

// AppOption is a function that configures the App function.
// This type implements the functional options pattern, allowing for flexible
// and extensible configuration of applications.
//
// The ezapp framework provides several built-in options:
//   - WithAppShutdownTimeout: Sets the timeout for application shutdown
//   - WithAppStartupTimeout: Sets the timeout for application startup
//   - WithLogger: Sets the logger for the application
//   - WithLogAttrs: Adds attributes to log entries
//   - WithShutdownSignal: Sets the channel for receiving shutdown signals
//
// Example:
//
//	app, err := wire.App(
//	    wire.Runnables(myRunnable),
//	    wire.WithAppShutdownTimeout(15*time.Second),
//	    wire.WithLogger(slog.New(slog.NewTextHandler(os.Stdout, nil))),
//	)
type AppOption func(*appOptions)
