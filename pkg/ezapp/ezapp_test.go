package ezapp

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"
	"time"
)

// MockRunnable is a mock implementation of the Runnable interface for testing
type MockRunnable struct {
	RunFunc       func(context.Context) error
	HandleErrFunc func(error) error
	RunCalled     bool
	HandleErrCalled bool
	mu            sync.Mutex
}

func (m *MockRunnable) Run(ctx context.Context) error {
	m.mu.Lock()
	m.RunCalled = true
	m.mu.Unlock()
	return m.RunFunc(ctx)
}

func (m *MockRunnable) HandleError(err error) error {
	m.mu.Lock()
	m.HandleErrCalled = true
	m.mu.Unlock()
	return m.HandleErrFunc(err)
}

func TestEzApp_Run_Success(t *testing.T) {
	// Create a mock runnable that succeeds
	runnable := &MockRunnable{
		RunFunc: func(ctx context.Context) error {
			return nil
		},
		HandleErrFunc: func(err error) error {
			return nil
		},
	}

	// Create a cleanup function that succeeds
	cleanupCalled := false
	cleanupFunc := func() error {
		cleanupCalled = true
		return nil
	}

	// Create an EzApp with the mock runnable
	app := &EzApp{
		runnables:    []Runnable{runnable},
		errorHandler: nil,
		cleanupFunc:  cleanupFunc,
	}

	// Run the app
	app.Run()

	// Verify that the runnable was called
	if !runnable.RunCalled {
		t.Error("Expected Run to be called on the runnable, but it wasn't")
	}

	// Verify that the cleanup function was called
	if !cleanupCalled {
		t.Error("Expected cleanup function to be called, but it wasn't")
	}
}

func TestEzApp_Run_Error(t *testing.T) {
	// Create an error to return
	testErr := errors.New("test error")

	// Create a mock runnable that returns an error
	runnable := &MockRunnable{
		RunFunc: func(ctx context.Context) error {
			return testErr
		},
		HandleErrFunc: func(err error) error {
			// Return nil to indicate the error was handled
			return nil
		},
	}

	// Create a cleanup function that succeeds
	cleanupCalled := false
	cleanupFunc := func() error {
		cleanupCalled = true
		return nil
	}

	// Create an EzApp with the mock runnable
	app := &EzApp{
		runnables:    []Runnable{runnable},
		errorHandler: nil,
		cleanupFunc:  cleanupFunc,
	}

	// Run the app
	app.Run()

	// Verify that the runnable was called
	if !runnable.RunCalled {
		t.Error("Expected Run to be called on the runnable, but it wasn't")
	}

	// Verify that the error handler was called
	if !runnable.HandleErrCalled {
		t.Error("Expected HandleError to be called on the runnable, but it wasn't")
	}

	// Verify that the cleanup function was called
	if !cleanupCalled {
		t.Error("Expected cleanup function to be called, but it wasn't")
	}
}

func TestEzApp_Run_ErrorPropagation(t *testing.T) {
	// Create an error to return
	testErr := errors.New("test error")

	// Create a mock runnable that returns an error
	runnable := &MockRunnable{
		RunFunc: func(ctx context.Context) error {
			return testErr
		},
		HandleErrFunc: func(err error) error {
			// Return the error to indicate it wasn't handled
			return err
		},
	}

	// Track if the app-level error handler was called
	appErrorHandlerCalled := false

	// Create an app-level error handler
	appErrorHandler := func(err error) error {
		appErrorHandlerCalled = true
		// Return nil to indicate the error was handled
		return nil
	}

	// Create a cleanup function that succeeds
	cleanupCalled := false
	cleanupFunc := func() error {
		cleanupCalled = true
		return nil
	}

	// Create an EzApp with the mock runnable and error handler
	app := &EzApp{
		runnables:    []Runnable{runnable},
		errorHandler: appErrorHandler,
		cleanupFunc:  cleanupFunc,
	}

	// Run the app
	app.Run()

	// Verify that the runnable was called
	if !runnable.RunCalled {
		t.Error("Expected Run to be called on the runnable, but it wasn't")
	}

	// Verify that the error handler was called
	if !runnable.HandleErrCalled {
		t.Error("Expected HandleError to be called on the runnable, but it wasn't")
	}

	// Verify that the app-level error handler was called
	if !appErrorHandlerCalled {
		t.Error("Expected app-level error handler to be called, but it wasn't")
	}

	// Verify that the cleanup function was called
	if !cleanupCalled {
		t.Error("Expected cleanup function to be called, but it wasn't")
	}
}

