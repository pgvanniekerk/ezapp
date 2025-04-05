package app

import (
	"fmt"
	"reflect"

	"github.com/pgvanniekerk/ezapp/pkg/ezapp"
)

// EnsureEmbedsRunnableStruct checks if the provided Runnable interface instance
// embeds the ezapp.Runnable struct. If not, it returns an error.
func EnsureEmbedsRunnableStruct(runnable Runnable) error {
	// Get the value of the runnable
	val := reflect.ValueOf(runnable)

	// If the runnable is a pointer, get the value it points to
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// If it's not a struct, we can't proceed
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("runnable is not a struct or a pointer to a struct")
	}

	// Get the type of the struct
	typ := val.Type()

	// Iterate through all fields of the struct
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// Check if this is an anonymous field
		if field.Anonymous {
			// Check if the field is of type ezapp.Runnable
			if field.Type == reflect.TypeOf(ezapp.Runnable{}) {
				// Found the embedded ezapp.Runnable struct
				return nil
			}
		}
	}

	// If we get here, the ezapp.Runnable struct is not embedded
	return fmt.Errorf("runnable does not embed ezapp.Runnable struct")
}
