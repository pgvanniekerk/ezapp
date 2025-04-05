package app

import (
	"log/slog"
	"os"
	"testing"
	"time"
)

// TestNew tests the New function
func TestNew(t *testing.T) {
	// Create a logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Create a channel for the shutdown signal
	shutdownSig := make(chan error, 1)

	// Create a good runnable that embeds ezapp.Runnable
	goodRunnable := &GoodRunnable{}

	// Create log attributes
	logAttr1 := slog.String("key1", "value1")
	logAttr2 := slog.Int("key2", 42)

	// Create params with the good runnable
	params := Params{
		ShutdownTimeout: 1 * time.Second,
		Runnables:       []Runnable{goodRunnable},
		ShutdownSig:     shutdownSig,
		Logger:          logger,
		LogAttrs:        []slog.Attr{logAttr1, logAttr2},
	}

	// Call New with these params
	app, err := New(params)

	// Check that there was no error
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check that the app is not nil
	if app == nil {
		t.Errorf("Expected non-nil app, got nil")
	}

	// Check that the app has the expected properties
	if app.shutdownTimeout != params.ShutdownTimeout {
		t.Errorf("Expected shutdownTimeout to be %v, got %v", params.ShutdownTimeout, app.shutdownTimeout)
	}

	if len(app.runnables) != len(params.Runnables) {
		t.Errorf("Expected %d runnables, got %d", len(params.Runnables), len(app.runnables))
	}

	if app.shutdownSig != params.ShutdownSig {
		t.Errorf("Expected shutdownSig to be %v, got %v", params.ShutdownSig, app.shutdownSig)
	}

	// Create a bad runnable that doesn't embed ezapp.Runnable
	badRunnable := &BadRunnable{}

	// Create params with the bad runnable
	params = Params{
		ShutdownTimeout: 1 * time.Second,
		Runnables:       []Runnable{badRunnable},
		ShutdownSig:     shutdownSig,
		Logger:          logger,
	}

	// Call New with these params
	app, err = New(params)

	// Check that there was an error
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// Check that the app is nil
	if app != nil {
		t.Errorf("Expected nil app, got %v", app)
	}
}

// Note: We're reusing the GoodRunnable and BadRunnable types from ensureembedsrunnablestruct_test.go
