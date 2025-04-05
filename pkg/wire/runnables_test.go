package wire

import (
	"context"
	"github.com/pgvanniekerk/ezapp/pkg/ezapp"
	"testing"
)

// TestRunnables tests the Runnables function
func TestRunnables(t *testing.T) {
	// Create a mock Runnable
	mockRunnable := &MockRunnable{}

	// Call Runnables with the mock
	runnablesFunc := Runnables(mockRunnable)

	// Get the runnables from the function
	runnables := runnablesFunc()

	// Check that we got the expected number of runnables
	if len(runnables) != 1 {
		t.Errorf("Expected 1 runnable, got %d", len(runnables))
	}

	// Check that the runnable is the one we passed in
	if _, ok := runnables[0].(*MockRunnable); !ok {
		t.Errorf("Expected runnable to be of type *MockRunnable, got %T", runnables[0])
	}

	// Test with multiple runnables
	mockRunnable2 := &MockRunnable{}
	mockRunnable3 := &MockRunnable{}
	runnablesFunc = Runnables(mockRunnable, mockRunnable2, mockRunnable3)
	runnables = runnablesFunc()

	// Check that we got the expected number of runnables
	if len(runnables) != 3 {
		t.Errorf("Expected 3 runnables, got %d", len(runnables))
	}

	// Check that the runnables are of the expected type
	for i, r := range runnables {
		if _, ok := r.(*MockRunnable); !ok {
			t.Errorf("Expected runnable at index %d to be of type *MockRunnable, got %T", i, r)
		}
	}
}

// MockRunnable is a mock implementation of app.Runnable for testing
type MockRunnable struct {
	ezapp.Runnable
}

func (m *MockRunnable) Run() error {
	return nil
}

func (m *MockRunnable) Stop(ctx context.Context) error {
	return nil
}
