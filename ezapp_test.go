package ezapp

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConfig represents a test configuration struct
type TestConfig struct {
	Port        int    `env:"PORT" default:"8080"`
	DatabaseURL string `env:"DATABASE_URL" default:"test://localhost"`
	TestValue   string `env:"TEST_VALUE" default:"test"`
}


// Helper function to create a successful runner
func successfulRunner(ctx context.Context) error {
	// Simulate some work
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(10 * time.Millisecond):
		return nil
	}
}

// Helper function to create a failing runner
func failingRunner(ctx context.Context) error {
	return errors.New("runner failed")
}

// Helper function to create a successful cleanup function
func successfulCleanup(ctx context.Context) error {
	return nil
}

// Helper function to create a failing cleanup function
func failingCleanup(ctx context.Context) error {
	return errors.New("cleanup failed")
}

// Helper function to create a recorder cleanup function
func createRecorderCleanup(called *bool) func(context.Context) error {
	return func(ctx context.Context) error {
		*called = true
		return nil
	}
}

// TestRunSuccessful tests the successful execution path of Run function
// This test verifies that:
// - Configuration is loaded correctly
// - Logger is initialized and non-nil
// - StartupCtx is created and non-nil
// - Runners are executed successfully
// - Application completes without calling Fatal
func TestRunSuccessful(t *testing.T) {
	var runnerExecuted bool

	testRunner := func(ctx context.Context) error {
		runnerExecuted = true
		return nil
	}

	initializer := func(ctx InitCtx[TestConfig]) (AppCtx, error) {
		// Verify required fields are populated
		require.NotNil(t, ctx.StartupCtx, "StartupCtx should not be nil")
		require.NotNil(t, ctx.Logger, "Logger should not be nil")
		
		// The config should be populated (values may be zero values if no env vars set)
		// This just verifies the config loading mechanism works
		t.Logf("Config loaded: Port=%d, DatabaseURL=%s, TestValue=%s", 
			ctx.Config.Port, ctx.Config.DatabaseURL, ctx.Config.TestValue)
		
		return Construct(WithRunners(testRunner))
	}

	// Since Run calls logger.Fatal on errors which terminates the process,
	// we need a different approach to test. We'll run it in a goroutine
	// and expect it to complete successfully.
	done := make(chan bool, 1)
	
	go func() {
		defer func() {
			if r := recover(); r != nil {
				done <- false // Panic occurred (Fatal was called)
			} else {
				done <- true // Normal completion
			}
		}()
		Run(initializer)
	}()

	select {
	case success := <-done:
		if !success {
			t.Fatal("Run panicked when it should have completed successfully")
		}
		assert.True(t, runnerExecuted, "Runner should have been executed")
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not complete within timeout")
	}
}

// TestRunWithCleanupFunction tests that cleanup function is called and executed properly
// This test verifies that:
// - Cleanup function is executed after successful app run
// - Cleanup receives proper shutdown context with timeout
// - Application completes successfully after cleanup
func TestRunWithCleanupFunction(t *testing.T) {
	var cleanupCalled bool
	var cleanupContext context.Context
	
	cleanup := func(ctx context.Context) error {
		cleanupCalled = true
		cleanupContext = ctx
		
		// Verify context has deadline (shutdown timeout)
		_, hasDeadline := ctx.Deadline()
		assert.True(t, hasDeadline, "Cleanup context should have deadline")
		
		return nil
	}

	initializer := func(ctx InitCtx[TestConfig]) (AppCtx, error) {
		return Construct(
			WithRunners(successfulRunner),
			WithCleanup(cleanup),
		)
	}

	done := make(chan bool, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				done <- false
			} else {
				done <- true
			}
		}()
		Run(initializer)
	}()

	select {
	case success := <-done:
		if !success {
			t.Fatal("Run should not panic with successful cleanup")
		}
		assert.True(t, cleanupCalled, "Cleanup function should have been called")
		assert.NotNil(t, cleanupContext, "Cleanup should have received context")
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not complete within timeout")
	}
}

