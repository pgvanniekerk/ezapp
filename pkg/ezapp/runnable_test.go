package ezapp

import (
	"context"
	"log/slog"
	"os"
	"testing"
)

// TestRunnable tests the Runnable struct and its methods
func TestRunnable(t *testing.T) {
	// Create a logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Create a Runnable with the logger
	runnable := Runnable{
		Logger: logger,
	}

	// Test the Run method
	err := runnable.Run()
	if err != nil {
		t.Errorf("Expected nil error from Run, got %v", err)
	}

	// Create a context
	ctx := context.Background()

	// Test the Stop method
	err = runnable.Stop(ctx)
	if err != nil {
		t.Errorf("Expected nil error from Stop, got %v", err)
	}

	// Test with a nil logger
	runnable = Runnable{
		Logger: nil,
	}

	// Test the Run method with nil logger
	err = runnable.Run()
	if err != nil {
		t.Errorf("Expected nil error from Run with nil logger, got %v", err)
	}

	// Test the Stop method with nil logger
	err = runnable.Stop(ctx)
	if err != nil {
		t.Errorf("Expected nil error from Stop with nil logger, got %v", err)
	}
}
