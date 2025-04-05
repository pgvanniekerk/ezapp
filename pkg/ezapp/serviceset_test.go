package ezapp

import (
	"context"
	"errors"
	"testing"

	"github.com/pgvanniekerk/ezapp/internal/app"
)

// ==================== Mock Implementations ====================

// mockService is a mock implementation of the app.Service interface for testing
type mockService struct{}

func (m *mockService) Run() error {
	return nil
}

func (m *mockService) Stop(ctx context.Context) error {
	return nil
}

// mockServiceForNewServiceSet is a mock implementation of the app.Service interface for testing
type mockServiceForNewServiceSet struct{}

func (m *mockServiceForNewServiceSet) Run() error {
	return nil
}

func (m *mockServiceForNewServiceSet) Stop(ctx context.Context) error {
	return nil
}

// mockServiceForWithServices is a mock implementation of the app.Service interface for testing
type mockServiceForWithServices struct{}

func (m *mockServiceForWithServices) Run() error {
	return nil
}

func (m *mockServiceForWithServices) Stop(ctx context.Context) error {
	return nil
}

// ==================== ServiceSet Tests ====================

// TestServiceSetStructure tests the structure of the ServiceSet struct
func TestServiceSetStructure(t *testing.T) {
	// Create a new ServiceSet
	s := ServiceSet{
		Services:    make([]app.Service, 0),
		CleanupFunc: func() error { return nil },
	}

	// Check that the Services field is initialized correctly
	if s.Services == nil {
		t.Error("ServiceSet.Services should not be nil")
	}

	// Check that the CleanupFunc field is set correctly
	if s.CleanupFunc == nil {
		t.Error("ServiceSet.CleanupFunc should not be nil")
	}

	// Test the CleanupFunc
	err := s.CleanupFunc()
	if err != nil {
		t.Errorf("CleanupFunc returned an error: %v", err)
	}
}

// TestServiceSetWithError tests the ServiceSet with a CleanupFunc that returns an error
func TestServiceSetWithError(t *testing.T) {
	// Create a new ServiceSet with a CleanupFunc that returns an error
	expectedErr := errors.New("cleanup error")
	s := ServiceSet{
		Services:    make([]app.Service, 0),
		CleanupFunc: func() error { return expectedErr },
	}

	// Test the CleanupFunc
	err := s.CleanupFunc()
	if err != expectedErr {
		t.Errorf("CleanupFunc returned %v, want %v", err, expectedErr)
	}
}

// ==================== ServiceSetOption Tests ====================

// TestServiceSetOption tests the ServiceSetOption type
func TestServiceSetOption(t *testing.T) {
	// Create a new ServiceSet
	s := &ServiceSet{
		Services: make([]app.Service, 0),
	}

	// Create a test ServiceSetOption
	testOption := func(s *ServiceSet) {
		s.CleanupFunc = func() error { return nil }
	}

	// Apply the option
	testOption(s)

	// Check that the option was applied correctly
	if s.CleanupFunc == nil {
		t.Error("ServiceSetOption did not set CleanupFunc")
	}

	// Test the CleanupFunc
	err := s.CleanupFunc()
	if err != nil {
		t.Errorf("CleanupFunc returned an error: %v", err)
	}
}

// TestMultipleServiceSetOptions tests applying multiple ServiceSetOptions
func TestMultipleServiceSetOptions(t *testing.T) {
	// Create a new ServiceSet
	s := &ServiceSet{
		Services: make([]app.Service, 0),
	}

	// Create test ServiceSetOptions
	option1 := func(s *ServiceSet) {
		s.CleanupFunc = func() error { return nil }
	}

	// Create a mock Service implementation for testing
	mockSvc := &mockService{}

	option2 := func(s *ServiceSet) {
		s.Services = append(s.Services, mockSvc)
	}

	// Apply the options
	option1(s)
	option2(s)

	// Check that the options were applied correctly
	if s.CleanupFunc == nil {
		t.Error("option1 did not set CleanupFunc")
	}

	if len(s.Services) != 1 {
		t.Errorf("option2 did not add service, got %d services", len(s.Services))
	}
}

