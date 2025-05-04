package container

import (
	"reflect"
)

// createParamsFromDigStruct creates a new instance of a params struct and populates it with values from the dig struct.
//
// This function is a crucial part of the dependency injection process. It:
//  1. Creates a new instance of the params struct based on the provided paramsType
//  2. Copies values from the dig struct (which contains dependencies injected by dig)
//     into the corresponding fields of the params struct
//  3. Returns the populated params struct as a reflect.Value for use in component initialization
//
// The dig struct (created by createDigInStructType) has the same field names as the params struct,
// allowing this function to match fields by name and copy their values. This creates a clean
// separation between the dig-specific input struct and the component's parameter struct.
//
// This approach allows components to define their dependencies using a regular struct
// without any dig-specific annotations, making the component code cleaner and more maintainable.
//
// Parameters:
//   - paramsType: The reflect.Type of the component's parameters struct
//   - digStruct: A reflect.Value containing the dig struct with injected dependencies
//
// Returns:
//   - reflect.Value: A new instance of the params struct with values copied from the dig struct
func createParamsFromDigStruct(paramsType reflect.Type, digStruct reflect.Value) reflect.Value {
	// Create a new instance of the params struct using reflection
	// This creates a zero-initialized struct that we'll populate with values
	paramsInstance := reflect.New(paramsType).Elem()

	// Copy values from the dig struct to the params struct
	// We iterate through each field in the params struct and look for a matching
	// field in the dig struct to copy the value from
	for i := 0; i < paramsType.NumField(); i++ {
		field := paramsType.Field(i)

		// Find the corresponding field in the dig struct by name
		inField := digStruct.FieldByName(field.Name)

		// If the field exists in the dig struct, copy its value to the params struct
		if inField.IsValid() {
			paramsInstance.FieldByName(field.Name).Set(inField)
		}
		// If the field doesn't exist, it remains zero-initialized
	}

	// Return the populated params struct as a reflect.Value
	return paramsInstance
}
