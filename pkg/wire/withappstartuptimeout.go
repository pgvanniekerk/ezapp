package wire

import (
	"time"
)

// WithAppStartupTimeout returns an AppOption that sets the startup timeout for the App
func WithAppStartupTimeout(timeout time.Duration) AppOption {
	return func(o *appOptions) {
		o.appConf.StartupTimeout = timeout
	}
}
