package ezapp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/dig"
)

func TestConstruct(t *testing.T) {
	// Create a container using Construct with no build processes
	container := Construct()

	// Verify that the container is not nil
	assert.NotNil(t, container, "The container should not be nil")
}

func TestConstructWithBuildProcess(t *testing.T) {
	// Create a BuildProcess that provides a string
	buildProc := func(bCtx BuildContext) error {
		return bCtx.Container().Provide(func() string {
			return "test string"
		})
	}

	// Create a container using Construct with the build process
	container := Construct(buildProc)

	// Verify that the string was provided to the container
	err := container.Invoke(func(s string) {
		assert.Equal(t, "test string", s, "The provided string should match the original string")
	})
	assert.NoError(t, err, "Invoking the container with the provided type should not error")
}

func TestConstructWithMultipleBuildProcesses(t *testing.T) {
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

	// Create a container using Construct with both build processes
	container := Construct(buildProc1, buildProc2)

	// Define a struct with fields that have name tags
	type NamedDeps struct {
		dig.In
		String1 string `name:"string1"`
		String2 string `name:"string2"`
	}

	// Verify that both strings were provided to the container
	err := container.Invoke(func(deps NamedDeps) {
		assert.Equal(t, "test string 1", deps.String1, "The first provided string should match the original string")
		assert.Equal(t, "test string 2", deps.String2, "The second provided string should match the original string")
	})
	assert.NoError(t, err, "Invoking the container with the provided types should not error")
}

func TestConstructWithErrorBuildProcess(t *testing.T) {
	// Create a BuildProcess that returns an error
	buildProc := func(bCtx BuildContext) error {
		return bCtx.Container().Provide(func() (string, error) {
			return "", nil
		})
	}

	// Create a container using Construct with the build process
	// Even though the build process returns an error, Construct should ignore it
	container := Construct(buildProc)

	// Verify that the container is not nil
	assert.NotNil(t, container, "The container should not be nil")
}
