package main

import (
	"context"
	"fmt"
	"time"
)

// MockService is a mock service that implements the Service interface
type MockService struct {
	name string
	db   *MockDB
}

// NewMockService creates a new MockService
func NewMockService(name string, db *MockDB) *MockService {
	return &MockService{
		name: name,
		db:   db,
	}
}

// Run starts the mock service
// It should only return an error in exceptional circumstances such as dependency failures or timeouts (application-impacting errors)
func (s *MockService) Run() error {
	fmt.Printf("Starting mock service: %s\n", s.name)

	// Simulate some work
	go func() {
		for {
			// Perform a database query to demonstrate using the DB
			now, err := s.db.Query(context.Background(), "SELECT NOW()")
			if err != nil {
				fmt.Printf("Error querying database: %v\n", err)
			} else {
				fmt.Printf("Current database time: %v\n", now)
			}

			// Sleep for a while
			time.Sleep(5 * time.Second)
		}
	}()

	return nil
}

// Stop gracefully shuts down the mock service
// If it returns an error, it will be reported during shutdown
// If the context timeout is reached, the application will force close
func (s *MockService) Stop(ctx context.Context) error {
	fmt.Printf("Stopping mock service: %s\n", s.name)

	// Simulate cleanup work
	select {
	case <-time.After(1 * time.Second):
		fmt.Println("Mock service cleanup completed")
		return nil
	case <-ctx.Done():
		fmt.Println("Mock service cleanup timed out")
		return ctx.Err()
	}
}
