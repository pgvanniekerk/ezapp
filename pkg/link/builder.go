package link

import (
	"context"
	"fmt"
	"github.com/pgvanniekerk/ezapp/pkg/ezapp"
	"reflect"
)

// provideBuildFunction provides the Build function to the container
func provideBuildFunction[B builder[T], T any](bCtx ezapp.BuildContext) error {
	err := bCtx.Container().Provide(
		func(b B) (T, error) {
			t, err := b.Build(bCtx.InitTimeout())
			if err != nil {
				return t, fmt.Errorf("error building %T: %w", t, err)
			}
			return t, nil
		},
	)
	return err
}

// retrieveStructType validates that B is a struct or pointer to a struct and returns its reflect.Type
func retrieveStructType[B builder[T], T any]() (reflect.Type, error) {
	var b B
	bType := reflect.TypeOf(b)

	// Handle pointer type
	if bType.Kind() == reflect.Ptr {
		bType = bType.Elem()
	}

	// Validate that B is a struct
	if bType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct type, got %v", bType.Kind())
	}

	return bType, nil
}

// collateModuleDependencies collects module dependencies for a given type
func collateModuleDependencies(bType reflect.Type, bCtx ezapp.BuildContext) ([]reflect.Type, error) {
	modDeps := make([]reflect.Type, 0, len(bCtx.Modules()))
	for _, mod := range bCtx.Modules() {
		modDep, err := mod.GetDependencies(bType)
		if err != nil {
			return nil, fmt.Errorf("error getting dependencies for %v: %w", bType, err)
		}
		if modDep != nil {
			modDeps = append(modDeps, modDep)
		}
	}

	return modDeps, nil
}

func Builder[B builder[T], T any](bCtx ezapp.BuildContext) error {

	// First get the struct type and validate it
	bType, err := retrieveStructType[B, T]()
	if err != nil {
		return err
	}

	// Provide the Build function
	if err := provideBuildFunction[B, T](bCtx); err != nil {
		return err
	}

	// Collate module dependencies
	_, err = collateModuleDependencies(bType, bCtx)
	if err != nil {
		return err
	}

	return nil
}

type builder[T any] interface {
	Build(context.Context) (T, error)
}
