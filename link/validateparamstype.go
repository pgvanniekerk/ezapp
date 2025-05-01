package link

import (
	"reflect"
)

// validateParamsType checks if the provided type parameter is a struct.
//
// This function is an internal helper used by the Component function to ensure
// that the Params type is a struct, which is a requirement for dependency injection.
// The dig container expects structs for dependency injection because:
//  1. Struct fields represent the dependencies to be injected
//  2. Field names and tags are used to match dependencies in the container
//  3. Structs provide a clear, type-safe way to group related dependencies
//
// If the Params type is not a struct (e.g., it's a primitive type, pointer, or interface),
// the Component function will return an error, preventing invalid usage.
//
// Parameters:
//   - Params: The type to validate
//
// Returns:
//   - bool: true if the type is a struct, false otherwise
func validateParamsType[Params any]() bool {
	// Create a zero value of the Params type
	var params Params

	// Get the reflect.Type of the Params type
	paramsType := reflect.TypeOf(params)

	// Check if the kind of the type is a struct
	return paramsType.Kind() == reflect.Struct
}
