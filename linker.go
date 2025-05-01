package ezapp

import "go.uber.org/dig"

// Linker is a function type that registers components with a dependency injection container.
//
// This type represents functions that can register one or more components with a dig.Container.
// The Component function in this package returns a function of this type, which can then be
// used to register components with the application's container.
//
// Linker functions are typically used in the application's startup code to register
// all components with the dependency injection container before running the application.
//
// Example usage:
//
//	// Create a container
//	container := dig.New()
//
//	// Create linkers for different components
//	userServiceLinker := link.Component[*UserService, UserServiceParams](container)
//	authServiceLinker := link.Component[*AuthService, AuthServiceParams](container)
//
//	// Register components with the container
//	if err := userServiceLinker(container); err != nil {
//	    return err
//	}
//	if err := authServiceLinker(container); err != nil {
//	    return err
//	}
type Linker func(*dig.Container) error
