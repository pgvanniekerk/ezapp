package ezapp

import (
	"context"
	"go.uber.org/dig"
)

// Construct creates and returns a new Dig container with the provided options applied
func Construct(bProcs ...BuildProcess) *dig.Container {

	container := dig.New()

	// Create a build context with the container
	bCtx := &buildContext{
		container: container,
		initCtx:   context.Background(),
	}

	// Apply all build processes
	for _, bProc := range bProcs {
		_ = bProc(bCtx)
	}

	return container
}