func TestEzApp_Run_ContextCancellation(t *testing.T) {
	// Create a channel to signal when the runnable has started
	started := make(chan struct{})

	// Create a channel to signal when the context is canceled
	canceled := make(chan struct{})

	// Create a mock runnable that blocks until the context is canceled
	runnable := &MockRunnable{
		RunFunc: func(ctx context.Context) error {
			// Signal that the runnable has started
			close(started)

			// Wait for the context to be canceled
			<-ctx.Done()

			// Signal that the context was canceled
			close(canceled)

			return ctx.Err()
		},
		HandleErrFunc: func(err error) error {
			return err
		},
	}

	// Create a cleanup function that succeeds
	cleanupCalled := false
	cleanupFunc := func() error {
		cleanupCalled = true
		return nil
	}

	// Create an EzApp with the mock runnable
	app := &EzApp{
		runnables:    []Runnable{runnable},
		errorHandler: nil,
		cleanupFunc:  cleanupFunc,
	}

	// Create a channel to signal when the app.Run() has completed
	done := make(chan struct{})

	// Run the app in a goroutine
	go func() {
		app.Run()
		close(done)
	}()

	// Wait for the runnable to start
	<-started

	// Wait a bit to ensure the goroutine is running
	time.Sleep(100 * time.Millisecond)

	// Simulate a Ctrl+C signal by sending an interrupt signal
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("Failed to find process: %v", err)
	}

	// Send the interrupt signal
	err = p.Signal(os.Interrupt)
	if err != nil {
		t.Fatalf("Failed to send signal: %v", err)
	}

	// Wait for the context to be canceled (with a timeout)
	select {
	case <-canceled:
		// Context was canceled, which is what we expect
	case <-time.After(5 * time.Second):
		t.Error("Context was not canceled within the expected time")
	}

	// Wait for the app.Run() to complete
	select {
	case <-done:
		// app.Run() completed, which is what we expect
	case <-time.After(5 * time.Second):
		t.Error("app.Run() did not complete within the expected time")
	}

	// Verify that the cleanup function was called
	if !cleanupCalled {
		t.Error("Expected cleanup function to be called, but it wasn't")
	}
}

func TestEzApp_Run_CleanupError(t *testing.T) {
	// Create a mock runnable that succeeds
	runnable := &MockRunnable{
		RunFunc: func(ctx context.Context) error {
			return nil
		},
		HandleErrFunc: func(err error) error {
			return nil
		},
	}

	// Create a cleanup function that returns an error
	cleanupErr := errors.New("cleanup error")
	cleanupFunc := func() error {
		return cleanupErr
	}

	// Create an EzApp with the mock runnable
	app := &EzApp{
		runnables:    []Runnable{runnable},
		errorHandler: nil,
		cleanupFunc:  cleanupFunc,
	}

	// Run the app and expect a panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected Run to panic when cleanup function returns an error, but it didn't")
		}
	}()

	app.Run()
}

