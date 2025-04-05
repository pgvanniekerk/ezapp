package buildoption

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

// TestWithoutOptions tests that WithoutOptions returns a BuildOptions with default values
func TestWithoutOptions(t *testing.T) {
	options := WithoutOptions()

	// Check that the returned value is of type *Options
	_, ok := options.(*Options)
	if !ok {
		t.Errorf("WithoutOptions() returned %T, want *Options", options)
	}

	// Check default values
	if options.GetStartupTimeout() != DefaultStartupTimeout {
		t.Errorf("GetStartupTimeout() = %v, want %v", options.GetStartupTimeout(), DefaultStartupTimeout)
	}

	if options.GetEnvVarPrefix() != DefaultEnvVarPrefix {
		t.Errorf("GetEnvVarPrefix() = %v, want %v", options.GetEnvVarPrefix(), DefaultEnvVarPrefix)
	}

	// Check that GetShutdownSignal returns a non-nil channel
	if options.GetShutdownSignal() == nil {
		t.Errorf("GetShutdownSignal() = nil, want non-nil")
	}

	// Check that the error handler is set to DefaultErrorHandler
	// We can't directly compare function values, so we'll check the function pointer
	defaultHandlerPtr := reflect.ValueOf(DefaultErrorHandler).Pointer()
	actualHandlerPtr := reflect.ValueOf(options.GetErrorHandler()).Pointer()
	if defaultHandlerPtr != actualHandlerPtr {
		t.Errorf("GetErrorHandler() points to %v, want %v", actualHandlerPtr, defaultHandlerPtr)
	}
}

// TestWithOptions tests that WithOptions applies the given options
func TestWithOptions(t *testing.T) {
	// Custom values for testing
	customTimeout := 30 * time.Second
	customPrefix := "TEST_"
	customShutdownSignal := make(chan struct{})
	customErrorHandler := func(err error) error { return nil }

	// Create options with custom values
	options := WithOptions(
		WithStartupTimeout(customTimeout),
		WithEnvVarPrefix(customPrefix),
		WithShutdownSignal(customShutdownSignal),
		WithErrorHandler(customErrorHandler),
	)

	// Check that the returned value is of type *Options
	_, ok := options.(*Options)
	if !ok {
		t.Errorf("WithOptions() returned %T, want *Options", options)
	}

	// Check custom values
	if options.GetStartupTimeout() != customTimeout {
		t.Errorf("GetStartupTimeout() = %v, want %v", options.GetStartupTimeout(), customTimeout)
	}

	if options.GetEnvVarPrefix() != customPrefix {
		t.Errorf("GetEnvVarPrefix() = %v, want %v", options.GetEnvVarPrefix(), customPrefix)
	}

	if options.GetShutdownSignal() != customShutdownSignal {
		t.Errorf("GetShutdownSignal() = %v, want %v", options.GetShutdownSignal(), customShutdownSignal)
	}

	// Check that the error handler is set to customErrorHandler
	// We can't directly compare function values, so we'll check the function pointer
	customHandlerPtr := reflect.ValueOf(customErrorHandler).Pointer()
	actualHandlerPtr := reflect.ValueOf(options.GetErrorHandler()).Pointer()
	if customHandlerPtr != actualHandlerPtr {
		t.Errorf("GetErrorHandler() points to %v, want %v", actualHandlerPtr, customHandlerPtr)
	}
}

// TestWithOptionsPartial tests that WithOptions applies only the given options and keeps defaults for others
func TestWithOptionsPartial(t *testing.T) {
	// Custom values for testing
	customTimeout := 30 * time.Second

	// Create options with only custom timeout
	options := WithOptions(
		WithStartupTimeout(customTimeout),
	)

	// Check custom timeout
	if options.GetStartupTimeout() != customTimeout {
		t.Errorf("GetStartupTimeout() = %v, want %v", options.GetStartupTimeout(), customTimeout)
	}

	// Check default values for other options
	if options.GetEnvVarPrefix() != DefaultEnvVarPrefix {
		t.Errorf("GetEnvVarPrefix() = %v, want %v", options.GetEnvVarPrefix(), DefaultEnvVarPrefix)
	}

	// Check that the error handler is set to DefaultErrorHandler
	defaultHandlerPtr := reflect.ValueOf(DefaultErrorHandler).Pointer()
	actualHandlerPtr := reflect.ValueOf(options.GetErrorHandler()).Pointer()
	if defaultHandlerPtr != actualHandlerPtr {
		t.Errorf("GetErrorHandler() points to %v, want %v", actualHandlerPtr, defaultHandlerPtr)
	}
}

// TestGetErrorHandler tests the GetErrorHandler method
func TestGetErrorHandler(t *testing.T) {
	// Create a custom error handler
	customErrorHandler := func(err error) error { return nil }

	// Create options with custom error handler
	options := &Options{
		ErrorHandler: customErrorHandler,
	}

	// Check that GetErrorHandler returns the custom error handler
	if reflect.ValueOf(options.GetErrorHandler()).Pointer() != reflect.ValueOf(customErrorHandler).Pointer() {
		t.Errorf("GetErrorHandler() returned wrong function")
	}
}

