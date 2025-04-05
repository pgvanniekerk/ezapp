package wire

import (
	"errors"
	"testing"
)

// TestWithShutdownSignal tests the WithShutdownSignal function
func TestWithShutdownSignal(t *testing.T) {
	// Create a channel
	errChan := make(chan error, 1)

	// Call WithShutdownSignal with this channel
	option := WithShutdownSignal(errChan)

	// Create an appOptions struct
	opts := &appOptions{
		shutdownSig: nil, // Initialize with nil
	}

	// Apply the option to the appOptions struct
	option(opts)

	// Check that the shutdownSig field has been set correctly
	if opts.shutdownSig != errChan {
		t.Errorf("Expected shutdownSig to be %v, got %v", errChan, opts.shutdownSig)
	}

	// Test with a different channel
	errChan2 := make(chan error, 2)
	option = WithShutdownSignal(errChan2)
	option(opts)

	// Check that the shutdownSig field has been updated
	if opts.shutdownSig != errChan2 {
		t.Errorf("Expected shutdownSig to be %v, got %v", errChan2, opts.shutdownSig)
	}

	// Test that the channel works as expected
	testErr := errors.New("test error")
	go func() {
		errChan2 <- testErr
	}()

	// Read from the channel
	receivedErr := <-opts.shutdownSig

	// Check that we received the expected error
	if receivedErr != testErr {
		t.Errorf("Expected to receive error %v, got %v", testErr, receivedErr)
	}
}
