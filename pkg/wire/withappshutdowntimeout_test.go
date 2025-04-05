package wire

import (
	"testing"
	"time"

	"github.com/pgvanniekerk/ezapp/internal/conf"
)

// TestWithAppShutdownTimeout tests the WithAppShutdownTimeout function
func TestWithAppShutdownTimeout(t *testing.T) {
	// Create a time.Duration value
	timeout := 30 * time.Second

	// Call WithAppShutdownTimeout with this value
	option := WithAppShutdownTimeout(timeout)

	// Create an appOptions struct with an appConf field
	opts := &appOptions{
		appConf: conf.AppConf{
			ShutdownTimeout: 15 * time.Second, // Default value
		},
	}

	// Apply the option to the appOptions struct
	option(opts)

	// Check that the ShutdownTimeout field has been set correctly
	if opts.appConf.ShutdownTimeout != timeout {
		t.Errorf("Expected shutdown timeout to be %v, got %v", timeout, opts.appConf.ShutdownTimeout)
	}

	// Test with a different timeout
	timeout = 45 * time.Second
	option = WithAppShutdownTimeout(timeout)
	option(opts)

	// Check that the ShutdownTimeout field has been updated
	if opts.appConf.ShutdownTimeout != timeout {
		t.Errorf("Expected shutdown timeout to be %v, got %v", timeout, opts.appConf.ShutdownTimeout)
	}
}
