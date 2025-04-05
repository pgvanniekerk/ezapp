package app

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

// MockService is a mock implementation of the Service interface for testing
type MockService struct {
	RunCalled  bool
	StopCalled bool
	RunError   error
	StopError  error
	// Channel to block Run method until test signals it to return
	RunBlock chan struct{}
	// Channel to signal when Run is called
	RunStarted chan struct{}
	// Channel to signal when Stop is called
	StopSignal chan struct{}
	// Delay before Run returns
	RunDelay time.Duration
	// Delay before Stop returns
	StopDelay time.Duration
}

func NewMockService() *MockService {
	return &MockService{
		RunBlock:   make(chan struct{}),
		RunStarted: make(chan struct{}, 1),
		StopSignal: make(chan struct{}, 1),
	}
}

func (m *MockService) Run() error {
	m.RunCalled = true
	// Signal that Run was called
	select {
	case m.RunStarted <- struct{}{}:
	default:
	}

	// Simulate delay if specified
	if m.RunDelay > 0 {
		time.Sleep(m.RunDelay)
	}

	// If RunBlock channel is not nil, wait for it to be closed or receive a value
	if m.RunBlock != nil {
		<-m.RunBlock
	}

	return m.RunError
}

func (m *MockService) Stop(ctx context.Context) error {
	m.StopCalled = true
	// Signal that Stop was called
	select {
	case m.StopSignal <- struct{}{}:
	default:
	}

	// Simulate delay if specified
	if m.StopDelay > 0 {
		select {
		case <-time.After(m.StopDelay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return m.StopError
}

// ErrorHandlerMock is a mock implementation of the ErrorHandler function
type ErrorHandlerMock struct {
	Called    bool
	CalledErr error
	ReturnErr error
	mu        sync.Mutex
}

func (e *ErrorHandlerMock) Handle(err error) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.Called = true
	e.CalledErr = err
	return e.ReturnErr
}

// TestNewApp tests the creation of a new App instance with both default and custom configurations.
// It verifies that the App struct is properly initialized with the expected values for its fields,
// including services, shutdownTimeout, shutdownSig, and errorHandler.
func TestNewApp(t *testing.T) {
	// This test verifies that when creating an App with default values (nil errorHandler),
	// the App is properly initialized with the expected default values:
	// - services are correctly set
	// - shutdownTimeout is set to the default 15 seconds
	// - shutdownSig is correctly set
	// - a default errorHandler is provided
	t.Run("with default values", func(t *testing.T) {
		services := []Service{NewMockService()}
		shutdownSig := make(chan struct{})

		app, err := NewApp(services, nil, shutdownSig)
		if err != nil {
			t.Fatalf("NewApp returned an error: %v", err)
		}

		if app.services == nil || len(app.services) != 1 {
			t.Errorf("Expected services to be set, got %v", app.services)
		}

		if app.shutdownTimeout != 15*time.Second {
			t.Errorf("Expected shutdownTimeout to be 15s, got %v", app.shutdownTimeout)
		}

		if app.shutdownSig != shutdownSig {
			t.Errorf("Expected shutdownSig to be set correctly")
		}

		if app.errorHandler == nil {
			t.Errorf("Expected default errorHandler to be set")
		}
	})

	// This test verifies that when creating an App with a custom error handler,
	// the App is properly initialized with that custom handler and the handler
	// functions correctly when invoked with an error.
	t.Run("with custom error handler", func(t *testing.T) {
		services := []Service{NewMockService()}
		shutdownSig := make(chan struct{})
		customErrorHandler := func(err error) error {
			return errors.New("custom error")
		}

		app, err := NewApp(services, customErrorHandler, shutdownSig)
		if err != nil {
			t.Fatalf("NewApp returned an error: %v", err)
		}

		if app.errorHandler == nil {
			t.Errorf("Expected custom errorHandler to be set")
		}

		// Test that the custom error handler is used
		result := app.errorHandler(errors.New("test"))
		if result.Error() != "custom error" {
			t.Errorf("Custom error handler not working as expected")
		}
	})

	// Test validation of inputs
	t.Run("with nil services", func(t *testing.T) {
		shutdownSig := make(chan struct{})
		_, err := NewApp(nil, nil, shutdownSig)
		if err == nil {
			t.Errorf("Expected error when services is nil")
		}
	})

	t.Run("with empty services", func(t *testing.T) {
		services := []Service{}
		shutdownSig := make(chan struct{})
		_, err := NewApp(services, nil, shutdownSig)
		if err == nil {
			t.Errorf("Expected error when services is empty")
		}
	})

	t.Run("with nil shutdownSig", func(t *testing.T) {
		services := []Service{NewMockService()}
		_, err := NewApp(services, nil, nil)
		if err == nil {
			t.Errorf("Expected error when shutdownSig is nil")
		}
	})
}

// TestApp_Run tests the Run method of the App struct under various scenarios.
// It verifies that the App correctly starts and stops services, handles errors from services,
// manages timeouts during shutdown, and coordinates multiple services.
func TestApp_Run(t *testing.T) {
	// This test verifies the normal operation and graceful shutdown of the App.
	// It checks that:
	// - The service's Run method is called when the App starts
	// - The App responds to the shutdown signal
	// - The service's Stop method is called during shutdown
	// - The App waits for all services to finish before returning
	t.Run("normal operation and shutdown", func(t *testing.T) {
		mockService := NewMockService()
		services := []Service{mockService}
		shutdownSig := make(chan struct{})

		app, err := NewApp(services, nil, shutdownSig)
		if err != nil {
			t.Fatalf("NewApp returned an error: %v", err)
		}

		// Start the app in a goroutine
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			app.Run()
		}()

		// Wait for the service to start
		select {
		case <-mockService.RunStarted:
			// Service started successfully
		case <-time.After(100 * time.Millisecond):
			t.Fatalf("Timed out waiting for service to start")
		}

		// Trigger shutdown
		close(shutdownSig)

		// Allow the Run method to return
		close(mockService.RunBlock)

		// Wait for the app to finish
		wg.Wait()

		// Verify that the service was started and stopped
		if !mockService.RunCalled {
			t.Errorf("Expected service Run to be called")
		}

		if !mockService.StopCalled {
			t.Errorf("Expected service Stop to be called")
		}
	})

	// This test verifies that the App properly handles errors from services.
	// It checks that:
	// - When a service returns an error from its Run method, the error handler is called with that error
	// - The error handler receives the correct error message
	// - The service's Stop method is still called even after an error occurs
	t.Run("service error handling", func(t *testing.T) {
		mockService := NewMockService()
		mockService.RunError = errors.New("service error")

		services := []Service{mockService}
		shutdownSig := make(chan struct{})

		errorHandler := &ErrorHandlerMock{}
		app, err := NewApp(services, errorHandler.Handle, shutdownSig)
		if err != nil {
			t.Fatalf("NewApp returned an error: %v", err)
		}

		// Close the RunBlock channel to allow the Run method to return immediately with the error
		close(mockService.RunBlock)

		// Start the app
		app.Run()

		// Verify that the error handler was called with the correct error
		if !errorHandler.Called {
			t.Errorf("Expected error handler to be called")
		}

		if errorHandler.CalledErr == nil || errorHandler.CalledErr.Error() != "service error" {
			t.Errorf("Expected error handler to be called with 'service error', got %v", errorHandler.CalledErr)
		}

		// Verify that the service was stopped
		if !mockService.StopCalled {
			t.Errorf("Expected service Stop to be called after error")
		}
	})

	// This test verifies that the App properly handles errors that occur during service shutdown.
	// It checks that:
	// - When a service returns an error from its Stop method, the error handler is called with that error
	// - The error handler receives the correct error message
	// - The App continues with the shutdown process despite the error
	t.Run("stop error handling", func(t *testing.T) {
		mockService := NewMockService()
		mockService.StopError = errors.New("stop error")

		services := []Service{mockService}
		shutdownSig := make(chan struct{})

		errorHandler := &ErrorHandlerMock{}
		app, err := NewApp(services, errorHandler.Handle, shutdownSig)
		if err != nil {
			t.Fatalf("NewApp returned an error: %v", err)
		}

		// Start the app in a goroutine
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			app.Run()
		}()

		// Wait for the service to start
		select {
		case <-mockService.RunStarted:
			// Service started successfully
		case <-time.After(100 * time.Millisecond):
			t.Fatalf("Timed out waiting for service to start")
		}

		// Trigger shutdown
		close(shutdownSig)

		// Allow the Run method to return
		close(mockService.RunBlock)

		// Wait for the app to finish
		wg.Wait()

		// Verify that the error handler was called with the stop error
		if !errorHandler.Called {
			t.Errorf("Expected error handler to be called for stop error")
		}

		if errorHandler.CalledErr == nil || errorHandler.CalledErr.Error() != "stop error" {
			t.Errorf("Expected error handler to be called with 'stop error', got %v", errorHandler.CalledErr)
		}
	})

	// This test verifies that the App properly handles timeouts during service shutdown.
	// It checks that:
	// - When a service takes longer to stop than the shutdown timeout, the context deadline is exceeded
	// - The App handles the timeout gracefully without deadlocking
	// - The App continues with the shutdown process despite the timeout
	t.Run("shutdown timeout", func(t *testing.T) {
		mockService := NewMockService()
		mockService.StopDelay = 50 * time.Millisecond

		services := []Service{mockService}
		shutdownSig := make(chan struct{})

		// Create a custom error handler that doesn't panic on context.DeadlineExceeded
		errorHandler := func(err error) error {
			// Just return the error without panicking
			return err
		}

		app, err := NewApp(services, errorHandler, shutdownSig)
		if err != nil {
			t.Fatalf("NewApp returned an error: %v", err)
		}
		// Set a very short shutdown timeout to test timeout handling
		app.shutdownTimeout = 10 * time.Millisecond

		// Start the app in a goroutine
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			app.Run()
		}()

		// Wait for the service to start
		select {
		case <-mockService.RunStarted:
			// Service started successfully
		case <-time.After(100 * time.Millisecond):
			t.Fatalf("Timed out waiting for service to start")
		}

		// Trigger shutdown
		close(shutdownSig)

		// Allow the Run method to return
		close(mockService.RunBlock)

		// Wait for the app to finish
		wg.Wait()

		// The test passes if we reach here without deadlocking
		// The context timeout should have been triggered
	})

	// This test verifies that the App correctly manages multiple services.
	// It checks that:
	// - All services' Run methods are called when the App starts
	// - The App waits for all services to start
	// - All services' Stop methods are called during shutdown
	// - The App waits for all services to finish before returning
	t.Run("multiple services", func(t *testing.T) {
		mockService1 := NewMockService()
		mockService2 := NewMockService()

		services := []Service{mockService1, mockService2}
		shutdownSig := make(chan struct{})

		app, err := NewApp(services, nil, shutdownSig)
		if err != nil {
			t.Fatalf("NewApp returned an error: %v", err)
		}

		// Start the app in a goroutine
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			app.Run()
		}()

		// Wait for both services to start
		select {
		case <-mockService1.RunStarted:
			// Service 1 started
		case <-time.After(100 * time.Millisecond):
			t.Fatalf("Timed out waiting for service 1 to start")
		}

		select {
		case <-mockService2.RunStarted:
			// Service 2 started
		case <-time.After(100 * time.Millisecond):
			t.Fatalf("Timed out waiting for service 2 to start")
		}

		// Trigger shutdown
		close(shutdownSig)

		// Allow the Run methods to return
		close(mockService1.RunBlock)
		close(mockService2.RunBlock)

		// Wait for the app to finish
		wg.Wait()

		// Verify that both services were started and stopped
		if !mockService1.RunCalled || !mockService2.RunCalled {
			t.Errorf("Expected both services' Run methods to be called")
		}

		if !mockService1.StopCalled || !mockService2.StopCalled {
			t.Errorf("Expected both services' Stop methods to be called")
		}
	})
}
