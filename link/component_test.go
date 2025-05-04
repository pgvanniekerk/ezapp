package link

import (
	"context"
	"github.com/pgvanniekerk/ezapp/internal/primitive"
	"github.com/stretchr/testify/suite"
	"go.uber.org/dig"
	"testing"
)

// Ensure we're using the primitive package
var _ primitive.Component[MockParams] = (*MockComponent)(nil)

// Mock component and params for testing
type MockComponent struct {
	InitCalled    bool
	CleanupCalled bool
	InitError     error
	CleanupError  error
	Params        MockParams
	InitContext   context.Context
	// Embed primitive.Component to satisfy the validation
	primitive.Component[MockParams]
}

func (m *MockComponent) Init(ctx context.Context, params MockParams) error {
	m.InitCalled = true
	m.Params = params
	m.InitContext = ctx
	return m.InitError
}

func (m *MockComponent) Cleanup(ctx context.Context) error {
	m.CleanupCalled = true
	return m.CleanupError
}

// MockParams is a struct type for testing
type MockParams struct {
	Value1 string
	Value2 int
}

// NonStructParams is a non-struct type for testing error case
type NonStructParams int

// MockComponentForNonStruct is a mock component that implements primitive.Component[NonStructParams]
// This is used to test the error case when Params is not a struct
type MockComponentForNonStruct struct{}

func (m *MockComponentForNonStruct) Init(ctx context.Context, params NonStructParams) error {
	return nil
}

func (m *MockComponentForNonStruct) Cleanup(ctx context.Context) error {
	return nil
}

// MockComponentWithoutEmbedding is a mock component that implements primitive.Component[MockParams]
// but doesn't embed the primitive.Component field. This is used to test the error case
// when Comp doesn't embed primitive.Component.
type MockComponentWithoutEmbedding struct{}

// Ensure MockComponentWithoutEmbedding implements primitive.Component[MockParams]
var _ primitive.Component[MockParams] = (*MockComponentWithoutEmbedding)(nil)

func (m *MockComponentWithoutEmbedding) Init(ctx context.Context, params MockParams) error {
	return nil
}

func (m *MockComponentWithoutEmbedding) Cleanup(ctx context.Context) error {
	return nil
}

// ComponentTestSuite is a test suite for the Component function
type ComponentTestSuite struct {
	suite.Suite
	container *dig.Container
}

func (s *ComponentTestSuite) SetupTest() {
	// Create a new dig container for each test
	s.container = dig.New()
}

func (s *ComponentTestSuite) TestComponentRegistration() {
	// Create test parameters
	testParams := MockParams{
		Value1: "test value",
		Value2: 42,
	}

	// Provide the individual fields of MockParams to the container
	err := s.container.Provide(func() string {
		return testParams.Value1
	})
	s.NoError(err, "Value1 provider registration should succeed")

	err = s.container.Provide(func() int {
		return testParams.Value2
	})
	s.NoError(err, "Value2 provider registration should succeed")

	// Provide a context with the name "ezapp_initCtx" to the container
	err = s.container.Provide(func() context.Context {
		return context.Background()
	}, dig.Name("ezapp_initCtx"))
	s.NoError(err, "Context provider registration should succeed")

	// Test successful component registration
	err = Component[*MockComponent, MockParams](s.container)
	s.NoError(err, "Component registration should succeed")

	// Define a custom type to avoid conflicts with existing providers
	type TestResult struct {
		Success bool
	}

	// Verify that the component was registered by checking if we can provide another component
	// that depends on it without error
	err = s.container.Provide(func(c *MockComponent) TestResult {
		return TestResult{Success: true}
	})
	s.NoError(err, "Dependent provider registration should succeed")
}

func (s *ComponentTestSuite) TestComponentRegistrationWithNonStructParams() {
	// Test error case when Params is not a struct
	err := Component[*MockComponentForNonStruct, NonStructParams](s.container)
	s.Error(err, "Component registration should fail with non-struct Params")
	s.Contains(err.Error(), "must be a struct", "Error message should indicate that Params must be a struct")
}

func (s *ComponentTestSuite) TestComponentRegistrationWithoutEmbedding() {
	// Test error case when Comp doesn't embed primitive.Component
	err := Component[*MockComponentWithoutEmbedding, MockParams](s.container)
	s.Error(err, "Component registration should fail when Comp doesn't embed primitive.Component")
	s.Contains(err.Error(), "must embed a field of type primitive.Component", "Error message should indicate that Comp must embed primitive.Component")
}

func (s *ComponentTestSuite) TestComponentInvocation() {
	// Test that we can invoke the type from the dig container
	// This verifies that the Component function correctly registers the component

	// Create a new container for this test to avoid interference with other tests
	container := dig.New()

	// Create test parameters
	testParams := MockParams{
		Value1: "test value",
		Value2: 42,
	}

	// Provide the individual fields of MockParams to the container
	err := container.Provide(func() string {
		return testParams.Value1
	})
	s.NoError(err, "Value1 provider registration should succeed")

	err = container.Provide(func() int {
		return testParams.Value2
	})
	s.NoError(err, "Value2 provider registration should succeed")

	// Create a custom context with a value that can be verified
	type contextKey string
	const testKey contextKey = "testKey"
	testValue := "testValue"
	testCtx := context.WithValue(context.Background(), testKey, testValue)

	// Provide the custom context with the name "ezapp_initCtx" to the container
	err = container.Provide(func() context.Context {
		return testCtx
	}, dig.Name("ezapp_initCtx"))
	s.NoError(err, "Context provider registration should succeed")

	// Register the component with the container using the Component function
	err = Component[*MockComponent, MockParams](container)
	s.NoError(err, "Component registration should succeed")

	// Now invoke a function that tries to receive the component from the container
	// This should succeed if the Component function correctly registered the component
	var componentFromContainer *MockComponent
	err = container.Invoke(func(component *MockComponent) {
		componentFromContainer = component
	})

	// Assert that there was no error retrieving the component
	s.NoError(err, "Component should be successfully retrieved from the container")

	// Verify that the component is not nil
	s.NotNil(componentFromContainer, "Component should not be nil")

	// Verify that the component was initialized
	s.True(componentFromContainer.InitCalled, "Component's Init method should have been called")

	// Verify that the component's parameters were correctly set
	s.Equal(testParams.Value1, componentFromContainer.Params.Value1, "Value1 should match the provided value")
	s.Equal(testParams.Value2, componentFromContainer.Params.Value2, "Value2 should match the provided value")

	// Verify that the context passed to Init is the same as the one provided to the container
	s.NotNil(componentFromContainer.InitContext, "InitContext should not be nil")
	s.Equal(testValue, componentFromContainer.InitContext.Value(testKey), "InitContext should contain the test value")
}

// Run the test suite
func TestComponentSuite(t *testing.T) {
	suite.Run(t, new(ComponentTestSuite))
}
