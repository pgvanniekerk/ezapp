package container

import (
	"context"
	"errors"
	"go.uber.org/dig"
	"reflect"
)

// ErrInitMethodNotFound is returned when a component doesn't have an Init method
var ErrInitMethodNotFound = errors.New("component does not have an Init method")

// InitContextIn is a struct that embeds dig.In and has a field for the initialization context.
// It is used by BuildProvideFunc to inject a context with the name "ezapp_initCtx" into components.
type InitContextIn struct {
	dig.In

	// InitCtx is the context used for component initialization.
	// It must be provided to the container with the name "ezapp_initCtx".
	InitCtx context.Context `name:"ezapp_initCtx"`
}

// buildProvideFunc creates a function that can be used with dig.Provide to register a component.
//
// This function is the core of the dependency injection mechanism. It dynamically creates
// a provider function at runtime using reflection. The created function:
//  1. Accepts two structs with dig.In embedded:
//     - First struct contains the component's parameters (created by createDigInStructType)
//     - Second struct is InitContextIn which contains a context field with tag name:"ezapp_initCtx"
//  2. Creates a new instance of the component
//  3. Extracts dependencies from the first dig struct into a Params struct
//  4. Initializes the component with the injected context and Params struct
//  5. Returns the initialized component and an error
//
// The function signature is designed to match what dig.Provide expects:
//
//	func(paramsDigStruct, ctxDigStruct) (component, error)
//
// When the dig container resolves dependencies, it:
//   - Calls this function with dig structs containing the required dependencies
//   - Gets back an initialized component
//   - Stores the component in the container for other components to use
//
// This approach allows components to be created and initialized with their
// dependencies automatically, without manual wiring.
//
// Parameters:
//   - compType: The reflect.Type of the component to create
//   - paramsType: The reflect.Type of the component's parameters struct
//
// Returns:
//   - interface{}: A function that can be passed to dig.Provide
func buildProvideFunc(compType reflect.Type, paramsType reflect.Type) interface{} {
	// Create a function using reflection with the signature:
	// func(paramsDigStruct, ctxDigStruct) (component, error)
	return reflect.MakeFunc(
		// Define the function type (signature)
		reflect.FuncOf(
			// Input parameter types: the params dig struct and the context dig struct
			[]reflect.Type{createDigInStructType(paramsType), reflect.TypeOf(InitContextIn{})},
			// Return value types: the component and an error
			[]reflect.Type{reflect.PointerTo(compType), reflect.TypeOf((*error)(nil)).Elem()},
			// Not a variadic function
			false,
		),
		// Define the function implementation
		func(args []reflect.Value) []reflect.Value {
			// Extract the dig structs from the arguments
			paramsDigStruct := args[0]
			ctxDigStruct := args[1]

			// Extract the initialization context from the context dig struct
			initCtx := ctxDigStruct.Interface().(InitContextIn).InitCtx

			// Convert the dig struct to a Params struct
			params := createParamsFromDigStruct(paramsType, paramsDigStruct)

			// Create a new instance of the component
			compInstance := reflect.New(compType).Interface()

			// Find the Init method on the component
			initMethod := reflect.ValueOf(compInstance).MethodByName("Init")
			if !initMethod.IsValid() {
				// Return an error if the Init method doesn't exist
				return []reflect.Value{
					reflect.Zero(reflect.PointerTo(compType)),
					reflect.ValueOf(ErrInitMethodNotFound),
				}
			}

			// Call the Init method with the context and params
			results := initMethod.Call([]reflect.Value{
				reflect.ValueOf(initCtx),
				params,
			})

			// Check if there was an error during initialization
			errValue := results[0]
			if !errValue.IsNil() {
				// Return zero value for the component and the error
				return []reflect.Value{
					reflect.Zero(reflect.PointerTo(compType)),
					errValue,
				}
			}

			// Return the initialized component and nil error
			return []reflect.Value{
				reflect.ValueOf(compInstance),
				reflect.Zero(reflect.TypeOf((*error)(nil)).Elem()),
			}
		},
	).Interface()
}
