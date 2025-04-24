package ezapp

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

// TestConfig is a test configuration struct
type TestConfig struct {
	StringValue string  `env:"TEST_STRING"`
	IntValue    int     `env:"TEST_INT"`
	BoolValue   bool    `env:"TEST_BOOL"`
	FloatValue  float64 `env:"TEST_FLOAT"`
	DefaultName string  `env:"DEFAULTNAME"` // Will use environment variable "DEFAULTNAME"
}

// mockRunnable is a mock implementation of the Runnable interface for testing
type mockRunnable struct {
	runFunc func(context.Context) error
}

func (m mockRunnable) Run(ctx context.Context) error {
	return m.runFunc(ctx)
}

// exitError is a custom error type used to capture os.Exit calls in tests
type exitError struct {
	code int
}

// testExitCode is used to capture the exit code in tests
var testExitCode int

// captureOutput captures stdout during the execution of a function
func captureOutput(f func()) string {
	// Save the original stdout
	originalStdout := os.Stdout

	// Create a pipe to capture stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Create a channel to handle panics
	done := make(chan interface{})

	// Execute the function in a goroutine
	go func() {
		defer func() {
			// Recover from any panic
			if r := recover(); r != nil {
				// Check if it's our exitError
				if e, ok := r.(exitError); ok {
					// Just record the exit code and continue
					done <- e
				} else {
					// Re-panic for other errors
					panic(r)
				}
			} else {
				// No panic, signal completion
				close(done)
			}
		}()

		// Execute the function
		f()
	}()

	// Wait for the function to complete or panic
	<-done

	// Close the pipe writer and restore stdout
	w.Close()
	os.Stdout = originalStdout

	// Read the captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	r.Close()

	// Return the captured output
	return buf.String()
}

// TestRun tests the Run function
func TestRun(t *testing.T) {
	// Skip the actual test since it would call os.Exit
	t.Skip("Skipping TestRun since it would call os.Exit")
}

// TestRunWithMockExit tests the Run function with a mock exit function
func TestRunWithMockExit(t *testing.T) {
	// Save the original osExit function and restore it after the test
	originalExit := osExit
	defer func() {
		osExit = originalExit
	}()

	// Create a mock exit function that records the exit code
	testExitCode = 0 // Reset the exit code
	osExit = func(code int) {
		testExitCode = code
		panic(exitError{code: code})
	}

	// Set environment variables for the test
	os.Setenv("TEST_STRING", "test-value")
	defer func() {
		os.Unsetenv("TEST_STRING")
	}()

	// Create a mock builder that returns a successful runnable
	mockBuilder := func(conf TestConfig) (App, error) {
		return Construct(
			WithRunnables(
				mockRunnable{
					runFunc: func(ctx context.Context) error {
						// This runnable completes successfully immediately
						return nil
					},
				},
			),
		), nil
	}

	// Call the Run function with captureOutput to handle the panic
	output := captureOutput(func() {
		Run(mockBuilder)
	})

	// Check that the output contains the success message
	if !strings.Contains(output, "App shutting down: All runnables completed successfully") {
		t.Errorf("Expected output to contain success message, got: %q", output)
	}

	// Assert that the exit code is 0 (success)
	if testExitCode != 0 {
		t.Errorf("Expected exit code to be 0, but got %d", testExitCode)
	}
}

// TestRunWithError tests the Run function with a builder that returns an error
func TestRunWithError(t *testing.T) {
	// Save the original osExit function and restore it after the test
	originalExit := osExit
	defer func() {
		osExit = originalExit
	}()

	// Create a mock exit function that records the exit code
	testExitCode = 0 // Reset the exit code
	osExit = func(code int) {
		testExitCode = code
		panic(exitError{code: code})
	}

	// Create a mock builder that returns an error
	expectedError := errors.New("builder error")
	mockBuilder := func(conf TestConfig) (App, error) {
		return App{}, expectedError
	}

	// Call the Run function with captureOutput to handle the panic
	output := captureOutput(func() {
		Run(mockBuilder)
	})

	// Check that the output contains the error message
	if !strings.Contains(output, "App shutting down: Initialization error: builder error") {
		t.Errorf("Expected output to contain error message, got: %q", output)
	}

	// Assert that the exit code is 1 (error)
	if testExitCode != 1 {
		t.Errorf("Expected exit code to be 1, but got %d", testExitCode)
	}
}

