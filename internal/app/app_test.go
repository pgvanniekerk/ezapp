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

func TestApp_Run_Success(t *testing.T) {
	// Create a logger that writes to a buffer
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Create a mock runnable that succeeds
	runCalled := false
	mockR := &mockRunnable{
		name: "SuccessRunnable",
		runFunc: func(ctx context.Context) error {
			runCalled = true
			return nil
		},
	}

	// Create an app with the mock runnable
	app := New([]Runnable{mockR}, logger)

	// Run the app in a goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	var runErr error
	go func() {
		defer wg.Done()
		runErr = app.Run()
	}()

	// Wait a bit for the app to start
	time.Sleep(100 * time.Millisecond)

	// Cancel the context to stop the app
	app.stopCtxCancel()

	// Wait for the app to finish
	wg.Wait()

	// Check that Run was called and no error was returned
	if !runCalled {
		t.Error("Run was not called on the runnable")
	}
	if runErr != nil {
		t.Errorf("Expected no error, got %v", runErr)
	}
}

func TestApp_Run_RunnableError(t *testing.T) {
	// Create a logger that writes to a buffer
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Create a mock runnable that returns an error
	expectedErr := errors.New("runnable error")
	mockR := &mockRunnable{
		name: "ErrorRunnable",
		runFunc: func(ctx context.Context) error {
			return expectedErr
		},
		handleErrFunc: func(err error) error {
			// Return the error unchanged
			return err
		},
	}

	// Create an app with the mock runnable
	app := New([]Runnable{mockR}, logger)

	// Run the app
	err := app.Run()

	// Check that the error was propagated
	if err != expectedErr {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}

func TestApp_Run_RunnableErrorHandled(t *testing.T) {
	// Create a logger that writes to a buffer
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Create a mock runnable that returns an error but handles it
	runErr := errors.New("runnable error")
	mockR := &mockRunnable{
		name: "HandledErrorRunnable",
		runFunc: func(ctx context.Context) error {
			return runErr
		},
		handleErrFunc: func(err error) error {
			// Return nil to indicate the error was handled
			return nil
		},
	}

	// Create an app with the mock runnable
	app := New([]Runnable{mockR}, logger)

	// Run the app in a goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	var appErr error
	go func() {
		defer wg.Done()
		appErr = app.Run()
	}()

	// Wait a bit for the app to start and handle the error
	time.Sleep(100 * time.Millisecond)

	// Cancel the context to stop the app
	app.stopCtxCancel()

	// Wait for the app to finish
	wg.Wait()

	// Check that no error was returned
	if appErr != nil {
		t.Errorf("Expected no error, got %v", appErr)
	}
}

func TestApp_Run_SigTerm(t *testing.T) {
	// Create a logger that writes to a buffer
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Create a mock runnable that blocks until context is done
	runCalled := false
	ctxDone := false
	mockR := &mockRunnable{
		name: "BlockingRunnable",
		runFunc: func(ctx context.Context) error {
			runCalled = true
			<-ctx.Done()
			ctxDone = true
			return nil
		},
	}

	// Create an app with the mock runnable
	app := New([]Runnable{mockR}, logger)

	// Run the app in a goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	var runErr error
	go func() {
		defer wg.Done()
		runErr = app.Run()
	}()

	// Wait a bit for the app to start
	time.Sleep(100 * time.Millisecond)

	// Send a SIGTERM signal
	app.sigTerm <- syscall.SIGTERM

	// Wait for the app to finish
	wg.Wait()

	// Check that Run was called, context was done, and no error was returned
	if !runCalled {
		t.Error("Run was not called on the runnable")
	}
	if !ctxDone {
		t.Error("Context was not cancelled")
	}
	if runErr != nil {
		t.Errorf("Expected no error, got %v", runErr)
	}
}

func TestApp_Run_MultipleRunnables(t *testing.T) {
	// Create a logger that writes to a buffer
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Create multiple mock runnables
	run1Called := false
	mockR1 := &mockRunnable{
		name: "Runnable1",
		runFunc: func(ctx context.Context) error {
			run1Called = true
			return nil
		},
	}

	run2Called := false
	mockR2 := &mockRunnable{
		name: "Runnable2",
		runFunc: func(ctx context.Context) error {
			run2Called = true
			return nil
		},
	}

	// Create an app with the mock runnables
	app := New([]Runnable{mockR1, mockR2}, logger)

	// Run the app in a goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	var runErr error
	go func() {
		defer wg.Done()
		runErr = app.Run()
	}()

	// Wait a bit for the app to start
	time.Sleep(100 * time.Millisecond)

	// Cancel the context to stop the app
	app.stopCtxCancel()

	// Wait for the app to finish
	wg.Wait()

	// Check that both Run methods were called and no error was returned
	if !run1Called {
		t.Error("Run was not called on the first runnable")
	}
	if !run2Called {
		t.Error("Run was not called on the second runnable")
	}
	if runErr != nil {
		t.Errorf("Expected no error, got %v", runErr)
	}
}

func TestApp_Run_ContextCanceled(t *testing.T) {
	// Create a logger that writes to a buffer
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Create a mock runnable that returns context.Canceled
	mockR := &mockRunnable{
		name: "ContextCanceledRunnable",
		runFunc: func(ctx context.Context) error {
			return context.Canceled
		},
	}

	// Create an app with the mock runnable
	app := New([]Runnable{mockR}, logger)

	// Run the app
	err := app.Run()

	// Check that context.Canceled was ignored
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}