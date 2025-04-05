package app

import (
	"log/slog"
	"reflect"
	"strings"

	"github.com/pgvanniekerk/ezapp/pkg/ezapp"
)

// setRunnableLogger examines a Runnable object using reflection to find an anonymous field
// for ezapp.Runnable with the tag toggle:"useEzAppLogger". If found, it sets the Logger field
// of the ezapp.Runnable struct to the provided logger with additional type and package information.
func setRunnableLogger(runnable Runnable, logger *slog.Logger) {
	// Get the value of the runnable
	val := reflect.ValueOf(runnable)

	// Get the type name and package path of the runnable
	runnableType := reflect.TypeOf(runnable)
	// If the runnable is a pointer, get the element type
	if runnableType.Kind() == reflect.Ptr {
		runnableType = runnableType.Elem()
	}
	typeName := runnableType.Name()
	packagePath := runnableType.PkgPath()

	// Extract just the package name from the path
	packageName := packagePath
	if lastSlash := strings.LastIndex(packagePath, "/"); lastSlash >= 0 {
		packageName = packagePath[lastSlash+1:]
	}

	// Add type and package information to the logger
	loggerWithAttrs := logger.With(
		slog.String("typeName", typeName),
		slog.String("packagePath", packageName),
	)

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
				// Check if the field has the tag toggle:"useEzAppLogger"
				if field.Tag.Get("toggle") == "useEzAppLogger" {
					// Get the field value
					fieldVal := val.Field(i)

					// Find the Logger field in the ezapp.Runnable struct
					loggerField := fieldVal.FieldByName("Logger")

					// If the Logger field exists and is settable, set it to the logger with attributes
					if loggerField.IsValid() && loggerField.CanSet() {
						loggerField.Set(reflect.ValueOf(loggerWithAttrs))
					}

					return
				}
			}
		}
	}
}
