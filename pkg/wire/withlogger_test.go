package wire

import (
	"log/slog"
	"os"
	"testing"
)

// TestWithLogger tests the WithLogger function
func TestWithLogger(t *testing.T) {
	// Create a logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Call WithLogger with this logger
	option := WithLogger(logger)

	// Create an appOptions struct
	opts := &appOptions{
		logger: nil, // Initialize with nil
	}

	// Apply the option to the appOptions struct
	option(opts)

	// Check that the logger field has been set correctly
	if opts.logger != logger {
		t.Errorf("Expected logger to be %v, got %v", logger, opts.logger)
	}

	// Test with a different logger
	logger2 := slog.New(slog.NewTextHandler(os.Stderr, nil))
	option = WithLogger(logger2)
	option(opts)

	// Check that the logger field has been updated
	if opts.logger != logger2 {
		t.Errorf("Expected logger to be %v, got %v", logger2, opts.logger)
	}
}
