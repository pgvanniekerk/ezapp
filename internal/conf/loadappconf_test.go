package conf

import (
	"os"
	"testing"
	"time"
)

// TestLoadAppConf tests the LoadAppConf function
func TestLoadAppConf(t *testing.T) {
	// Save the original environment variables
	originalShutdownTimeout := os.Getenv("EZAPP_SHUTDOWN_TIMEOUT")
	originalStartupTimeout := os.Getenv("EZAPP_STARTUP_TIMEOUT")

	// Restore the original environment variables when the test is done
	defer func() {
		os.Setenv("EZAPP_SHUTDOWN_TIMEOUT", originalShutdownTimeout)
		os.Setenv("EZAPP_STARTUP_TIMEOUT", originalStartupTimeout)
	}()

	// Test with custom environment variables
	os.Setenv("EZAPP_SHUTDOWN_TIMEOUT", "30s")
	os.Setenv("EZAPP_STARTUP_TIMEOUT", "20s")

	// Call LoadAppConf
	conf, err := LoadAppConf()

	// Check that there was no error
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check that the configuration has the expected values
	expectedShutdownTimeout := 30 * time.Second
	if conf.ShutdownTimeout != expectedShutdownTimeout {
		t.Errorf("Expected ShutdownTimeout to be %v, got %v", expectedShutdownTimeout, conf.ShutdownTimeout)
	}

	expectedStartupTimeout := 20 * time.Second
	if conf.StartupTimeout != expectedStartupTimeout {
		t.Errorf("Expected StartupTimeout to be %v, got %v", expectedStartupTimeout, conf.StartupTimeout)
	}

	// Test with default values
	os.Unsetenv("EZAPP_SHUTDOWN_TIMEOUT")
	os.Unsetenv("EZAPP_STARTUP_TIMEOUT")

	// Call LoadAppConf
	conf, err = LoadAppConf()

	// Check that there was no error
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check that the configuration has the default values
	expectedShutdownTimeout = 15 * time.Second
	if conf.ShutdownTimeout != expectedShutdownTimeout {
		t.Errorf("Expected ShutdownTimeout to be %v, got %v", expectedShutdownTimeout, conf.ShutdownTimeout)
	}

	expectedStartupTimeout = 15 * time.Second
	if conf.StartupTimeout != expectedStartupTimeout {
		t.Errorf("Expected StartupTimeout to be %v, got %v", expectedStartupTimeout, conf.StartupTimeout)
	}

	// Test with invalid values
	os.Setenv("EZAPP_SHUTDOWN_TIMEOUT", "invalid")

	// Call LoadAppConf
	_, err = LoadAppConf()

	// Check that there was an error
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