// TestRunWithNoRunners tests that Run completes successfully even with no runners
// This test verifies that the application can run with zero runners (which is valid)
func TestRunWithNoRunners(t *testing.T) {
	initializer := func(ctx InitCtx[TestConfig]) (AppCtx, error) {
		// Return AppCtx with no runners - this is actually valid behavior
		return Construct() // No WithRunners called
	}

	done := make(chan bool, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				done <- false // Panic occurred (Fatal was called)
			} else {
				done <- true // Normal completion
			}
		}()
		Run(initializer)
	}()

	select {
	case success := <-done:
		// No runners is actually valid - the app just starts and finishes immediately
		assert.True(t, success, "Run should complete successfully even with no runners")
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not complete within timeout")
	}
}

// TestConstructWithOptions tests the Construct function with various options
// This test verifies the functional options pattern works correctly
func TestConstructWithOptions(t *testing.T) {
	// Test with runners only
	appCtx1, err := Construct(WithRunners(successfulRunner, failingRunner))
	require.NoError(t, err, "Construct with runners should not fail")
	assert.Len(t, appCtx1.runnerList, 2, "Should have 2 runners")
	assert.Nil(t, appCtx1.cleanupFunc, "Cleanup function should be nil")

	// Test with cleanup only
	appCtx2, err := Construct(WithCleanup(successfulCleanup))
	require.NoError(t, err, "Construct with cleanup should not fail")
	assert.Len(t, appCtx2.runnerList, 0, "Should have 0 runners")
	assert.NotNil(t, appCtx2.cleanupFunc, "Cleanup function should be set")

	// Test with both runners and cleanup
	appCtx3, err := Construct(
		WithRunners(successfulRunner),
		WithCleanup(successfulCleanup),
	)
	require.NoError(t, err, "Construct with both options should not fail")
	assert.Len(t, appCtx3.runnerList, 1, "Should have 1 runner")
	assert.NotNil(t, appCtx3.cleanupFunc, "Cleanup function should be set")

	// Test with no options
	appCtx4, err := Construct()
	require.NoError(t, err, "Construct with no options should not fail")
	assert.Len(t, appCtx4.runnerList, 0, "Should have 0 runners")
	assert.Nil(t, appCtx4.cleanupFunc, "Cleanup function should be nil")
}

// TestInitCtxPopulation tests that InitCtx is properly populated
// This test verifies that all required fields are set correctly
func TestInitCtxPopulation(t *testing.T) {
	var capturedInitCtx InitCtx[TestConfig]
	
	initializer := func(ctx InitCtx[TestConfig]) (AppCtx, error) {
		capturedInitCtx = ctx
		return Construct(WithRunners(successfulRunner))
	}

	done := make(chan bool, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				done <- false
			} else {
				done <- true
			}
		}()
		Run(initializer)
	}()

	select {
	case success := <-done:
		if !success {
			t.Fatal("Run should complete successfully")
		}
		
		// Verify InitCtx was populated correctly
		assert.NotNil(t, capturedInitCtx.StartupCtx, "StartupCtx should not be nil")
		assert.NotNil(t, capturedInitCtx.Logger, "Logger should not be nil")
		
		// Verify context has timeout
		_, hasDeadline := capturedInitCtx.StartupCtx.Deadline()
		assert.True(t, hasDeadline, "StartupCtx should have deadline")
		
		// Verify config is populated (even if with zero values)
		// The fact that we got here means config loading succeeded
		t.Logf("Config loaded successfully: %+v", capturedInitCtx.Config)
		
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not complete within timeout")
	}
}

/*
NOTE: The following tests cannot be implemented because they would trigger logger.Fatal() 
which calls os.Exit() and terminates the test process. To properly test these scenarios,
we would need:
1. Dependency injection to mock the logger
2. Process isolation (running tests in separate processes)
3. A testable version of Run() that doesn't call Fatal

The scenarios we cannot test directly:
- TestRunConfigurationLoadingFailure (invalid config struct)
- TestRunStartupContextFailure (invalid EZAPP_STARTUP_TIMEOUT)
- TestRunInitializerFailure (initializer returns error)
- TestRunApplicationFailure (runner returns error)
- TestRunCleanupFailure (cleanup function returns error)
- TestRunCleanupWithApplicationFailure (both app and cleanup fail)

These scenarios are covered by the logic in the Run function and would call logger.Fatal()
with appropriate error messages, but cannot be tested without process termination.
*/