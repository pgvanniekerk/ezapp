package main

import (
	"context"
	"fmt"
	"time"
)

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
