package app

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"sync"
	"testing"
	"time"
)

// TestAppRun tests the Run method of the App struct
func TestAppRun(t *testing.T) {
	// Create a logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Create a channel for the shutdown signal
	shutdownSig := make(chan error, 1)

	// Create mock runnables
	mockRunnable1 := &MockRunnable{
		name: "runnable1",
	}
	mockRunnable2 := &MockRunnable{
		name: "runnable2",
	}

	// Create an App with the mock runnables and shutdown signal
	app := &App{
		shutdownTimeout: 1 * time.Second,
		runnables:       []Runnable{mockRunnable1, mockRunnable2},
		shutdownSig:     shutdownSig,
		logger:          logger,
	}

	// Run the app in a goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		app.Run()
	}()

	// Wait a bit for the runnables to start
	time.Sleep(100 * time.Millisecond)

	// Send a shutdown signal
	shutdownSig <- errors.New("shutdown")

	// Wait for the app to finish
	wg.Wait()

	// Check that the runnables were started and stopped
	if !mockRunnable1.runCalled {
		t.Errorf("Expected Run to be called on runnable1, but it wasn't")
	}
	if !mockRunnable2.runCalled {
		t.Errorf("Expected Run to be called on runnable2, but it wasn't")
	}
	if !mockRunnable1.stopCalled {
		t.Errorf("Expected Stop to be called on runnable1, but it wasn't")
	}
	if !mockRunnable2.stopCalled {
		t.Errorf("Expected Stop to be called on runnable2, but it wasn't")
	}
}

// MockRunnable is a mock implementation of Runnable for testing
type MockRunnable struct {
	name       string
	runCalled  bool
	stopCalled bool
}

func (m *MockRunnable) Run() error {
	m.runCalled = true
	return nil
}

func (m *MockRunnable) Stop(ctx context.Context) error {
	m.stopCalled = true
	return nil
}