// TestGetStartupTimeout tests the GetStartupTimeout method
func TestGetStartupTimeout(t *testing.T) {
	// Create options with custom timeout
	customTimeout := 30 * time.Second
	options := &Options{
		StartupTimeout: customTimeout,
	}

	// Check that GetStartupTimeout returns the custom timeout
	if options.GetStartupTimeout() != customTimeout {
		t.Errorf("GetStartupTimeout() = %v, want %v", options.GetStartupTimeout(), customTimeout)
	}
}

// TestGetEnvVarPrefix tests the GetEnvVarPrefix method
func TestGetEnvVarPrefix(t *testing.T) {
	// Create options with custom prefix
	customPrefix := "TEST_"
	options := &Options{
		EnvVarPrefix: customPrefix,
	}

	// Check that GetEnvVarPrefix returns the custom prefix
	if options.GetEnvVarPrefix() != customPrefix {
		t.Errorf("GetEnvVarPrefix() = %v, want %v", options.GetEnvVarPrefix(), customPrefix)
	}
}

// TestGetShutdownSignal tests the GetShutdownSignal method
func TestGetShutdownSignal(t *testing.T) {
	// Test with custom shutdown signal
	customShutdownSignal := make(chan struct{})
	options := &Options{
		ShutdownSignal: customShutdownSignal,
	}

	// Check that GetShutdownSignal returns the custom shutdown signal
	if options.GetShutdownSignal() != customShutdownSignal {
		t.Errorf("GetShutdownSignal() returned wrong channel")
	}

	// Test with nil shutdown signal (should return DefaultShutdownSignal)
	options = &Options{
		ShutdownSignal: nil,
	}

	// Check that GetShutdownSignal returns a non-nil channel
	if options.GetShutdownSignal() == nil {
		t.Errorf("GetShutdownSignal() = nil, want non-nil")
	}
}

// TestDefaultShutdownSignal tests the DefaultShutdownSignal function
func TestDefaultShutdownSignal(t *testing.T) {
	// Get the default shutdown signal
	shutdownChan := DefaultShutdownSignal()

	// Check that it's not nil
	if shutdownChan == nil {
		t.Errorf("DefaultShutdownSignal() = nil, want non-nil")
	}

	// Test that the channel closes when a signal is received
	// This is a bit tricky to test without actually sending a signal
	// We'll just check that the channel is of the right type
	if reflect.TypeOf(shutdownChan).String() != "<-chan struct {}" {
		t.Errorf("DefaultShutdownSignal() returned channel of type %v, want <-chan struct {}", reflect.TypeOf(shutdownChan))
	}
}

// TestDefaultErrorHandler tests the DefaultErrorHandler function
func TestDefaultErrorHandler(t *testing.T) {
	// DefaultErrorHandler panics, so we need to recover from the panic
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("DefaultErrorHandler did not panic")
		}
	}()

	// Call DefaultErrorHandler with an error
	err := DefaultErrorHandler(errors.New("test error"))
	if err != nil {
		t.Errorf("DefaultErrorHandler returned an error: %v", err)
	}
}

// TestWithErrorHandler tests the WithErrorHandler function
func TestWithErrorHandler(t *testing.T) {
	// Create a custom error handler
	customErrorHandler := func(err error) error { return nil }

	// Create an option with the custom error handler
	option := WithErrorHandler(customErrorHandler)

	// Apply the option to an Options struct
	options := &Options{}
	option(options)

	// Check that the error handler was set correctly
	if reflect.ValueOf(options.ErrorHandler).Pointer() != reflect.ValueOf(customErrorHandler).Pointer() {
		t.Errorf("WithErrorHandler did not set the error handler correctly")
	}
}

// TestWithStartupTimeout tests the WithStartupTimeout function
func TestWithStartupTimeout(t *testing.T) {
	// Create a custom timeout
	customTimeout := 30 * time.Second

	// Create an option with the custom timeout
	option := WithStartupTimeout(customTimeout)

	// Apply the option to an Options struct
	options := &Options{}
	option(options)

	// Check that the timeout was set correctly
	if options.StartupTimeout != customTimeout {
		t.Errorf("WithStartupTimeout() set timeout to %v, want %v", options.StartupTimeout, customTimeout)
	}
}

// TestWithEnvVarPrefix tests the WithEnvVarPrefix function
func TestWithEnvVarPrefix(t *testing.T) {
	// Create a custom prefix
	customPrefix := "TEST_"

	// Create an option with the custom prefix
	option := WithEnvVarPrefix(customPrefix)

	// Apply the option to an Options struct
	options := &Options{}
	option(options)

	// Check that the prefix was set correctly
	if options.EnvVarPrefix != customPrefix {
		t.Errorf("WithEnvVarPrefix() set prefix to %v, want %v", options.EnvVarPrefix, customPrefix)
	}
}

// TestWithShutdownSignal tests the WithShutdownSignal function
func TestWithShutdownSignal(t *testing.T) {
	// Create a custom shutdown signal
	customShutdownSignal := make(chan struct{})

	// Create an option with the custom shutdown signal
	option := WithShutdownSignal(customShutdownSignal)

	// Apply the option to an Options struct
	options := &Options{}
	option(options)

	// Check that the shutdown signal was set correctly
	if options.ShutdownSignal != customShutdownSignal {
		t.Errorf("WithShutdownSignal did not set the shutdown signal correctly")
	}
}
