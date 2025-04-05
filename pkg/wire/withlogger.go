package wire

import (
	"log/slog"
)

// WithLogger returns an AppOption that sets the logger for the App
func WithLogger(logger *slog.Logger) AppOption {
	return func(o *appOptions) {
		o.logger = logger
	}
}
