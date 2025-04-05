package app

import (
	"context"
	"testing"

	"github.com/pgvanniekerk/ezapp/pkg/ezapp"
)

// TestEnsureEmbedsRunnableStruct tests the EnsureEmbedsRunnableStruct function
func TestEnsureEmbedsRunnableStruct(t *testing.T) {
	// Test with a struct that embeds ezapp.Runnable
	goodRunnable := &GoodRunnable{}
	err := EnsureEmbedsRunnableStruct(goodRunnable)
	if err != nil {
		t.Errorf("Expected no error for GoodRunnable, got %v", err)
	}

	// Test with a struct that doesn't embed ezapp.Runnable
	badRunnable := &BadRunnable{}
	err = EnsureEmbedsRunnableStruct(badRunnable)
	if err == nil {
		t.Errorf("Expected error for BadRunnable, got nil")
	}

	// Test with a non-struct
	nonStruct := NonStruct(func() {})
	err = EnsureEmbedsRunnableStruct(nonStruct)
	if err == nil {
		t.Errorf("Expected error for NonStruct, got nil")
	}
}

// GoodRunnable is a struct that embeds ezapp.Runnable and implements the Runnable interface
type GoodRunnable struct {
	ezapp.Runnable
}

func (g *GoodRunnable) Run() error {
	return nil
}

func (g *GoodRunnable) Stop(ctx context.Context) error {
	return nil
}

// BadRunnable is a struct that doesn't embed ezapp.Runnable but still implements the Runnable interface
type BadRunnable struct{}

func (b *BadRunnable) Run() error {
	return nil
}

func (b *BadRunnable) Stop(ctx context.Context) error {
	return nil
}

// NonStruct is a type that's not a struct but implements the Runnable interface
type NonStruct func()

func (n NonStruct) Run() error {
	return nil
}

func (n NonStruct) Stop(ctx context.Context) error {
	return nil
}
