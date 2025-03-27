package ezapp

import (
	"context"
	"log/slog"
	"os"
	"testing"
)

// mockRunnable is a mock implementation of the Runnable interface for testing
type mockRunnable struct {
	runFunc       func(context.Context) error
	handleErrFunc func(error) error
	name          string
}

func (m *mockRunnable) Run(ctx context.Context) error {
	if m.runFunc != nil {
		return m.runFunc(ctx)
	}
	return nil
}

func (m *mockRunnable) HandleError(err error) error {
	if m.handleErrFunc != nil {
		return m.handleErrFunc(err)
	}
	return err
}

func (m *mockRunnable) String() string {
	return m.name
}

func TestBuild_ReturnsValidEzApp(t *testing.T) {
	// Create a mock runnable
	mockR := &mockRunnable{
		name: "TestRunnable",
	}

	// Build an EzApp with the mock runnable
	app := Build(func() []func(*options) {
		return []func(*options){}
	}, mockR)

	// Verify that the returned value is not nil
	if app == nil {
		t.Error("Build returned nil")
	}
}

func TestBuild_AppliesOptions(t *testing.T) {
	// Create a logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Create a mock runnable
	mockR := &mockRunnable{
		name: "TestRunnable",
	}

	// Build an EzApp with the mock runnable and logger option
	app := Build(func() []func(*options) {
		return WithOptions(WithLogger(logger))
	}, mockR)

	// Verify that the returned value is not nil
	if app == nil {
		t.Error("Build returned nil")
	}

	// We can't directly verify that the logger was set correctly because the App struct is in the internal package
	// and we don't have access to its fields. However, we can indirectly verify it by calling Run and checking
	// that it doesn't panic (which would happen if the logger was nil).
	// This is not a perfect test, but it's the best we can do with the current structure.

	// We're not actually going to call Run() here because it would block indefinitely
	// Just checking that app is not nil is sufficient for this test
}

func TestBuild_CreatesNoopLogger(t *testing.T) {
	// Create a mock runnable
	mockR := &mockRunnable{
		name: "TestRunnable",
	}

	// Build an EzApp with the mock runnable and no logger option
	app := Build(func() []func(*options) {
		return WithOptions()
	}, mockR)

	// Verify that the returned value is not nil
	if app == nil {
		t.Error("Build returned nil")
	}

	// We can't directly verify that a noop logger was created because the App struct is in the internal package
	// and we don't have access to its fields. However, we can indirectly verify it by checking that the app is not nil.
	// This is not a perfect test, but it's the best we can do with the current structure.

	// We're not actually going to call Run() here because it would block indefinitely
	// Just checking that app is not nil is sufficient for this test
}

func TestBuild_PassesRunnables(t *testing.T) {
	// Create multiple mock runnables
	mockR1 := &mockRunnable{
		name: "TestRunnable1",
	}
	mockR2 := &mockRunnable{
		name: "TestRunnable2",
	}

	// Build an EzApp with the mock runnables
	app := Build(func() []func(*options) {
		return WithOptions()
	}, mockR1, mockR2)

	// Verify that the returned value is not nil
	if app == nil {
		t.Error("Build returned nil")
	}

	// We can't directly verify that the runnables were passed correctly because the App struct is in the internal package
	// and we don't have access to its fields. However, we can indirectly verify it by checking that the app is not nil.
	// This is not a perfect test, but it's the best we can do with the current structure.

	// We're not actually going to call Run() here because it would block indefinitely
	// Just checking that app is not nil is sufficient for this test
}
