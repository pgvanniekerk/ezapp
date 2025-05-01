package link

import (
	"go.uber.org/dig"
	"reflect"
)

// CreateDigInStructType creates a new reflect struct type that embeds dig.In and has the same fields as Params.
//
// This function is a critical part of the dependency injection mechanism. It dynamically
// creates a new struct type at runtime that:
//  1. Embeds dig.In (required by dig to mark a struct as an input parameter)
//  2. Contains all the fields from the Params struct
//
// The resulting struct type is used by BuildProvideFunc to create a function signature
// that dig can use to inject dependencies. When dig calls this function, it will:
//   - Create an instance of this struct
//   - Fill its fields with dependencies from the container
//   - Pass this struct to the provider function
//
// This approach allows the Component function to work with any Params struct type
// without requiring manual creation of dig.In structs for each component.
//
// Parameters:
//   - Params: The type parameter representing the component's parameters
//
// Returns:
//   - reflect.Type: A new struct type that embeds dig.In and has the same fields as Params
func CreateDigInStructType[Params any]() reflect.Type {
	// Get the reflect.Type of the Params type
	var params Params
	paramsType := reflect.TypeOf(params)

	// Create a slice to hold the fields for the new struct type
	// We need capacity for all fields in Params plus one for the embedded dig.In
	fields := make([]reflect.StructField, 0, paramsType.NumField()+1)

	// Add dig.In as an embedded field
	// This marks the struct as an input parameter for dig
	fields = append(fields, reflect.StructField{
		Name:      "In",
		Type:      reflect.TypeOf(dig.In{}),
		Anonymous: true, // Anonymous field makes it an embedded field
	})

	// Add all fields from the Params struct to the new struct type
	// These fields will be filled with dependencies from the container
	for i := 0; i < paramsType.NumField(); i++ {
		field := paramsType.Field(i)
		fields = append(fields, field)
	}

	// Create and return the new struct type using reflection
	return reflect.StructOf(fields)
}
