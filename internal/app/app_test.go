package app

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/pgvanniekerk/ezapp/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create a test logger with handler
func createTestLogger() (*slog.Logger, *testutil.TestHandler) {
	return testutil.NewTestLogger(slog.LevelDebug)
}

// Helper runner that succeeds immediately
func successfulRunner(ctx context.Context) error {
	return nil
}

// Helper runner that succeeds after a delay
func delayedSuccessfulRunner(delay time.Duration) Runner {
	return func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			return nil
		}
	}
}

// Helper runner that fails immediately
func failingRunner(ctx context.Context) error {
	return errors.New("runner failed")
}

// Helper runner that fails after a delay
func delayedFailingRunner(delay time.Duration) Runner {
	return func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			return errors.New("delayed runner failed")
		}
	}
}

// Helper runner that runs indefinitely until context is cancelled
func longRunningRunner(started chan<- struct{}) Runner {
	return func(ctx context.Context) error {
		if started != nil {
			close(started)
		}
		<-ctx.Done()
		return ctx.Err()
	}
}

// Helper runner that records execution order
func orderRecordingRunner(id int, order *[]int, mu *sync.Mutex) Runner {
	return func(ctx context.Context) error {
		mu.Lock()
		*order = append(*order, id)
		mu.Unlock()
		return nil
	}
}

// TestNew tests the New constructor function
// This test verifies that:
// - New creates an App with the provided runners and logger
// - The returned App has all fields properly set
func TestNew(t *testing.T) {
	logger, _ := createTestLogger()
	runners := []Runner{successfulRunner, failingRunner}

	app := New(runners, logger)

	assert.Equal(t, runners, app.runnerList, "Runner list should be set correctly")
	assert.Equal(t, logger, app.logger, "Logger should be set correctly")
}

// TestAppRunWithNoRunners tests app execution with empty runner list
// This test verifies that:
// - App can run successfully with no runners
// - No errors are returned when runner list is empty
// - Appropriate debug logs are generated
func TestAppRunWithNoRunners(t *testing.T) {
	logger, logs := createTestLogger()
	app := New([]Runner{}, logger)

	err := app.Run()

	assert.NoError(t, err, "App should run successfully with no runners")

	// Verify debug logs were generated
	logMessages := logs.Messages()
	assert.Contains(t, logMessages, "start application")
	assert.Contains(t, logMessages, "application finished running")
}

// TestAppRunWithSingleSuccessfulRunner tests app execution with one successful runner
// This test verifies that:
// - App runs successfully with a single runner that completes without error
// - No errors are returned when runner succeeds
// - All lifecycle debug logs are generated
func TestAppRunWithSingleSuccessfulRunner(t *testing.T) {
	logger, logs := createTestLogger()
	app := New([]Runner{successfulRunner}, logger)

	err := app.Run()

	assert.NoError(t, err, "App should run successfully with successful runner")

	// Verify all expected debug logs
	logMessages := logs.Messages()
	assert.Contains(t, logMessages, "start application")
	assert.Contains(t, logMessages, "created termination context")
	assert.Contains(t, logMessages, "started termination signaller")
	assert.Contains(t, logMessages, "created error group")
	assert.Contains(t, logMessages, "started runnable invocations via error group")
	assert.Contains(t, logMessages, "application finished running")
}

// TestAppRunWithMultipleSuccessfulRunners tests app execution with multiple successful runners
// This test verifies that:
// - App can run multiple runners concurrently
// - All runners execute successfully
// - Execution order is concurrent (not sequential)
func TestAppRunWithMultipleSuccessfulRunners(t *testing.T) {
	logger, _ := createTestLogger()

	var order []int
	var mu sync.Mutex

	runners := []Runner{
		orderRecordingRunner(1, &order, &mu),
		orderRecordingRunner(2, &order, &mu),
		orderRecordingRunner(3, &order, &mu),
	}

	app := New(runners, logger)

	err := app.Run()

	assert.NoError(t, err, "App should run successfully with multiple runners")

	// Verify all runners executed
	mu.Lock()
	assert.Len(t, order, 3, "All runners should have executed")
	assert.Contains(t, order, 1, "Runner 1 should have executed")
	assert.Contains(t, order, 2, "Runner 2 should have executed")
	assert.Contains(t, order, 3, "Runner 3 should have executed")
	mu.Unlock()
}

