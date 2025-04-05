package app

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/pgvanniekerk/ezapp/pkg/ezapp"
)

// TestSetRunnableLogger tests the setRunnableLogger function
func TestSetRunnableLogger(t *testing.T) {
	// Create a logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Create a runnable with the ezapp.Runnable struct embedded with the toggle:"useEzAppLogger" tag
	runnable := &RunnableWithLogger{}

	// Call setRunnableLogger
	setRunnableLogger(runnable, logger)

	// Check that the Logger field has been set
	if runnable.Logger == nil {
		t.Errorf("Expected Logger field to be set, but it's nil")
	}

	// Create a runnable with the ezapp.Runnable struct embedded but without the toggle:"useEzAppLogger" tag
	runnableWithoutTag := &RunnableWithoutTag{}

	// Call setRunnableLogger
	setRunnableLogger(runnableWithoutTag, logger)

	// Check that the Logger field has not been set
	if runnableWithoutTag.Logger != nil {
		t.Errorf("Expected Logger field to not be set, but it's %v", runnableWithoutTag.Logger)
	}

	// Create a runnable that doesn't embed the ezapp.Runnable struct
	runnableWithoutEmbedding := &RunnableWithoutEmbedding{}

	// Call setRunnableLogger
	setRunnableLogger(runnableWithoutEmbedding, logger)

	// No need to check anything here, as the function should just return without error

	// Create a non-struct runnable
	nonStructRunnable := NonStructRunnable(func() {})

	// Call setRunnableLogger
	setRunnableLogger(nonStructRunnable, logger)

	// No need to check anything here, as the function should just return without error
}

// RunnableWithLogger is a struct that embeds ezapp.Runnable with the toggle:"useEzAppLogger" tag
type RunnableWithLogger struct {
	ezapp.Runnable `toggle:"useEzAppLogger"`
}

func (r *RunnableWithLogger) Run() error {
	return nil
}

func (r *RunnableWithLogger) Stop(ctx context.Context) error {
	return nil
}

// RunnableWithoutTag is a struct that embeds ezapp.Runnable but without the toggle:"useEzAppLogger" tag
type RunnableWithoutTag struct {
	ezapp.Runnable
}

func (r *RunnableWithoutTag) Run() error {
	return nil
}

func (r *RunnableWithoutTag) Stop(ctx context.Context) error {
	return nil
}

// RunnableWithoutEmbedding is a struct that doesn't embed ezapp.Runnable
type RunnableWithoutEmbedding struct{}

func (r *RunnableWithoutEmbedding) Run() error {
	return nil
}

func (r *RunnableWithoutEmbedding) Stop(ctx context.Context) error {
	return nil
}

// NonStructRunnable is a type that's not a struct but implements the Runnable interface
type NonStructRunnable func()

func (n NonStructRunnable) Run() error {
	return nil
}

func (n NonStructRunnable) Stop(ctx context.Context) error {
	return nil
}
