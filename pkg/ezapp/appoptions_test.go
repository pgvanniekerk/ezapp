package ezapp

import (
	"context"
	"testing"
)

// TestConstruct tests the Construct function
func TestConstruct(t *testing.T) {
	// Create an App with no options
	app := Construct()
	if len(app.runnables) != 0 {
		t.Errorf("Expected 0 runnables, got %d", len(app.runnables))
	}
	if app.cleanup != nil {
		t.Errorf("Expected nil cleanup function, but it was not nil")
	}
}

// TestWithRunnables tests the WithRunnables function
func TestWithRunnables(t *testing.T) {
	// Create a mock runnable
	mockRun := mockRunnable{
		runFunc: func(ctx context.Context) error {
			return nil
		},
	}

	// Create an App with one runnable
	app := Construct(WithRunnables(mockRun))
	if len(app.runnables) != 1 {
		t.Errorf("Expected 1 runnable, got %d", len(app.runnables))
	}

	// Create an App with multiple runnables
	app = Construct(WithRunnables(mockRun, mockRun))
	if len(app.runnables) != 2 {
		t.Errorf("Expected 2 runnables, got %d", len(app.runnables))
	}

	// Create an App with multiple WithRunnables calls
	app = Construct(WithRunnables(mockRun), WithRunnables(mockRun))
	if len(app.runnables) != 2 {
		t.Errorf("Expected 2 runnables, got %d", len(app.runnables))
	}
}

// TestWithCleanup tests the WithCleanup function
func TestWithCleanup(t *testing.T) {
	// Create a flag to track if cleanup was called
	cleanupCalled := false

	// Create a cleanup function
	cleanup := func(ctx context.Context) error {
		cleanupCalled = true
		return nil
	}

	// Create an App with a cleanup function
	app := Construct(WithCleanup(cleanup))
	if app.cleanup == nil {
		t.Errorf("Expected non-nil cleanup function")
	}

	// Call the cleanup function
	err := app.cleanup(context.Background())
	if err != nil {
		t.Errorf("Expected no error from cleanup, got %v", err)
	}
	if !cleanupCalled {
		t.Errorf("Expected cleanup function to be called")
	}
}

// TestMultipleOptions tests using multiple options together
func TestMultipleOptions(t *testing.T) {
	// Create a mock runnable
	mockRun := mockRunnable{
		runFunc: func(ctx context.Context) error {
			return nil
		},
	}

	// Create a flag to track if cleanup was called
	cleanupCalled := false

	// Create a cleanup function
	cleanup := func(ctx context.Context) error {
		cleanupCalled = true
		return nil
	}

	// Create an App with both runnables and cleanup
	app := Construct(
		WithRunnables(mockRun, mockRun),
		WithCleanup(cleanup),
	)

	// Check that both options were applied
	if len(app.runnables) != 2 {
		t.Errorf("Expected 2 runnables, got %d", len(app.runnables))
	}
	if app.cleanup == nil {
		t.Errorf("Expected non-nil cleanup function")
	}

	// Call the cleanup function
	err := app.cleanup(context.Background())
	if err != nil {
		t.Errorf("Expected no error from cleanup, got %v", err)
	}
	if !cleanupCalled {
		t.Errorf("Expected cleanup function to be called")
	}
}
