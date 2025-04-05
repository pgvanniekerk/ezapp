package wire

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/pgvanniekerk/ezapp/internal/app"
)

// TestApp tests the App function
func TestApp(t *testing.T) {
	// Create a mock runnable
	mockRunnable := &MockRunnable{}

	// Create a function that returns a slice of runnables
	runnablesFunc := func() []app.Runnable {
		return []app.Runnable{mockRunnable}
	}

	// Create some options
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	shutdownTimeout := 30 * time.Second
	startupTimeout := 20 * time.Second
	errChan := make(chan error, 1)
	logAttr := slog.String("key", "value")

	// Call App with these parameters
	appInstance, err := App(
		runnablesFunc,
		WithLogger(logger),
		WithAppShutdownTimeout(shutdownTimeout),
		WithAppStartupTimeout(startupTimeout),
		WithShutdownSignal(errChan),
		WithLogAttrs(logAttr),
	)

	// Check that there was no error
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check that the app instance is not nil
	if appInstance == nil {
		t.Errorf("Expected non-nil app instance, got nil")
	}

	// Test with no options
	appInstance, err = App(runnablesFunc)

	// Check that there was no error
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check that the app instance is not nil
	if appInstance == nil {
		t.Errorf("Expected non-nil app instance, got nil")
	}

	// Test with a nil runnablesFunc
	_, err = App(nil)

	// Check that there was an error
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// Check that the error message is as expected
	expectedErrMsg := "runnablesFunc cannot be nil"
	if err.Error() != expectedErrMsg {
		t.Errorf("Expected error message %q, got %q", expectedErrMsg, err.Error())
	}
}

// TestAppDefaultOptionsError tests the App function when defaultOptions returns an error
func TestAppDefaultOptionsError(t *testing.T) {
	// Save the original environment variable
	originalValue := os.Getenv("EZAPP_SHUTDOWN_TIMEOUT")

	// Restore the original environment variable when the test is done
	t.Cleanup(func() {
		os.Setenv("EZAPP_SHUTDOWN_TIMEOUT", originalValue)
	})

	// Set the environment variable to an invalid value
	os.Setenv("EZAPP_SHUTDOWN_TIMEOUT", "invalid")

	// Create a function that returns a slice of runnables
	runnablesFunc := func() []app.Runnable {
		return []app.Runnable{&MockRunnable{}}
	}

	// Call App
	_, err := App(runnablesFunc)

	// Check that there was an error
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// Check that the error message contains the expected text
	expectedErrMsg := "failed to retrieve default options for app:"
	errMsg := err.Error()
	if len(errMsg) < len(expectedErrMsg) || errMsg[:len(expectedErrMsg)] != expectedErrMsg {
		t.Errorf("Expected error message to start with %q, got %q", expectedErrMsg, errMsg)
	}
}
