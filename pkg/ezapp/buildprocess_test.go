package ezapp

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/dig"
)

func TestBuildProcess(t *testing.T) {
	// Create a mock BuildContext
	mockContainer := dig.New()
	mockBuildContext := &buildContext{
		container: mockContainer,
	}

	// Create a BuildProcess that provides a string
	buildProc := func(bCtx BuildContext) error {
		return bCtx.Container().Provide(func() string {
			return "test string"
		})
	}

	// Apply the BuildProcess to the mock BuildContext
	err := buildProc(mockBuildContext)
	assert.NoError(t, err, "Applying a valid BuildProcess should not error")

	// Verify that the string was provided to the container
	err = mockContainer.Invoke(func(s string) {
		assert.Equal(t, "test string", s, "The provided string should match the original string")
	})
	assert.NoError(t, err, "Invoking the container with the provided type should not error")
}

func TestBuildProcessError(t *testing.T) {
	// Create a mock BuildContext
	mockContainer := dig.New()
	mockBuildContext := &buildContext{
		container: mockContainer,
	}

	// Create a BuildProcess that returns an error
	expectedErr := errors.New("test error")
	buildProc := func(bCtx BuildContext) error {
		return expectedErr
	}

	// Apply the BuildProcess to the mock BuildContext
	err := buildProc(mockBuildContext)
	assert.Equal(t, expectedErr, err, "The BuildProcess should return the expected error")
}

func TestBuildProcessComposition(t *testing.T) {
	// Create a mock BuildContext
	mockContainer := dig.New()
	mockBuildContext := &buildContext{
		container: mockContainer,
	}

	// Create a BuildProcess that provides a string
	buildProc1 := func(bCtx BuildContext) error {
		return bCtx.Container().Provide(func() string {
			return "test string 1"
		}, dig.Name("string1"))
	}

	// Create another BuildProcess that provides a different string
	buildProc2 := func(bCtx BuildContext) error {
		return bCtx.Container().Provide(func() string {
			return "test string 2"
		}, dig.Name("string2"))
	}

	// Apply both BuildProcesses to the mock BuildContext
	err := buildProc1(mockBuildContext)
	assert.NoError(t, err, "Applying the first BuildProcess should not error")
	err = buildProc2(mockBuildContext)
	assert.NoError(t, err, "Applying the second BuildProcess should not error")

	// Define a struct with fields that have name tags
	type NamedDeps struct {
		dig.In
		String1 string `name:"string1"`
		String2 string `name:"string2"`
	}

	// Verify that both strings were provided to the container
	err = mockContainer.Invoke(func(deps NamedDeps) {
		assert.Equal(t, "test string 1", deps.String1, "The first provided string should match the original string")
		assert.Equal(t, "test string 2", deps.String2, "The second provided string should match the original string")
	})
	assert.NoError(t, err, "Invoking the container with the provided types should not error")
}
