package container

import (
	"context"
	"reflect"
)

// implementsComponent checks if the provided type implements the primitive.Component interface.
//
// Parameters:
//   - t: The reflect.Type to check
//
// Returns:
//   - bool: true if the type implements primitive.Component, false otherwise
func implementsComponent(t reflect.Type) bool {
	// If t is nil, it can't implement the interface
	if t == nil {
		return false
	}

	// Create a pointer to the type if it's not already a pointer
	ptrType := t
	if t.Kind() != reflect.Ptr {
		ptrType = reflect.PointerTo(t)
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
	return validateInit(initMethod.Type) && validateCleanup(cleanupMethod.Type)
}

// validateInit checks if the method signature matches the Init method of primitive.Component.
//
// Parameters:
//   - methodType: The reflect.Type of the method to check
//
// Returns:
//   - bool: true if the method signature matches, false otherwise
func validateInit(methodType reflect.Type) bool {
	// Init method should have 3 parameters:
	// 1. Receiver (the component itself)
	// 2. context.Context
	// 3. Params (which should be a struct)
	if methodType.NumIn() != 3 {
		return false
	}

	// Check that the second parameter is context.Context
	if !methodType.In(1).Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
		return false
	}

	// Check that the third parameter is a struct
	if methodType.In(2).Kind() != reflect.Struct {
		return false
	}

	// Init method should return 1 value: error
	if methodType.NumOut() != 1 {
		return false
	}

	// Check that the return value is error
	return methodType.Out(0).Implements(reflect.TypeOf((*error)(nil)).Elem())
}

// validateCleanup checks if the method signature matches the Cleanup method of primitive.Component.
//
// Parameters:
//   - methodType: The reflect.Type of the method to check
//
// Returns:
//   - bool: true if the method signature matches, false otherwise
func validateCleanup(methodType reflect.Type) bool {
	// Cleanup method should have 2 parameters:
	// 1. Receiver (the component itself)
	// 2. context.Context
	if methodType.NumIn() != 2 {
		return false
	}

	// Check that the second parameter is context.Context
	if !methodType.In(1).Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
		return false
	}

	// Cleanup method should return 1 value: error
	if methodType.NumOut() != 1 {
		return false
	}

	// Check that the return value is error
	return methodType.Out(0).Implements(reflect.TypeOf((*error)(nil)).Elem())
}
