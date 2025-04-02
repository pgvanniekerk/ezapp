package app

import "context"

// Service defines the interface for components that can be run by the application.
type Service interface {

	// Run starts the service and should only return an error in exceptional circumstances
	// such as dependency failures or timeouts (application-impacting errors).
	Run() error

	// Stop gracefully shuts down the service. If it returns an error, it will be reported
	// during shutdown. If the context timeout is reached, the application will force close.
	Stop(context.Context) error
}
