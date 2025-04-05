package ezapp

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/pgvanniekerk/ezapp/internal/app"
)

// ==================== Mock Implementations ====================

// mockBuildOptions is a mock implementation of the BuildOptions interface for testing
type mockBuildOptions struct {
	errorHandler   app.ErrorHandler
	startupTimeout time.Duration
	envVarPrefix   string
	shutdownSignal chan struct{}
}

func (m *mockBuildOptions) GetErrorHandler() app.ErrorHandler {
	return m.errorHandler
}

func (m *mockBuildOptions) GetStartupTimeout() time.Duration {
	return m.startupTimeout
}

func (m *mockBuildOptions) GetEnvVarPrefix() string {
	return m.envVarPrefix
}

func (m *mockBuildOptions) GetShutdownSignal() <-chan struct{} {
	return m.shutdownSignal
}

// mockServiceForBuild is a mock implementation of the app.Service interface for testing
type mockServiceForBuild struct {
	runCalled  bool
	stopCalled bool
	runError   error
	stopError  error
}

func (m *mockServiceForBuild) Run() error {
	m.runCalled = true
	return m.runError
}

func (m *mockServiceForBuild) Stop(ctx context.Context) error {
	m.stopCalled = true
	return m.stopError
}

// mockConfig is a mock configuration struct for testing
type mockConfig struct {
	TestValue string `envconfig:"TEST_VALUE" default:"default_value"`
}

// ==================== Test Build Function ====================

// TestBuild tests the Build function with valid options
func TestBuild(t *testing.T) {
	// Create a mock BuildOptions
	shutdownChan := make(chan struct{})
	mockOptions := &mockBuildOptions{
		errorHandler:   func(err error) error { return err },
		startupTimeout: 5 * time.Second,
		envVarPrefix:   "TEST",
		shutdownSignal: shutdownChan,
	}

	// Create a mock WireFunc
	mockService := &mockServiceForBuild{}
	mockWireFunc := func(ctx context.Context, config mockConfig) (ServiceSet, error) {
		return ServiceSet{
			Services: []app.Service{mockService},
		}, nil
	}

	// Call Build with the mock options and WireFunc
	ezApp := Build(mockWireFunc, mockOptions)

	// Check that ezApp is not nil
	if ezApp == nil {
		t.Error("Build returned nil")
	}
}

// TestBuildWithNilOptions tests that Build panics when options is nil
func TestBuildWithNilOptions(t *testing.T) {
	// Create a mock WireFunc
	mockWireFunc := func(ctx context.Context, config mockConfig) (ServiceSet, error) {
		return ServiceSet{}, nil
	}

	// Expect a panic when options is nil
	defer func() {
		if r := recover(); r == nil {
			t.Error("Build did not panic with nil options")
		}
	}()

	// Call Build with nil options
	_ = Build(mockWireFunc, nil)
}

// ==================== Test invokeWireFunc Function ====================

// TestInvokeWireFunc tests the invokeWireFunc function with valid parameters
func TestInvokeWireFunc(t *testing.T) {
	// Create a mock WireFunc
	mockService := &mockServiceForBuild{}
	mockWireFunc := func(ctx context.Context, config mockConfig) (ServiceSet, error) {
		// Verify that the context has a deadline
		_, hasDeadline := ctx.Deadline()
		if !hasDeadline {
			t.Error("Context does not have a deadline")
		}

		return ServiceSet{
			Services: []app.Service{mockService},
		}, nil
	}

	// Call invokeWireFunc with the mock WireFunc
	serviceSet := invokeWireFunc(mockWireFunc, 5*time.Second, "TEST")

	// Check that serviceSet has the expected service
	if len(serviceSet.Services) != 1 {
		t.Errorf("invokeWireFunc returned ServiceSet with %d services, expected 1", len(serviceSet.Services))
	}
}

// TestInvokeWireFuncWithError tests that invokeWireFunc panics when WireFunc returns an error
func TestInvokeWireFuncWithError(t *testing.T) {
	// Create a mock WireFunc that returns an error
	expectedErr := errors.New("wire func error")
	mockWireFunc := func(ctx context.Context, config mockConfig) (ServiceSet, error) {
		return ServiceSet{}, expectedErr
	}

	// Expect a panic when WireFunc returns an error
	defer func() {
		if r := recover(); r == nil {
			t.Error("invokeWireFunc did not panic when WireFunc returned an error")
		}
	}()

	// Call invokeWireFunc with the mock WireFunc
	_ = invokeWireFunc(mockWireFunc, 5*time.Second, "TEST")
}
