package wire

import (
	"log/slog"
	"os"
	"testing"
)

// TestDefaultOptions tests the defaultOptions function
func TestDefaultOptions(t *testing.T) {
	// Save the original environment variables
	originalShutdownTimeout := os.Getenv("EZAPP_SHUTDOWN_TIMEOUT")
	originalStartupTimeout := os.Getenv("EZAPP_STARTUP_TIMEOUT")

	// Restore the original environment variables when the test is done
	t.Cleanup(func() {
		os.Setenv("EZAPP_SHUTDOWN_TIMEOUT", originalShutdownTimeout)
		os.Setenv("EZAPP_STARTUP_TIMEOUT", originalStartupTimeout)
	})

	// Unset the environment variables to ensure we're testing with default values
	os.Unsetenv("EZAPP_SHUTDOWN_TIMEOUT")
	os.Unsetenv("EZAPP_STARTUP_TIMEOUT")
	// Call defaultOptions
	options, err := defaultOptions()

	// Check that there was no error
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check that the options are not nil
	if options == nil {
		t.Errorf("Expected non-nil options, got nil")
	}

	// Check that the appConf field is not nil
	if options.appConf.ShutdownTimeout == 0 {
		t.Errorf("Expected non-zero ShutdownTimeout, got 0")
	}

	if options.appConf.StartupTimeout == 0 {
		t.Errorf("Expected non-zero StartupTimeout, got 0")
	}

	// Check that the logger field is not nil
	if options.logger == nil {
		t.Errorf("Expected non-nil logger, got nil")
	}

	// Check that the logAttrs field is initialized to an empty slice
	if options.logAttrs == nil {
		t.Errorf("Expected non-nil logAttrs, got nil")
	}

	if len(options.logAttrs) != 0 {
		t.Errorf("Expected empty logAttrs, got %v", options.logAttrs)
	}

	// Check that the shutdownSig field is nil
	if options.shutdownSig != nil {
		t.Errorf("Expected nil shutdownSig, got %v", options.shutdownSig)
	}

	// Test that the logger has the expected level
	handler, ok := options.logger.Handler().(*slog.TextHandler)
	if !ok {
		t.Errorf("Expected logger handler to be *slog.TextHandler, got %T", options.logger.Handler())
	} else {
		// Note: In a real test, you might want to check the level of the handler,
		// but this requires accessing unexported fields, which is not recommended.
		// Instead, we just check that the handler is of the expected type.
		_ = handler
	}
}
