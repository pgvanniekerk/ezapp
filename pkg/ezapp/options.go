package ezapp

import "log/slog"

type options struct {
	logger *slog.Logger
}

func WithLogger(logger *slog.Logger) func(*options) {
	return func(o *options) {
		o.logger = logger
	}
}
