package link

import (
	"context"
	"testing"

	"github.com/pgvanniekerk/ezapp/pkg/ezapp"
	"github.com/stretchr/testify/assert"
)

// MockStruct is a simple struct for testing
type MockStruct struct {
	Value string
}

// MockBuilder implements the builder interface for MockStruct
type MockBuilder struct {
	Value string
}

// Build returns a MockStruct
func (b MockBuilder) Build(ctx context.Context) (MockStruct, error) {
	return MockStruct{Value: b.Value}, nil
}

func TestBuilderWithStruct(t *testing.T) {
	// Create a MockBuilder
	_ = MockBuilder{Value: "test value"}

	// Create a build process with the Builder function
	buildProc := Builder[MockBuilder]

	// Create a container and apply the build process
	container := ezapp.Construct(buildProc)

	// Verify that the container is not nil
	assert.NotNil(t, container, "The container should not be nil")

	// TODO: Once the Builder function is implemented, add assertions to verify
	// that the MockStruct was provided to the container
}

// MockBuilderPtr implements the builder interface for *MockStruct
type MockBuilderPtr struct {
	Value string
}

// Build returns a pointer to MockStruct
func (b MockBuilderPtr) Build(ctx context.Context) (*MockStruct, error) {
	return &MockStruct{Value: b.Value}, nil
}

func TestBuilderWithPointer(t *testing.T) {
	// Create a MockBuilderPtr
	_ = MockBuilderPtr{Value: "test pointer value"}

	// Create a build process with the Builder function
	buildProc := Builder[MockBuilderPtr]

	// Create a container and apply the build process
	container := ezapp.Construct(buildProc)

	// Verify that the container is not nil
	assert.NotNil(t, container, "The container should not be nil")

	// TODO: Once the Builder function is implemented, add assertions to verify
	// that the *MockStruct was provided to the container
}
