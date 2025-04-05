package wire

import (
	"log/slog"
	"testing"
)

// TestWithLogAttrs tests the WithLogAttrs function
func TestWithLogAttrs(t *testing.T) {
	// Create some log attributes
	attr1 := slog.String("key1", "value1")
	attr2 := slog.Int("key2", 42)

	// Call WithLogAttrs with these attributes
	option := WithLogAttrs(attr1, attr2)

	// Create an appOptions struct
	opts := &appOptions{
		logAttrs: []slog.Attr{}, // Initialize with empty slice
	}

	// Apply the option to the appOptions struct
	option(opts)

	// Check that the logAttrs field has been set correctly
	if len(opts.logAttrs) != 2 {
		t.Errorf("Expected 2 log attributes, got %d", len(opts.logAttrs))
	}

	// Check the first attribute
	if opts.logAttrs[0].Key != "key1" {
		t.Errorf("Expected key1, got %s", opts.logAttrs[0].Key)
	}
	if opts.logAttrs[0].Value.String() != `value1` {
		t.Errorf("Expected value1, got %s", opts.logAttrs[0].Value.String())
	}

	// Check the second attribute
	if opts.logAttrs[1].Key != "key2" {
		t.Errorf("Expected key2, got %s", opts.logAttrs[1].Key)
	}
	if opts.logAttrs[1].Value.String() != "42" {
		t.Errorf("Expected 42, got %s", opts.logAttrs[1].Value.String())
	}

	// Test with no attributes
	option = WithLogAttrs()
	opts = &appOptions{
		logAttrs: []slog.Attr{attr1, attr2}, // Initialize with some attributes
	}
	option(opts)

	// Check that the logAttrs field has been set to an empty slice
	if len(opts.logAttrs) != 0 {
		t.Errorf("Expected 0 log attributes, got %d", len(opts.logAttrs))
	}
}