func TestEzApp_Run_MultipleRunnables(t *testing.T) {
	// Create a counter to track how many runnables were called
	var counter int
	var mu sync.Mutex

	// Create a function to increment the counter
	incrementCounter := func() {
		mu.Lock()
		counter++
		mu.Unlock()
	}

	// Create multiple mock runnables
	runnable1 := &MockRunnable{
		RunFunc: func(ctx context.Context) error {
			incrementCounter()
			return nil
		},
		HandleErrFunc: func(err error) error {
			return nil
		},
	}

	runnable2 := &MockRunnable{
		RunFunc: func(ctx context.Context) error {
			incrementCounter()
			return nil
		},
		HandleErrFunc: func(err error) error {
			return nil
		},
	}

	runnable3 := &MockRunnable{
		RunFunc: func(ctx context.Context) error {
			incrementCounter()
			return nil
		},
		HandleErrFunc: func(err error) error {
			return nil
		},
	}

	// Create a cleanup function that succeeds
	cleanupCalled := false
	cleanupFunc := func() error {
		cleanupCalled = true
		return nil
	}

	// Create an EzApp with the mock runnables
	app := &EzApp{
		runnables:    []Runnable{runnable1, runnable2, runnable3},
		errorHandler: nil,
		cleanupFunc:  cleanupFunc,
	}

	// Run the app
	app.Run()

	// Verify that all runnables were called
	if counter != 3 {
		t.Errorf("Expected all 3 runnables to be called, but got %d", counter)
	}

	// Verify that the cleanup function was called
	if !cleanupCalled {
		t.Error("Expected cleanup function to be called, but it wasn't")
	}
}

// TestBuild_Success tests that the Build function successfully creates an EzApp with no options
func TestBuild_Success(t *testing.T) {
	// Create a mock runnable
	runnable := &MockRunnable{
		RunFunc: func(ctx context.Context) error {
			return nil
		},
		HandleErrFunc: func(err error) error {
			return nil
		},
	}

	// Create a mock cleanup function
	cleanupCalled := false
	cleanupFunc := func() error {
		cleanupCalled = true
		return nil
	}

	// Create a mock wire function that succeeds
	wireFunc := func(config struct{}) (WireBundle, error) {
		return WireBundle{
			Runnables:   []Runnable{runnable},
			CleanupFunc: cleanupFunc,
		}, nil
	}

	// Call the Build function with no options
	app := Build(wireFunc)

	// Verify that the app was created with the correct runnables
	if len(app.runnables) != 1 {
		t.Errorf("Expected app to have 1 runnable, but got %d", len(app.runnables))
	}

	// Verify that the app was created with no error handler
	if app.errorHandler != nil {
		t.Error("Expected app to have no error handler, but it does")
	}

	// Verify that the app was created with the correct cleanup function
	if app.cleanupFunc == nil {
		t.Error("Expected app to have a cleanup function, but it doesn't")
	}

	// Verify that the cleanup function is the one from the wire bundle
	if app.cleanupFunc != nil {
		// Call the cleanup function to verify it's the one from the wire bundle
		app.cleanupFunc()
		if !cleanupCalled {
			t.Error("Expected wire bundle cleanup function to be called, but it wasn't")
		}
	}
}

// TestBuild_WithOptions tests that the Build function successfully creates an EzApp with options
func TestBuild_WithOptions(t *testing.T) {
	// Create a mock runnable
	runnable := &MockRunnable{
		RunFunc: func(ctx context.Context) error {
			return nil
		},
		HandleErrFunc: func(err error) error {
			return nil
		},
	}

	// Create a mock cleanup function for the wire bundle
	wireCleanupFunc := func() error {
		return nil
	}

	// Create a mock cleanup function for the options
	optionsCleanupCalled := false
	optionsCleanupFunc := func() error {
		optionsCleanupCalled = true
		return nil
	}

	// Create a mock error handler
	errHandlerCalled := false
	errHandler := func(err error) error {
		errHandlerCalled = true
		return nil
	}

	// Create a mock wire function that succeeds
	wireFunc := func(config struct{}) (WireBundle, error) {
		return WireBundle{
			Runnables:   []Runnable{runnable},
			CleanupFunc: wireCleanupFunc,
		}, nil
	}

	// Call the Build function with options
	app := Build(wireFunc, WithErrorHandler(errHandler), WithCleanupFunc(optionsCleanupFunc))

	// Verify that the app was created with the correct runnables
	if len(app.runnables) != 1 {
		t.Errorf("Expected app to have 1 runnable, but got %d", len(app.runnables))
	}

	// Verify that the app was created with the correct error handler
	if app.errorHandler == nil {
		t.Error("Expected app to have an error handler, but it doesn't")
	}

	// Verify that the error handler is the one we provided
	if app.errorHandler != nil {
		// Call the error handler to verify it's the one we provided
		app.errorHandler(errors.New("test error"))
		if !errHandlerCalled {
			t.Error("Expected error handler to be called, but it wasn't")
		}
	}

	// Verify that the app was created with the correct cleanup function
	if app.cleanupFunc == nil {
		t.Error("Expected app to have a cleanup function, but it doesn't")
	}

	// Verify that the cleanup function is the one we provided in the options
	if app.cleanupFunc != nil {
		// Call the cleanup function to verify it's the one we provided
		app.cleanupFunc()
		if !optionsCleanupCalled {
			t.Error("Expected options cleanup function to be called, but it wasn't")
		}
	}
}