// ==================== NewServiceSet Tests ====================

// TestNewServiceSet tests the NewServiceSet function
func TestNewServiceSet(t *testing.T) {
	// Create a new ServiceSet with no options
	s := NewServiceSet()

	// Check that the Services field is initialized correctly
	if s.Services == nil {
		t.Error("ServiceSet.Services should not be nil")
	}

	// Check that the CleanupFunc field is nil
	if s.CleanupFunc != nil {
		t.Error("ServiceSet.CleanupFunc should be nil when no options are provided")
	}
}

// TestNewServiceSetWithOptions tests the NewServiceSet function with options
func TestNewServiceSetWithOptions(t *testing.T) {
	// Create a mock Service implementation for testing
	mockSvc := &mockServiceForNewServiceSet{}

	// Create a cleanup function for testing
	cleanupCalled := false
	cleanup := func() error {
		cleanupCalled = true
		return nil
	}

	// Create a new ServiceSet with options
	s := NewServiceSet(
		WithServices(mockSvc),
		WithCleanupFunc(cleanup),
	)

	// Check that the Services field contains the mock service
	if len(s.Services) != 1 {
		t.Errorf("ServiceSet.Services should contain 1 service, got %d", len(s.Services))
	}

	// Check that the CleanupFunc field is set correctly
	if s.CleanupFunc == nil {
		t.Error("ServiceSet.CleanupFunc should not be nil")
	}

	// Test the CleanupFunc
	err := s.CleanupFunc()
	if err != nil {
		t.Errorf("CleanupFunc returned an error: %v", err)
	}

	// Check that the cleanup function was called
	if !cleanupCalled {
		t.Error("CleanupFunc was not called")
	}
}

// ==================== WithCleanupFunc Tests ====================

// TestWithCleanupFunc tests the WithCleanupFunc function
func TestWithCleanupFunc(t *testing.T) {
	// Create a new ServiceSet
	s := &ServiceSet{}

	// Create a cleanup function for testing
	cleanupCalled := false
	cleanup := func() error {
		cleanupCalled = true
		return nil
	}

	// Create a ServiceSetOption using WithCleanupFunc
	option := WithCleanupFunc(cleanup)

	// Apply the option to the ServiceSet
	option(s)

	// Check that the CleanupFunc field is set correctly
	if s.CleanupFunc == nil {
		t.Error("WithCleanupFunc did not set CleanupFunc")
	}

	// Test the CleanupFunc
	err := s.CleanupFunc()
	if err != nil {
		t.Errorf("CleanupFunc returned an error: %v", err)
	}

	// Check that the cleanup function was called
	if !cleanupCalled {
		t.Error("CleanupFunc was not called")
	}
}

// TestWithCleanupFuncError tests the WithCleanupFunc function with a function that returns an error
func TestWithCleanupFuncError(t *testing.T) {
	// Create a new ServiceSet
	s := &ServiceSet{}

	// Create a cleanup function that returns an error
	expectedErr := errors.New("cleanup error")
	cleanup := func() error {
		return expectedErr
	}

	// Create a ServiceSetOption using WithCleanupFunc
	option := WithCleanupFunc(cleanup)

	// Apply the option to the ServiceSet
	option(s)

	// Check that the CleanupFunc field is set correctly
	if s.CleanupFunc == nil {
		t.Error("WithCleanupFunc did not set CleanupFunc")
	}

	// Test the CleanupFunc
	err := s.CleanupFunc()
	if err != expectedErr {
		t.Errorf("CleanupFunc returned %v, want %v", err, expectedErr)
	}
}

