package link

import (
	"github.com/pgvanniekerk/ezapp/internal/primitive"
	"reflect"
)

// getCoreType returns the core (non-pointer) type of a component.
//
// This function is an internal helper used by the Component function to handle
// both pointer and non-pointer component types uniformly. When creating a new
// instance of a component using reflection, we need the underlying type (not the pointer type).
//
// For example:
// - If Comp is *MyComponent, this returns the type of MyComponent
// - If Comp is MyComponent, this returns the type of MyComponent
//
// This allows the BuildProvideFunc to create a new instance of the component
// using reflect.New(compType) regardless of whether the component type was
// originally specified as a pointer or value type.
//
// Parameters:
//   - Comp: The component type (implements primitive.Component[Params])
//   - Params: The parameters type for the component
//
// Returns:
//   - reflect.Type: The core (non-pointer) type of the component
func getCoreType[Comp primitive.Component[Params], Params any]() reflect.Type {
	// Get the type of Comp using reflection
	var comp Comp
	compType := reflect.TypeOf(comp)

	// If the component type is a pointer, return its element type
	// This handles cases where Comp is *MyComponent
	if compType.Kind() == reflect.Ptr {
		return compType.Elem()
	}

	// Otherwise, return the type as is
	// This handles cases where Comp is MyComponent
	return compType
}
