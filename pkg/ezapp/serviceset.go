package ezapp

import (
	"github.com/pgvanniekerk/ezapp/internal/app"
)

// NewServiceSet creates a new ServiceSet with the given options
func NewServiceSet(opts ...ServiceSetOption) ServiceSet {
	s := &ServiceSet{
		Services: make([]app.Service, 0),
	}

	for _, opt := range opts {
		opt(s)
	}

	return *s
}

// ServiceSet represents a set of services and an optional cleanup function
type ServiceSet struct {
	Services    []app.Service
	CleanupFunc CleanupFunc
}

// ServiceSetOption is a function that configures a ServiceSet
type ServiceSetOption func(*ServiceSet)

// WithCleanupFunc sets the cleanup function for a ServiceSet
func WithCleanupFunc(cleanup CleanupFunc) ServiceSetOption {
	return func(s *ServiceSet) {
		s.CleanupFunc = cleanup
	}
}

// WithServices adds multiple services to the ServiceSet
func WithServices(services ...app.Service) ServiceSetOption {
	return func(s *ServiceSet) {
		s.Services = append(s.Services, services...)
	}
}

// CleanupFunc is a function that performs cleanup operations and returns an error if any occur
type CleanupFunc func() error