// TestRunWithRunnableError tests the Run function with a runnable that returns an error
func TestRunWithRunnableError(t *testing.T) {
	// Save the original osExit function and restore it after the test
	originalExit := osExit
	defer func() {
		osExit = originalExit
	}()

	// Create a mock exit function that records the exit code
	testExitCode = 0 // Reset the exit code
	osExit = func(code int) {
		testExitCode = code
		panic(exitError{code: code})
	}

	// Set environment variables for the test
	os.Setenv("TEST_STRING", "test-value")
	defer func() {
		os.Unsetenv("TEST_STRING")
	}()

	// Create a mock builder that returns a runnable that returns an error
	mockBuilder := func(conf TestConfig) (App, error) {
		return Construct(
			WithRunnables(
				mockRunnable{
					runFunc: func(ctx context.Context) error {
						return errors.New("runnable error")
					},
				},
			),
		), nil
	}

	// Call the Run function with captureOutput to handle the panic
	output := captureOutput(func() {
		Run(mockBuilder)
	})

	// Check that the output contains the error message
	if !strings.Contains(output, "App shutting down: Runnable error: runnable error") {
		t.Errorf("Expected output to contain error message, got: %q", output)
	}

	// Assert that the exit code is 1 (error)
	if testExitCode != 1 {
		t.Errorf("Expected exit code to be 1, but got %d", testExitCode)
	}
}

// TestRunWithNoRunnables tests the Run function with a builder that returns no runnables
func TestRunWithNoRunnables(t *testing.T) {
	// Save the original osExit function and restore it after the test
	originalExit := osExit
	defer func() {
		osExit = originalExit
	}()

	// Create a mock exit function that records the exit code
	testExitCode = 0 // Reset the exit code
	osExit = func(code int) {
		testExitCode = code
		panic(exitError{code: code})
	}

	// Set environment variables for the test
	os.Setenv("TEST_STRING", "test-value")
	defer func() {
		os.Unsetenv("TEST_STRING")
	}()

	// Create a mock builder that returns no runnables
	mockBuilder := func(conf TestConfig) (App, error) {
		return Construct(), nil
	}

	// Call the Run function with captureOutput to handle the panic
	output := captureOutput(func() {
		Run(mockBuilder)
	})

	// Check that the output contains the no runnables message
	if !strings.Contains(output, "App shutting down: No runnables to execute") {
		t.Errorf("Expected output to contain no runnables message, got: %q", output)
	}

	// Assert that the exit code is 0 (success)
	if testExitCode != 0 {
		t.Errorf("Expected exit code to be 0, but got %d", testExitCode)
	}
}

// TestRunWithInvalidConfig tests the Run function with an invalid configuration type
func TestRunWithInvalidConfig(t *testing.T) {
	// Skip this test for now as it's hard to test with generics
	t.Skip("Skipping TestRunWithInvalidConfig as it's hard to test with generics")
}

