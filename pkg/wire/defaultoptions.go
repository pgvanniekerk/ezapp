package wire

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/pgvanniekerk/ezapp/internal/conf"
)

// defaultOptions returns the default options for the App function
func defaultOptions() (*appOptions, error) {

	appConf, err := conf.LoadAppConf()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve app configuration: %w", err)
	}

	// Create a default logger that writes to stdout with INFO level
	opts := &slog.HandlerOptions{
		Level: slog.LevelError,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))

	// Create a nil channel for shutdownSig
	var shutdownSig <-chan error

	return &appOptions{
		appConf:     appConf,
		shutdownSig: shutdownSig,
		logger:      logger,
		logAttrs:    []slog.Attr{},
	}, nil
}
