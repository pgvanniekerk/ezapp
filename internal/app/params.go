package app

import (
	"log/slog"
	"time"
)

// Params holds the parameters for creating a new App instance.
// This struct is used by the app.New function to configure the App.
//
// The Params struct is typically created by the wire.App function
// based on the options provided by the user. It is not meant to be
// created directly by users of the ezapp framework.
type Params struct {
	// ShutdownTimeout is the maximum time allowed for stopping all runnables
	// during application shutdown. If this timeout is reached, any remaining
	// runnables will be forcibly terminated.
	//
	// Default: 15 seconds (set by wire.defaultOptions)
	ShutdownTimeout time.Duration

	// Runnables is a slice of components that implement the Runnable interface.
	// These components will be managed by the App, which will start them when
	// the application starts and stop them when the application shuts down.
	//
	// At least one Runnable is required for the App to be useful.
	Runnables []Runnable

	// ShutdownSig is a channel that signals when the application should shut down.
	// When a value is received on this channel, the App will initiate graceful
	// shutdown of all runnables.
	//
	// Default: A channel that receives a signal when os.Interrupt is received
	// (set by wire.defaultOptions)
	ShutdownSig <-chan error

	// Logger is used for application-level logging.
	// This logger will be used by the App to log startup, shutdown, and error events.
	//
	// Default: A logger that writes to os.Stderr with INFO level
	// (set by wire.defaultOptions)
	Logger *slog.Logger

	// LogAttrs is a slice of log attributes to be added to the logger.
	// These attributes will be included in all log entries created by the App.
	//
	// Default: Empty slice (set by wire.defaultOptions)
	LogAttrs []slog.Attr
}