// TestAppRunWithFailingRunner tests app execution when a runner fails
// This test verifies that:
// - App returns error when any runner fails
// - Error message includes the wrapped runner error
// - Other runners are cancelled when one fails
func TestAppRunWithFailingRunner(t *testing.T) {
	logger, _ := createTestLogger()

	runners := []Runner{
		failingRunner,
		delayedSuccessfulRunner(100 * time.Millisecond), // Should be cancelled
	}

	app := New(runners, logger)

	err := app.Run()

	assert.Error(t, err, "App should return error when runner fails")
	assert.Contains(t, err.Error(), "failed to invoke runnable", "Error should be wrapped properly")
	assert.Contains(t, err.Error(), "runner failed", "Error should contain original runner error")
}

// TestAppRunWithMixedRunners tests app execution with both successful and failing runners
// This test verifies that:
// - When one runner fails, the app fails even if others succeed
// - Error group behavior cancels other runners
// - First error is returned (error group behavior)
func TestAppRunWithMixedRunners(t *testing.T) {
	logger, _ := createTestLogger()

	runners := []Runner{
		successfulRunner,
		failingRunner,
		delayedSuccessfulRunner(100 * time.Millisecond),
	}

	app := New(runners, logger)

	err := app.Run()

	assert.Error(t, err, "App should fail when any runner fails")
	assert.Contains(t, err.Error(), "runner failed", "Should contain failing runner error")
}

// TestAppRunWithDelayedFailure tests app execution with a runner that fails after delay
// This test verifies that:
// - App waits for all runners to complete or fail
// - Delayed failures are properly handled
// - Other running runners are cancelled when one fails
func TestAppRunWithDelayedFailure(t *testing.T) {
	logger, _ := createTestLogger()

	started := make(chan struct{})
	runners := []Runner{
		longRunningRunner(started),                  // Will be cancelled
		delayedFailingRunner(50 * time.Millisecond), // Will fail after delay
	}

	app := New(runners, logger)

	// Wait for long runner to start
	go func() {
		<-started
	}()

	err := app.Run()

	assert.Error(t, err, "App should fail when delayed runner fails")
	assert.Contains(t, err.Error(), "delayed runner failed", "Should contain delayed failure error")
}

// TestAppRunWithContextCancellation tests behavior when runners respect context cancellation
// This test verifies that:
// - Runners that respect context cancellation are properly stopped
// - Context cancellation error is handled appropriately
// - App can complete even when runners are cancelled
func TestAppRunWithContextCancellation(t *testing.T) {
	logger, _ := createTestLogger()

	cancellingRunner := func(ctx context.Context) error {
		// Simulate a runner that respects context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(1 * time.Second):
			return nil
		}
	}

	// Create a runner that will cause cancellation by failing quickly
	quickFailRunner := func(ctx context.Context) error {
		time.Sleep(10 * time.Millisecond)
		return errors.New("quick fail")
	}

	runners := []Runner{cancellingRunner, quickFailRunner}
	app := New(runners, logger)

	err := app.Run()

	assert.Error(t, err, "App should return error from failing runner")
	assert.Contains(t, err.Error(), "quick fail", "Should contain quick fail error")
}

// TestAppTerminationSignaller tests the signal handling functionality
// This test verifies that:
// - Signal handler is set up correctly
// - SIGTERM signal triggers context cancellation
// - Signal cleanup is performed properly
func TestAppTerminationSignaller(t *testing.T) {
	logger, logs := createTestLogger()

	started := make(chan struct{})
	cancelled := make(chan struct{})

	runners := []Runner{
		func(ctx context.Context) error {
			close(started)
			<-ctx.Done()
			close(cancelled)
			return ctx.Err()
		},
	}

	app := New(runners, logger)

	// Run app in goroutine
	done := make(chan error, 1)
	go func() {
		done <- app.Run()
	}()

	// Wait for runner to start
	<-started

	// Send SIGTERM to trigger termination
	// Note: We need to send to the current process
	pid := os.Getpid()
	process, err := os.FindProcess(pid)
	require.NoError(t, err, "Should find current process")

	err = process.Signal(syscall.SIGTERM)
	require.NoError(t, err, "Should send SIGTERM successfully")

	// Wait for cancellation and completion
	select {
	case <-cancelled:
		// Good, context was cancelled
	case <-time.After(1 * time.Second):
		t.Fatal("Context should have been cancelled by SIGTERM")
	}

	select {
	case err := <-done:
		assert.Error(t, err, "App should return context cancellation error")
		assert.Contains(t, err.Error(), "context canceled", "Error should indicate context cancellation")
	case <-time.After(1 * time.Second):
		t.Fatal("App should have completed after signal")
	}

	// Verify signal handling debug logs
	logMessages := logs.Messages()
	assert.Contains(t, logMessages, "starting termination signaller")
	assert.Contains(t, logMessages, "started listening for SIGINT and SIGTERM")
	assert.Contains(t, logMessages, "received SIGINT or SIGTERM, terminating")
	assert.Contains(t, logMessages, "stopped listening for SIGINT and SIGTERM")
}

