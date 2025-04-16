package ezapp

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestConfig is a test configuration struct
type TestConfig struct {
	StringValue string  `env:"TEST_STRING"`
	IntValue    int     `env:"TEST_INT"`
	BoolValue   bool    `env:"TEST_BOOL"`
	FloatValue  float64 `env:"TEST_FLOAT"`
	DefaultName string  `env:"DEFAULTNAME"` // Will use environment variable "DEFAULTNAME"
}

// TestBuild tests the Build function
func TestBuild(t *testing.T) {

	// Set environment variables for the test
	os.Setenv("TEST_STRING", "test-value")
	os.Setenv("TEST_INT", "42")
	os.Setenv("TEST_BOOL", "true")
	os.Setenv("TEST_FLOAT", "3.14")
	os.Setenv("DEFAULTNAME", "default-name")

	defer func() {
		os.Unsetenv("TEST_STRING")
		os.Unsetenv("TEST_INT")
		os.Unsetenv("TEST_BOOL")
		os.Unsetenv("TEST_FLOAT")
		os.Unsetenv("DEFAULTNAME")
	}()

	// Create a mock builder that captures the config
	var capturedConfig TestConfig
	mockBuilder := func(conf TestConfig) ([]Runnable, error) {
		capturedConfig = conf
		return []Runnable{
			mockRunnable{
				runFunc: func(ctx context.Context) error {
					return nil
				},
			},
		}, nil
	}

	// Call the Build function
	app := Build(mockBuilder)

	// Assert that there was no error
	assert.Nil(t, app.initErr)
	assert.NotNil(t, app)

	// Assert that the config was populated correctly
	assert.Equal(t, "test-value", capturedConfig.StringValue)
	assert.Equal(t, 42, capturedConfig.IntValue)
	assert.Equal(t, true, capturedConfig.BoolValue)
	assert.Equal(t, 3.14, capturedConfig.FloatValue)
	assert.Equal(t, "default-name", capturedConfig.DefaultName)
}

// TestBuildNonStruct tests that Build returns an error when CONF is not a struct
func TestBuildNonStruct(t *testing.T) {
	// Create a mock builder that takes a non-struct type
	mockBuilder := func(conf string) ([]Runnable, error) {
		return nil, nil
	}

	// Call the Build function
	app := Build(mockBuilder)

	// Assert that there was an error
	assert.NotNil(t, app.initErr)
	assert.NotNil(t, app)
	assert.Contains(t, app.initErr.Error(), "CONF must be a struct")
}

// TestBuildEnvUnmarshalError tests that Build returns an error when environment variable unmarshaling fails
func TestBuildEnvUnmarshalError(t *testing.T) {

	// Define a test config with a field that will cause unmarshaling to fail
	type InvalidTestConfig struct {
		IntValue int `env:"TEST_INT"`
	}

	// Set environment variable with an invalid value for an int
	os.Setenv("TEST_INT", "not-an-int")
	defer func() {
		os.Unsetenv("TEST_INT")
	}()

	// Create a mock builder
	mockBuilder := func(conf InvalidTestConfig) ([]Runnable, error) {
		return []Runnable{
			mockRunnable{
				runFunc: func(ctx context.Context) error {
					return nil
				},
			},
		}, nil
	}

	// Call the Build function
	app := Build(mockBuilder)

	// Assert that there was an error
	assert.NotNil(t, app.initErr)
	assert.NotNil(t, app)
	assert.Contains(t, app.initErr.Error(), "failed to parse environment variables into CONF")
}

// TestBuildBuilderError tests that Build returns an error when the builder function returns an error
func TestBuildBuilderError(t *testing.T) {
	// Set environment variables for the test
	os.Setenv("TEST_STRING", "test-value")
	defer func() {
		os.Unsetenv("TEST_STRING")
	}()

	// Create a mock builder that returns an error
	expectedError := errors.New("builder error")
	mockBuilder := func(conf TestConfig) ([]Runnable, error) {
		return nil, expectedError
	}

	// Call the Build function
	app := Build(mockBuilder)

	// Assert that there was an error
	assert.NotNil(t, app.initErr)
	assert.NotNil(t, app)
	assert.Equal(t, expectedError, app.initErr)
}
