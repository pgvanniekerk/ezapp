package link

import (
	"errors"
	"github.com/pgvanniekerk/ezapp/internal/primitive"
	"go.uber.org/dig"
)

// Component registers a component with the dependency injection container.
//
// This function is a key part of the application's dependency injection system.
// It takes a component type and its parameters type as generic type parameters,
// and registers the component with the provided dig container.
//
// The function performs the following steps:
// 1. Validates that the Params type is a struct (required for dependency injection)
// 2. Gets the core type of the component (handling pointer types)
// 3. Builds a provider function that:
//   - Creates a new instance of the component
//   - Injects dependencies from the container into the component's parameters
//   - Uses a context with the name "ezapp_initCtx" for component initialization
//   - Initializes the component with the injected context and parameters
//
// 4. Registers this provider function with the dig container
//
// Note: This function requires a context with the name "ezapp_initCtx" to be provided
// to the container. This context will be used for component initialization.
//
// Usage example:
//
//	type MyComponent struct{}
//	type MyParams struct {
//	    Service1 *Service1
//	    Service2 *Service2
//	}
//
//	func (c *MyComponent) Init(ctx context.Context, params MyParams) error {
//	    // Initialize component with injected dependencies
//	    return nil
//	}
//
//	func (c *MyComponent) Cleanup(ctx context.Context) error {
//	    // Cleanup resources
//	    return nil
//	}
//
//	// Register the component with the container
//	err := link.Component[*MyComponent, MyParams](container)
//
// Parameters:
//   - Comp: A type that implements the primitive.Component[Params] interface
//   - Params: A struct type containing the component's dependencies
//   - digC: The dependency injection container
//
// Returns:
//   - error: An error if registration fails, nil otherwise
func Component[Comp primitive.Component[Params], Params any](digC *dig.Container) error {
	// Validate that Params is a struct - this is required for dependency injection
	// as the fields of the struct represent the dependencies to be injected
	if !validateParamsType[Params]() {
		return errors.New("type parameter 'Params' must be a struct")
	}

	// Provide the function to the dig container
	// This registers a provider function that will create and initialize the component
	// when its dependencies are requested
	return digC.Provide(
		BuildProvideFunc[Params](
			getCoreType[Comp, Params](),
		),
	)
}