// TestWithCleanupFuncNil tests the WithCleanupFunc function with a nil function
func TestWithCleanupFuncNil(t *testing.T) {
	// Create a new ServiceSet
	s := &ServiceSet{}

	// Create a ServiceSetOption using WithCleanupFunc with nil
	option := WithCleanupFunc(nil)

	// Apply the option to the ServiceSet
	option(s)

	// Check that the CleanupFunc field is set to nil
	if s.CleanupFunc != nil {
		t.Error("WithCleanupFunc did not set CleanupFunc to nil")
	}
}

// ==================== WithServices Tests ====================

// TestWithServices tests the WithServices function
func TestWithServices(t *testing.T) {
	// Create a new ServiceSet
	s := &ServiceSet{
		Services: make([]app.Service, 0),
	}

	// Create mock services for testing
	mockSvc1 := &mockServiceForWithServices{}
	mockSvc2 := &mockServiceForWithServices{}

	// Create a ServiceSetOption using WithServices
	option := WithServices(mockSvc1, mockSvc2)

	// Apply the option to the ServiceSet
	option(s)

	// Check that the Services field contains the mock services
	if len(s.Services) != 2 {
		t.Errorf("WithServices did not add services correctly, got %d services", len(s.Services))
	}
}

// TestWithServicesEmpty tests the WithServices function with no services
func TestWithServicesEmpty(t *testing.T) {
	// Create a new ServiceSet
	s := &ServiceSet{
		Services: make([]app.Service, 0),
	}

	// Create a ServiceSetOption using WithServices with no services
	option := WithServices()

	// Apply the option to the ServiceSet
	option(s)

	// Check that the Services field is still empty
	if len(s.Services) != 0 {
		t.Errorf("WithServices added services when none were provided, got %d services", len(s.Services))
	}
}

// TestWithServicesMultiple tests the WithServices function with multiple calls
func TestWithServicesMultiple(t *testing.T) {
	// Create a new ServiceSet
	s := &ServiceSet{
		Services: make([]app.Service, 0),
	}

	// Create mock services for testing
	mockSvc1 := &mockServiceForWithServices{}
	mockSvc2 := &mockServiceForWithServices{}
	mockSvc3 := &mockServiceForWithServices{}

	// Create ServiceSetOptions using WithServices
	option1 := WithServices(mockSvc1)
	option2 := WithServices(mockSvc2, mockSvc3)

	// Apply the options to the ServiceSet
	option1(s)
	option2(s)

	// Check that the Services field contains all the mock services
	if len(s.Services) != 3 {
		t.Errorf("WithServices did not add services correctly with multiple calls, got %d services", len(s.Services))
	}
}

// ==================== CleanupFunc Tests ====================

// TestCleanupFunc tests the CleanupFunc type
func TestCleanupFunc(t *testing.T) {
	// Create a CleanupFunc that returns nil
	cleanupCalled := false
	cleanup := CleanupFunc(func() error {
		cleanupCalled = true
		return nil
	})

	// Call the cleanup function
	err := cleanup()
	if err != nil {
		t.Errorf("CleanupFunc returned an error: %v", err)
	}

	// Check that the cleanup function was called
	if !cleanupCalled {
		t.Error("CleanupFunc was not called")
	}
}

// TestCleanupFuncError tests the CleanupFunc type with a function that returns an error
func TestCleanupFuncError(t *testing.T) {
	// Create a CleanupFunc that returns an error
	expectedErr := errors.New("cleanup error")
	cleanup := CleanupFunc(func() error {
		return expectedErr
	})

	// Call the cleanup function
	err := cleanup()
	if err != expectedErr {
		t.Errorf("CleanupFunc returned %v, want %v", err, expectedErr)
	}
}

// TestCleanupFuncNil tests that a nil CleanupFunc can be called without panicking
func TestCleanupFuncNil(t *testing.T) {
	// Create a nil CleanupFunc
	var cleanup CleanupFunc

	// Check that calling a nil CleanupFunc panics
	defer func() {
		if r := recover(); r == nil {
			t.Error("Calling a nil CleanupFunc should panic")
		}
	}()

	// Call the nil cleanup function
	_ = cleanup()
}
