package app

import (
	"context"
)

// Runnable is the internal interface that defines the contract for components
// that can be managed by the App. This interface is implemented by types that
// embed the ezapp.Runnable struct and override its methods.
//
// The Runnable interface defines three methods:
//   - Run: Called to start the component, returns an error if it fails
//   - Stop: Called to stop the component with a context that may include a deadline
//   - Sentinel: A marker method used to identify types that embed the ezapp.Runnable struct
//
// This interface is used internally by the ezapp framework and is not meant to be
// implemented directly by users. Instead, users should embed the ezapp.Runnable
// struct in their component types and override the Run and Stop methods.
type Runnable interface {
	// Run starts the component and returns an error if it fails.
	// This method is called in a separate goroutine by the App.
	// If this method returns an error, the App will initiate shutdown.
	Run() error

	// Stop gracefully stops the component.
	// This method is called with a context that may include a deadline,
	// after which the shutdown process will be aborted.
	// Implementations should respect context cancellation.
	Stop(context.Context) error

	// Sentinel is a marker method used to identify types that embed
	// the ezapp.Runnable struct. This method is used internally by
	// the framework and should not be called directly.
	Sentinel()
}