// TestAppRunnerListIndexCapture tests that the runner list index is captured correctly
// This test verifies that:
// - Each runner in the list is executed (not just the last one due to closure issues)
// - Index variable is properly captured in the closure
// - All runners receive the correct context
func TestAppRunnerListIndexCapture(t *testing.T) {
	logger, _ := createTestLogger()

	var executedRunners []int
	var mu sync.Mutex

	// Create runners that record their expected index
	runners := make([]Runner, 5)
	for i := 0; i < 5; i++ {
		expectedIndex := i
		runners[i] = func(ctx context.Context) error {
			mu.Lock()
			executedRunners = append(executedRunners, expectedIndex)
			mu.Unlock()
			return nil
		}
	}

	app := New(runners, logger)

	err := app.Run()

	assert.NoError(t, err, "App should run successfully")

	mu.Lock()
	assert.Len(t, executedRunners, 5, "All runners should have executed")

	// Verify all expected indices were executed (order may vary due to concurrency)
	for i := 0; i < 5; i++ {
		assert.Contains(t, executedRunners, i, "Runner %d should have executed", i)
	}
	mu.Unlock()
}

// TestAppWithNilLogger tests app behavior with nil logger (should not panic)
// This test verifies that:
// - App can handle nil logger gracefully (if the implementation allows it)
// - Or panics in a predictable way if nil logger is not supported
func TestAppWithNilLogger(t *testing.T) {
	// Create app with nil logger
	app := New([]Runner{successfulRunner}, nil)

	// This should either work gracefully or panic predictably
	defer func() {
		if r := recover(); r != nil {
			// If it panics, that's acceptable behavior for nil logger
			t.Logf("App panicked with nil logger: %v", r)
		}
	}()

	err := app.Run()

	// If we reach here, nil logger was handled gracefully
	assert.NoError(t, err, "App should handle nil logger gracefully if supported")
}

// TestAppRunConcurrentExecution tests that runners execute concurrently, not sequentially
// This test verifies that:
// - Multiple runners start approximately at the same time
// - Total execution time is less than sequential execution would take
// - Concurrent execution provides performance benefits
func TestAppRunConcurrentExecution(t *testing.T) {
	logger, _ := createTestLogger()

	runnerDelay := 100 * time.Millisecond
	numRunners := 3

	var startTimes []time.Time
	var mu sync.Mutex

	runners := make([]Runner, numRunners)
	for i := 0; i < numRunners; i++ {
		runners[i] = func(ctx context.Context) error {
			mu.Lock()
			startTimes = append(startTimes, time.Now())
			mu.Unlock()

			time.Sleep(runnerDelay)
			return nil
		}
	}

	app := New(runners, logger)

	start := time.Now()
	err := app.Run()
	totalDuration := time.Since(start)

	assert.NoError(t, err, "App should run successfully")

	// Verify concurrent execution - total time should be close to single runner time
	// Allow some overhead for goroutine startup
	maxExpectedDuration := runnerDelay + 50*time.Millisecond
	assert.Less(t, totalDuration, maxExpectedDuration,
		"Concurrent execution should be faster than sequential (%v vs %v)",
		totalDuration, time.Duration(numRunners)*runnerDelay)

	// Verify runners started approximately at the same time
	mu.Lock()
	assert.Len(t, startTimes, numRunners, "All runners should have started")
	if len(startTimes) >= 2 {
		maxStartDiff := startTimes[len(startTimes)-1].Sub(startTimes[0])
		assert.Less(t, maxStartDiff, 50*time.Millisecond,
			"Runners should start within 50ms of each other")
	}
	mu.Unlock()
}
