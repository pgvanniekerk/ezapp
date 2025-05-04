package container

import (
	"fmt"
	"go.uber.org/dig"
	"reflect"
)

type Container struct {
	digC *dig.Container
}

// NewContainer creates a new Container with a new dig.Container
func NewContainer() *Container {
	return &Container{
		digC: dig.New(),
	}
}

func (c *Container) Run() {

}

// getParamsType extracts the parameter type from a component type
func getParamsType(compType reflect.Type) (reflect.Type, error) {
	// Create a pointer to the type if it's not already a pointer
	ptrType := compType
	if compType.Kind() != reflect.Ptr {
		ptrType = reflect.PointerTo(compType)
	}

	// Find the Init method
	initMethod, hasInit := ptrType.MethodByName("Init")
	if !hasInit {
		return nil, fmt.Errorf("type %q does not have an Init method", compType)
	}

	// The Init method should have 3 parameters:
	// 1. Receiver (the component itself)
	// 2. context.Context
	// 3. Params (which should be a struct)
	if initMethod.Type.NumIn() != 3 {
		return nil, fmt.Errorf("Init method of type %q has wrong number of parameters", compType)
	}

	// Return the third parameter type (Params)
	return initMethod.Type.In(2), nil
}

func (c *Container) LinkComponent(compType reflect.Type) error {
	// Find the first anonymous field that implements primitive.Component
	embeddedCompType, found := embeddedComponent(compType)
	if !found {
		return fmt.Errorf("type %q does not embed a field that implements primitive.Component", compType)
	}

	// Check if the embedded field is a wire.Component.
	// If not, we need to call LinkComponent for the embedded field type.
	if !isWireComponent(embeddedCompType) {
		err := c.LinkComponent(embeddedCompType)
		if err != nil {
			return err
		}
	}

	// Get the parameter type for the component
	paramsType, err := getParamsType(compType)
	if err != nil {
		return err
	}

	// Create a provider function for the component and register it with the dig container
	err = c.digC.Provide(buildProvideFunc(compType, paramsType))
	if err != nil {
		return fmt.Errorf("failed to provide component %q: %w", compType, err)
	}

	return nil
}
