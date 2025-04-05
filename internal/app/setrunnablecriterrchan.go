package app

import (
	"reflect"

	"github.com/pgvanniekerk/ezapp/pkg/ezapp"
)

// setRunnableCritErrChan examines a Runnable object using reflection to find an anonymous field
// for ezapp.Runnable. If found, it sets the critErrChan field of the ezapp.Runnable struct
// to the provided channel.
func setRunnableCritErrChan(runnable Runnable, critErrChan chan<- error) {
	// Get the value of the runnable
	val := reflect.ValueOf(runnable)

	// If the runnable is a pointer, get the value it points to
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// If it's not a struct, we can't proceed
	if val.Kind() != reflect.Struct {
		return
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
				// Get the field value
				fieldVal := val.Field(i)

				// Find the critErrChan field in the ezapp.Runnable struct
				critErrChanField := fieldVal.FieldByName("critErrChan")

				// If the critErrChan field exists and is settable, set it to the provided channel
				if critErrChanField.IsValid() && critErrChanField.CanSet() {
					critErrChanField.Set(reflect.ValueOf(critErrChan))
				}

				return
			}
		}
	}
}
