package main

import (
	"context"
	"fmt"
	"github.com/pgvanniekerk/ezapp/pkg/ezapp"
)

// Config holds the application configuration
type Config struct {
	ServiceName string `envconfig:"SERVICE_NAME" default:"exampleapp"`
}

// WireFunc is the function used for dependency injection and wiring
// It connects services and dependencies together
func WireFunc(startupCtx context.Context, config Config) (ezapp.ServiceSet, error) {

	// Create a mock database
	fmt.Println("Creating mock database")
	db := NewMockDB()

	// Connect to the database using the startup context
	fmt.Println("Connecting to mock database")
	if err := db.Connect(startupCtx); err != nil {
		return ezapp.ServiceSet{}, fmt.Errorf("failed to connect to mock database: %w", err)
	}

	// Create a mock service
	mockService := NewMockService(config.ServiceName, db)

	// Create a cleanup function for the database
	cleanupFunc := func() error {
		fmt.Println("Closing mock database connection")
		return db.Close()
	}

	// Return a ServiceSet with the mock service and cleanup function
	return ezapp.NewServiceSet(
		ezapp.WithServices(mockService),
		ezapp.WithCleanupFunc(cleanupFunc),
	), nil
}