// TestBuild_Error tests that the Build function handles errors from the wire function
func TestBuild_Error(t *testing.T) {
	// Create an error to return
	testErr := errors.New("test error")

	// Create a mock wire function that returns an error
	wireFunc := func(config struct{}) (WireBundle, error) {
		return WireBundle{}, testErr
	}

	// Call the Build function and expect a panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected Build to panic when wire function returns an error, but it didn't")
		}
	}()

	Build(wireFunc)
}

// TestBuild_ErrorWithHandler tests that the Build function handles errors from the wire function with an error handler
func TestBuild_ErrorWithHandler(t *testing.T) {
	// Create an error to return
	testErr := errors.New("test error")

	// Create a mock cleanup function
	cleanupCalled := false
	cleanupFunc := func() error {
		cleanupCalled = true
		return nil
	}

	// Create a mock wire function that returns an error
	wireFunc := func(config struct{}) (WireBundle, error) {
		return WireBundle{
			CleanupFunc: cleanupFunc,
		}, testErr
	}

	// Create a mock error handler that resolves the error
	errHandlerCalled := false
	errHandler := func(err error) error {
		errHandlerCalled = true
		return nil
	}

	// Call the Build function with the error handler
	app := Build(wireFunc, WithErrorHandler(errHandler))

	// Verify that the error handler was called
	if !errHandlerCalled {
		t.Error("Expected error handler to be called, but it wasn't")
	}

	// Verify that the app was created
	if app == nil {
		t.Error("Expected app to be created, but it wasn't")
	}

	// Verify that the app was created with the correct cleanup function
	if app.cleanupFunc == nil {
		t.Error("Expected app to have a cleanup function, but it doesn't")
	}

	// Verify that the cleanup function is the one from the wire bundle
	if app.cleanupFunc != nil {
		// Call the cleanup function to verify it's the one from the wire bundle
		app.cleanupFunc()
		if !cleanupCalled {
			t.Error("Expected wire bundle cleanup function to be called, but it wasn't")
		}
	}
}

// TestBuild_ErrorWithHandlerNotResolved tests that the Build function handles errors from the wire function with an error handler that doesn't resolve the error
func TestBuild_ErrorWithHandlerNotResolved(t *testing.T) {
	// Create an error to return
	testErr := errors.New("test error")

	// Create a mock wire function that returns an error
	wireFunc := func(config struct{}) (WireBundle, error) {
		return WireBundle{}, testErr
	}

	// Create a mock error handler that doesn't resolve the error
	errHandlerCalled := false
	errHandler := func(err error) error {
		errHandlerCalled = true
		return err
	}

	// Call the Build function with the error handler and expect a panic
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected Build to panic when wire function returns an error and error handler doesn't resolve it, but it didn't")
		}
	}()

	Build(wireFunc, WithErrorHandler(errHandler))

	// Verify that the error handler was called
	if !errHandlerCalled {
		t.Error("Expected error handler to be called, but it wasn't")
	}
}
