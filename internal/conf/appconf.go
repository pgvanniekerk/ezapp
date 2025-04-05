package conf

import (
	"time"
)

// AppConf holds the application configuration for the ezapp framework.
// This struct defines configuration values that control the behavior of
// the application, such as timeouts for startup and shutdown.
//
// The configuration values can be set through environment variables using
// the envconfig tags. The EZAPP prefix is added to the environment variable
// names by the LoadAppConf function.
type AppConf struct {
	// ShutdownTimeout is the maximum time allowed for stopping all runnables
	// during application shutdown. If this timeout is reached, any remaining
	// runnables will be forcibly terminated.
	//
	// Environment variable: EZAPP_SHUTDOWN_TIMEOUT
	// Default: 15 seconds
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"15s"`

	// StartupTimeout is the maximum time allowed for starting all runnables
	// during application startup. This timeout is not currently used by the
	// framework but is reserved for future use.
	//
	// Environment variable: EZAPP_STARTUP_TIMEOUT
	// Default: 15 seconds
	StartupTimeout time.Duration `envconfig:"STARTUP_TIMEOUT" default:"15s"`
}
