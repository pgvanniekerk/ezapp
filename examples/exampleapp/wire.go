package main

import (
	"context"
	"fmt"
	"github.com/pgvanniekerk/ezapp/pkg/ezapp"
	"time"
)

// Config holds the application configuration
type Config struct {
	ServiceName string `envconfig:"SERVICE_NAME" default:"exampleapp"`
}

// MockDB is a simple mock database for demonstration purposes
type MockDB struct {
	isConnected bool
}

// NewMockDB creates a new mock database
func NewMockDB() *MockDB {
	return &MockDB{
		isConnected: false,
	}
}

// Connect simulates connecting to a database
func (db *MockDB) Connect(ctx context.Context) error {
	select {
	case <-time.After(500 * time.Millisecond):
		db.isConnected = true
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Close simulates closing a database connection
func (db *MockDB) Close() error {
	if !db.isConnected {
		return fmt.Errorf("database not connected")
	}
	db.isConnected = false
	return nil
}

// Query simulates a database query
func (db *MockDB) Query(ctx context.Context, _ string) (time.Time, error) {
	if !db.isConnected {
		return time.Time{}, fmt.Errorf("database not connected")
	}

	select {
	case <-time.After(100 * time.Millisecond):
		return time.Now(), nil
	case <-ctx.Done():
		return time.Time{}, ctx.Err()
	}
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