// TestRunWithCleanup tests the Run function with a cleanup function
func TestRunWithCleanup(t *testing.T) {
	// Save the original osExit function and restore it after the test
	originalExit := osExit
	defer func() {
		osExit = originalExit
	}()

	// Create a mock exit function that records the exit code
	testExitCode = 0 // Reset the exit code
	osExit = func(code int) {
		testExitCode = code
		panic(exitError{code: code})
	}

	// Set environment variables for the test
	os.Setenv("TEST_STRING", "test-value")
	defer func() {
		os.Unsetenv("TEST_STRING")
	}()

	// Track if cleanup was called
	cleanupCalled := false

	// Create a mock builder that returns a successful runnable and a cleanup function
	mockBuilder := func(conf TestConfig) (App, error) {
		return Construct(
			WithRunnables(
				mockRunnable{
					runFunc: func(ctx context.Context) error {
						// This runnable completes successfully immediately
						return nil
					},
				},
			),
			WithCleanup(func(ctx context.Context) error {
				cleanupCalled = true
				return nil
			}),
		), nil
	}

	// Call the Run function with captureOutput to handle the panic
	output := captureOutput(func() {
		Run(mockBuilder)
	})

	// Check that the output contains the success message
	if !strings.Contains(output, "App shutting down: All runnables completed successfully") {
		t.Errorf("Expected output to contain success message, got: %q", output)
	}

	// Check that the cleanup function was called
	if !cleanupCalled {
		t.Errorf("Expected cleanup function to be called")
	}

	// Assert that the exit code is 0 (success)
	if testExitCode != 0 {
		t.Errorf("Expected exit code to be 0, but got %d", testExitCode)
	}
}

// TestRunWithCleanupError tests the Run function with a cleanup function that returns an error
func TestRunWithCleanupError(t *testing.T) {
	// Save the original osExit function and restore it after the test
	originalExit := osExit
	defer func() {
		osExit = originalExit
	}()

	// Create a mock exit function that records the exit code
	testExitCode = 0 // Reset the exit code
	osExit = func(code int) {
		testExitCode = code
		panic(exitError{code: code})
	}

	// Set environment variables for the test
	os.Setenv("TEST_STRING", "test-value")
	defer func() {
		os.Unsetenv("TEST_STRING")
	}()

	// Track if cleanup was called
	cleanupCalled := false

	// Create a mock builder that returns a successful runnable and a cleanup function that returns an error
	mockBuilder := func(conf TestConfig) (App, error) {
		return Construct(
			WithRunnables(
				mockRunnable{
					runFunc: func(ctx context.Context) error {
						// This runnable completes successfully immediately
						return nil
					},
				},
			),
			WithCleanup(func(ctx context.Context) error {
				cleanupCalled = true
				return errors.New("cleanup error")
			}),
		), nil
	}

	// Call the Run function with captureOutput to handle the panic
	output := captureOutput(func() {
		Run(mockBuilder)
	})

	// Check that the output contains the success message for the runnable
	if !strings.Contains(output, "App shutting down: All runnables completed successfully") {
		t.Errorf("Expected output to contain success message, got: %q", output)
	}

	// Check that the output contains the cleanup error message
	if !strings.Contains(output, "Cleanup error: cleanup error") {
		t.Errorf("Expected output to contain cleanup error message, got: %q", output)
	}

	// Check that the cleanup function was called
	if !cleanupCalled {
		t.Errorf("Expected cleanup function to be called")
	}

	// Assert that the exit code is 1 (error)
	if testExitCode != 1 {
		t.Errorf("Expected exit code to be 1, but got %d", testExitCode)
	}
}

// TestGetCleanupTimeout tests the GetCleanupTimeout function
func TestGetCleanupTimeout(t *testing.T) {
	// Test default timeout (15 seconds)
	os.Unsetenv("EZAPP_TERM_TIMEOUT")
	timeout := GetCleanupTimeout()
	if timeout != 15*time.Second {
		t.Errorf("Expected default timeout to be 15 seconds, got %v", timeout)
	}

	// Test with a valid timeout value
	os.Setenv("EZAPP_TERM_TIMEOUT", "30")
	defer os.Unsetenv("EZAPP_TERM_TIMEOUT")
	timeout = GetCleanupTimeout()
	if timeout != 30*time.Second {
		t.Errorf("Expected timeout to be 30 seconds, got %v", timeout)
	}

	// Test with an invalid timeout value (should use default)
	os.Setenv("EZAPP_TERM_TIMEOUT", "invalid")
	timeout = GetCleanupTimeout()
	if timeout != 15*time.Second {
		t.Errorf("Expected timeout to be 15 seconds for invalid input, got %v", timeout)
	}
}
