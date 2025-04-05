package wire

import (
	"log/slog"
)

// WithLogAttrs returns an AppOption that sets the default log attributes for the App
func WithLogAttrs(attrs ...slog.Attr) AppOption {
	return func(o *appOptions) {
		o.logAttrs = attrs
	}
}
