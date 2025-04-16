package ezapp

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"
)

// Store the original osExit function
var originalExit = osExit

// Variable to capture the exit code
var testExitCode int

// Custom error type to capture exit code
type exitError struct {
	code int
}

// Override osExit for tests
func init() {
	osExit = func(code int) {
		testExitCode = code
		// Use panic to stop execution, similar to how os.Exit works
		panic(exitError{code: code})
	}
}

// Reset the exit code before each test
func resetExitCode() {
	testExitCode = 0
}

// mockRunnable is a mock implementation of the Runnable interface for testing
type mockRunnable struct {
	runFunc func(ctx context.Context) error
}

func (m mockRunnable) Run(ctx context.Context) error {
	return m.runFunc(ctx)
}

// newSuccessRunnable creates a mock runnable that succeeds after a delay
func newSuccessRunnable(delay time.Duration) mockRunnable {
	return mockRunnable{
		runFunc: func(ctx context.Context) error {
			select {
			case <-time.After(delay):
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		},
	}
}

// newErrorRunnable creates a mock runnable that returns an error after a delay
func newErrorRunnable(delay time.Duration, err error) mockRunnable {
	return mockRunnable{
		runFunc: func(ctx context.Context) error {
			select {
			case <-time.After(delay):
				return err
			case <-ctx.Done():
				return ctx.Err()
			}
		},
	}
}

// captureOutput captures stdout during the execution of a function
// It also recovers from exitError panics
func captureOutput(f func()) string {
	// Save the original stdout
	oldStdout := os.Stdout

	// Create a pipe to capture stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Create a channel to signal when we're done reading
	done := make(chan bool)

	// Create a buffer to store the output
	var buf strings.Builder

	// Start a goroutine to read from the pipe
	go func() {
		io.Copy(&buf, r)
		done <- true
	}()

	// Run the function and recover from exitError panics
	func() {
		defer func() {
			if r := recover(); r != nil {
				// Check if it's our exitError
				if _, ok := r.(exitError); !ok {
					// If it's not our exitError, re-panic
					panic(r)
				}
				// Otherwise, just continue
			}
			// Close the write end of the pipe to flush it
			w.Close()
		}()
		f()
	}()

	// Wait for the reader goroutine to finish
	<-done

	// Restore the original stdout
	os.Stdout = oldStdout

	return buf.String()
}

func TestEzAppRun_NoRunnables(t *testing.T) {
	resetExitCode()
	app := EzApp{
		runnableList: []Runnable{},
	}

	output := captureOutput(func() {
		app.Run()
	})

	expectedOutput := "App shutting down: No runnables to execute\n"
	if output != expectedOutput {
		t.Errorf("Expected output: %q, got: %q", expectedOutput, output)
	}

	// Check that the exit code is 0 (success)
	if testExitCode != 0 {
		t.Errorf("Expected exit code: 0, got: %d", testExitCode)
	}
}

func TestEzAppRun_SuccessfulRunnables(t *testing.T) {
	resetExitCode()
	app := EzApp{
		runnableList: []Runnable{
			newSuccessRunnable(50 * time.Millisecond),
			newSuccessRunnable(100 * time.Millisecond),
		},
	}

	output := captureOutput(func() {
		app.Run()
	})

	expectedOutput := "App shutting down: All runnables completed successfully\n"
	if output != expectedOutput {
		t.Errorf("Expected output: %q, got: %q", expectedOutput, output)
	}

	// Check that the exit code is 0 (success)
	if testExitCode != 0 {
		t.Errorf("Expected exit code: 0, got: %d", testExitCode)
	}
}

func TestEzAppRun_ErrorRunnable(t *testing.T) {
	testError := errors.New("test error")
	resetExitCode()
	app := EzApp{
		runnableList: []Runnable{
			newSuccessRunnable(100 * time.Millisecond),
			newErrorRunnable(50*time.Millisecond, testError),
		},
	}

	output := captureOutput(func() {
		app.Run()
	})

	// Check that the output contains the error message
	if output == "" || !contains(output, "App shutting down: Runnable error: test error") {
		t.Errorf("Expected output to contain error message, got: %q", output)
	}

	// Check that the output contains the shutdown complete message
	if !contains(output, "App shutdown complete") {
		t.Errorf("Expected output to contain shutdown complete message, got: %q", output)
	}

	// Check that the exit code is 1 (error)
	if testExitCode != 1 {
		t.Errorf("Expected exit code: 1, got: %d", testExitCode)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func TestEzAppRun_InitError(t *testing.T) {
	testError := errors.New("initialization error")
	resetExitCode()
	app := EzApp{
		initErr: testError,
	}

	output := captureOutput(func() {
		app.Run()
	})

	// Check that the output contains the initialization error message
	expectedOutput := "App shutting down: Initialization error: initialization error\n"
	if !contains(output, expectedOutput) {
		t.Errorf("Expected output to contain: %q, got: %q", expectedOutput, output)
	}

	// Check that the exit code is 1 (error)
	if testExitCode != 1 {
		t.Errorf("Expected exit code: 1, got: %d", testExitCode)
	}
}

func TestEzAppRun_SignalHandling(t *testing.T) {
	// This test is more complex and may not be reliable in all environments
	// Skip it if we're in a CI environment or if signals can't be properly tested
	if os.Getenv("CI") != "" {
		t.Skip("Skipping signal test in CI environment")
	}

	resetExitCode()
	app := EzApp{
		runnableList: []Runnable{
			newSuccessRunnable(5 * time.Second), // Long-running runnable
		},
	}

	var wg sync.WaitGroup
	wg.Add(1)

	var output string
	go func() {
		defer wg.Done()
		output = captureOutput(func() {
			app.Run()
		})
	}()

	// Give the app a moment to start
	time.Sleep(100 * time.Millisecond)

	// Send a SIGTERM signal to the process
	// Note: This is sending the signal to the current process, which may not be ideal
	// In a real environment, you might want to use a different approach
	p, _ := os.FindProcess(os.Getpid())
	p.Signal(syscall.SIGTERM)

	// Wait for the app to finish
	wg.Wait()

	// Check that the output contains the signal message
	if !contains(output, "App shutting down: Received signal") {
		t.Errorf("Expected output to contain signal message, got: %q", output)
	}

	// Check that the output contains the shutdown complete message
	if !contains(output, "App shutdown complete") {
		t.Errorf("Expected output to contain shutdown complete message, got: %q", output)
	}

	// Check that the exit code is 0 (success) for signal-based termination
	if testExitCode != 0 {
		t.Errorf("Expected exit code: 0, got: %d", testExitCode)
	}
}
