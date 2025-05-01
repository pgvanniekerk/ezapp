package link

import (
	"reflect"
)

// CreateParamsFromDigStruct creates a new instance of Params and populates it with values from the dig struct.
//
// This function is a crucial part of the dependency injection process. It:
//  1. Creates a new instance of the Params struct
//  2. Copies values from the dig struct (which contains dependencies injected by dig)
//     into the corresponding fields of the Params struct
//  3. Returns the populated Params struct for use in component initialization
//
// The dig struct (created by CreateDigInStructType) has the same field names as the Params struct,
// allowing this function to match fields by name and copy their values. This creates a clean
// separation between the dig-specific input struct and the component's parameter struct.
//
// This approach allows components to define their dependencies using a regular struct
// without any dig-specific annotations, making the component code cleaner and more maintainable.
//
// Parameters:
//   - Params: The type parameter representing the component's parameters
//   - digStruct: A reflect.Value containing the dig struct with injected dependencies
//
// Returns:
//   - Params: A new instance of the Params struct with values copied from the dig struct
func CreateParamsFromDigStruct[Params any](digStruct reflect.Value) Params {
	// Get the reflect.Type of the Params type
	var params Params
	paramsType := reflect.TypeOf(params)

	// Create a new instance of the Params struct using reflection
	// This creates a zero-initialized struct that we'll populate with values
	paramsInstance := reflect.New(paramsType).Elem()

	// Copy values from the dig struct to the Params struct
	// We iterate through each field in the Params struct and look for a matching
	// field in the dig struct to copy the value from
	for i := 0; i < paramsType.NumField(); i++ {
		field := paramsType.Field(i)

		// Find the corresponding field in the dig struct by name
		inField := digStruct.FieldByName(field.Name)

		// If the field exists in the dig struct, copy its value to the Params struct
		if inField.IsValid() {
			paramsInstance.FieldByName(field.Name).Set(inField)
		}
		// If the field doesn't exist, it remains zero-initialized
	}

	// Convert the reflect.Value back to the concrete Params type and return it
	return paramsInstance.Interface().(Params)
}
