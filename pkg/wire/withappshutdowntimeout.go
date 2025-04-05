package wire

import (
	"time"
)

// WithAppShutdownTimeout returns an AppOption that sets the shutdown timeout for the App
func WithAppShutdownTimeout(timeout time.Duration) AppOption {
	return func(o *appOptions) {
		o.appConf.ShutdownTimeout = timeout
	}
}
