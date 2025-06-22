package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestConfig is a test struct for LoadVar
type TestConfig struct {
	TestString string `env:"TEST_STRING"`
	TestInt    int    `env:"TEST_INT"`
	TestBool   bool   `env:"TEST_BOOL"`
}

func TestLoadVar(t *testing.T) {
	// Test case 1: Successful loading of configuration
	t.Run("successful loading", func(t *testing.T) {
		// Set environment variables
		os.Setenv("TEST_STRING", "test value")
		os.Setenv("TEST_INT", "42")
		os.Setenv("TEST_BOOL", "true")
		defer func() {
			os.Unsetenv("TEST_STRING")
			os.Unsetenv("TEST_INT")
			os.Unsetenv("TEST_BOOL")
		}()

		// Call the function
		config, err := LoadVar[TestConfig]()

		// Check results
		assert.NoError(t, err)
		assert.Equal(t, "test value", config.TestString)
		assert.Equal(t, 42, config.TestInt)
		assert.Equal(t, true, config.TestBool)
	})

	// Test case 2: Invalid environment variable
	t.Run("invalid environment variable", func(t *testing.T) {
		// Set environment variables with invalid value
		os.Setenv("TEST_INT", "not a number")
		defer os.Unsetenv("TEST_INT")

		// Call the function
		_, err := LoadVar[TestConfig]()

		// Check results
		assert.Error(t, err)
	})

	// Test case 3: Non-struct type
	t.Run("non-struct type", func(t *testing.T) {
		// Call the function with a non-struct type
		_, err := LoadVar[string]()

		// Check results
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config type must be a struct")
	})
}