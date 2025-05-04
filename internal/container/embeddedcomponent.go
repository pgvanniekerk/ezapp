package container

import (
	"reflect"
)

// embeddedComponent finds the first anonymous field in a struct type that implements the primitive.Component interface.
//
// Parameters:
//   - compType: The reflect.Type of the struct to check
//
// Returns:
//   - reflect.Type: The type of the first anonymous field that implements primitive.Component, or nil if none found
//   - bool: true if an anonymous field implementing primitive.Component was found, false otherwise
func embeddedComponent(compType reflect.Type) (reflect.Type, bool) {
	// If compType is nil or not a struct, return false
	if compType == nil || compType.Kind() != reflect.Struct {
		return nil, false
	}

	// Iterate through all fields of the struct
	for i := 0; i < compType.NumField(); i++ {
		field := compType.Field(i)

		// Check if the field is anonymous (embedded)
		if field.Anonymous {
			fieldType := field.Type

			// Check if the field type implements primitive.Component
			if implementsComponent(fieldType) {
				return fieldType, true
			}
		}
	}

	// No anonymous field implementing primitive.Component was found
	return nil, false
}
