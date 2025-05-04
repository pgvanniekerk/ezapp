package container

import (
	"github.com/pgvanniekerk/ezapp/wire"
	"reflect"
)

// isWireComponent checks if the provided type is wire.Component.
//
// Parameters:
//   - t: The reflect.Type to check
//
// Returns:
//   - bool: true if the type is wire.Component, false otherwise
func isWireComponent(t reflect.Type) bool {
	// If t is nil, it can't be wire.Component
	if t == nil {
		return false
	}

	// Get the wire.Component type
	wireComponentType := reflect.TypeOf(wire.Component{})

	// Compare the string representation of the types
	return t.String() == wireComponentType.String()
}
