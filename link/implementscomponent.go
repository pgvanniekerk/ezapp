package link

import (
	"github.com/pgvanniekerk/ezapp/internal/container"
	"reflect"
)

// implementsComponent checks if the provided type is or implements primitive.Component.
//
// Parameters:
//   - fieldType: The reflect.Type to check
//
// Returns:
//   - bool: true if the type is or implements primitive.Component, false otherwise
func implementsComponent(fieldType reflect.Type) bool {
	// Get the string representation of the field type
	fieldTypeStr := fieldType.String()

	// Check if the field type is primitive.Component or a generic instantiation of it
	if fieldTypeStr == "primitive.Component" ||
		fieldTypeStr == "github.com/pgvanniekerk/ezapp/internal/primitive.Component" ||
		(len(fieldTypeStr) > 58 && fieldTypeStr[:58] == "github.com/pgvanniekerk/ezapp/internal/primitive.Component[") {
		return true
	}

	// If the field is an interface type, we need to check if it's primitive.Component
	if fieldType.Kind() == reflect.Interface {
		// For interfaces, we can check the package path and name
		// The name might include generic type parameters, so we check if it starts with "Component"
		if fieldType.PkgPath() == "github.com/pgvanniekerk/ezapp/internal/primitive" &&
			(fieldType.Name() == "Component" || len(fieldType.Name()) > 9 && fieldType.Name()[:9] == "Component") {
			return true
		}
	}

	// If we're dealing with a struct type, check if it implements the interface
	// by having the required methods
	if fieldType.Kind() == reflect.Struct ||
		(fieldType.Kind() == reflect.Ptr && fieldType.Elem().Kind() == reflect.Struct) {
		// Create a pointer to the type if it's not already a pointer
		ptrType := fieldType
		if fieldType.Kind() != reflect.Ptr {
			ptrType = reflect.PtrTo(fieldType)
		}

		// Check if the type has Init and Cleanup methods
		initMethod, hasInit := ptrType.MethodByName("Init")
		if !hasInit {
			return false
		}

		cleanupMethod, hasCleanup := ptrType.MethodByName("Cleanup")
		if !hasCleanup {
			return false
		}

		// Validate Init and Cleanup method signatures
		return container.validateInit(initMethod.Type) && container.validateCleanup(cleanupMethod.Type)
	}

	return false
}
