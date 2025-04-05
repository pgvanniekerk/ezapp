package wire

import (
	"testing"
	"time"

	"github.com/pgvanniekerk/ezapp/internal/conf"
)

// TestWithAppStartupTimeout tests the WithAppStartupTimeout function
func TestWithAppStartupTimeout(t *testing.T) {
	// Create a time.Duration value
	timeout := 30 * time.Second

	// Call WithAppStartupTimeout with this value
	option := WithAppStartupTimeout(timeout)

	// Create an appOptions struct with an appConf field
	opts := &appOptions{
		appConf: conf.AppConf{
			StartupTimeout: 15 * time.Second, // Default value
		},
	}

	// Apply the option to the appOptions struct
	option(opts)

	// Check that the StartupTimeout field has been set correctly
	if opts.appConf.StartupTimeout != timeout {
		t.Errorf("Expected startup timeout to be %v, got %v", timeout, opts.appConf.StartupTimeout)
	}

	// Test with a different timeout
	timeout = 45 * time.Second
	option = WithAppStartupTimeout(timeout)
	option(opts)

	// Check that the StartupTimeout field has been updated
	if opts.appConf.StartupTimeout != timeout {
		t.Errorf("Expected startup timeout to be %v, got %v", timeout, opts.appConf.StartupTimeout)
	}
}
