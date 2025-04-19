package ezapp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/dig"
)

func TestBuildContext(t *testing.T) {
	// Create a dig container
	container := dig.New()

	// Create a buildContext with the container
	bCtx := &buildContext{
		container: container,
	}

	// Verify that the Container method returns the correct container
	assert.Equal(t, container, bCtx.Container(), "The Container method should return the container that was provided to the buildContext")
}

func TestBuildContextInterface(t *testing.T) {
	// Create a dig container
	container := dig.New()

	// Create a buildContext with the container
	bCtx := &buildContext{
		container: container,
	}

	// Verify that the buildContext implements the BuildContext interface
	var _ BuildContext = bCtx

	// Use the BuildContext interface to access the container
	var bCtxInterface BuildContext = bCtx
	assert.Equal(t, container, bCtxInterface.Container(), "The Container method should return the container that was provided to the buildContext")
}

func TestBuildContextWithContainer(t *testing.T) {
	// Create a dig container
	container := dig.New()

	// Provide a string to the container
	err := container.Provide(func() string {
		return "test string"
	})
	assert.NoError(t, err, "Providing a string to the container should not error")

	// Create a buildContext with the container
	bCtx := &buildContext{
		container: container,
	}

	// Verify that the Container method returns the container with the provided string
	returnedContainer := bCtx.Container()
	err = returnedContainer.Invoke(func(s string) {
		assert.Equal(t, "test string", s, "The provided string should match the original string")
	})
	assert.NoError(t, err, "Invoking the container with the provided type should not error")
}
