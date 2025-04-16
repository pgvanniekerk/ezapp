package wire

import (
	"errors"
	"testing"
)

// TestWithCriticalErrHandler tests the WithCriticalErrHandler function
func TestWithCriticalErrHandler(t *testing.T) {
	// Create a critical error handler function
	var handlerCalled bool
	var receivedErr error
	handler := func(err error) {
		handlerCalled = true
		receivedErr = err
	}

	// Call WithCriticalErrHandler with this handler
	option := WithCriticalErrHandler(handler)

	// Create an appOptions struct
	opts := &appOptions{
		criticalErrHandler: nil, // Initialize with nil
	}

	// Apply the option to the appOptions struct
	option(opts)

	// Check that the criticalErrHandler field has been set correctly
	if opts.criticalErrHandler == nil {
		t.Error("Expected criticalErrHandler to be set, got nil")
	}

	// Test the handler by calling it with an error
	testErr := errors.New("test error")
	opts.criticalErrHandler(testErr)

	// Check that the handler was called with the correct error
	if !handlerCalled {
		t.Error("Expected handler to be called, but it wasn't")
	}
	if !errors.Is(testErr, receivedErr) {
		t.Errorf("Expected handler to receive error %v, got %v", testErr, receivedErr)
	}

	// Test with a different handler
	var handler2Called bool
	var receivedErr2 error
	handler2 := func(err error) {
		handler2Called = true
		receivedErr2 = err
	}
	option = WithCriticalErrHandler(handler2)
	option(opts)

	// Check that the criticalErrHandler field has been updated
	if opts.criticalErrHandler == nil {
		t.Error("Expected criticalErrHandler to be set, got nil")
	}

	// Reset the flags
	handlerCalled = false
	receivedErr = nil

	// Test the new handler by calling it with a different error
	testErr2 := errors.New("another test error")
	opts.criticalErrHandler(testErr2)

	// Check that the new handler was called with the correct error
	if handlerCalled {
		t.Error("Expected original handler not to be called, but it was")
	}
	if !handler2Called {
		t.Error("Expected new handler to be called, but it wasn't")
	}
	if !errors.Is(testErr2, receivedErr2) {
		t.Errorf("Expected new handler to receive error %v, got %v", testErr2, receivedErr2)
	}
}
