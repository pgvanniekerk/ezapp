package ezapp

import (
	"context"
	"testing"
)

// TestRun tests the Run function
func TestRun(t *testing.T) {
	// Create a mock app
	mockApp := &MockApp{
		t: t,
	}

	// Create a mock builder that returns the mock app
	mockBuilder := func(ctx context.Context, conf struct{}) (*MockApp, error) {
		// Check that the context is not nil
		if ctx == nil {
			t.Errorf("Expected non-nil context, got nil")
		}

		return mockApp, nil
	}

	// Call Run with the mock builder
	// We need to use defer/recover to catch the panic that would occur if Run fails
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Run panicked: %v", r)
		}

		// Check that the Run method was called on the mock app
		if !mockApp.runCalled {
			t.Errorf("Expected Run method to be called on the mock app, but it wasn't")
		}
	}()

	Run(mockBuilder)
}

// MockApp is a mock implementation of EzApp for testing
type MockApp struct {
	t         *testing.T
	runCalled bool
}

func (m *MockApp) Run() {
	m.runCalled = true
}
