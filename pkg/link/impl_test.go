package link

import (
	"testing"

	"github.com/pgvanniekerk/ezapp/pkg/ezapp"
	"github.com/stretchr/testify/assert"
	"go.uber.org/dig"
)

// TestInterface is an interface used for testing the Impl function
type TestInterface interface {
	TestMethod() string
}

// TestImplementation is a struct that implements TestInterface
type TestImplementation struct {
	Value string
}

// TestMethod implements the TestInterface
func (t *TestImplementation) TestMethod() string {
	return t.Value
}

func TestImpl(t *testing.T) {
	// Create a test implementation
	testImpl := &TestImplementation{Value: "test value"}

	// Create a build process with the test implementation and the Impl Option
	buildProc := Object(testImpl, Impl[TestInterface](func(opts *[]dig.ProvideOption) {}))

	// Create a container and apply the build process
	container := ezapp.Construct(buildProc)

	// Verify that the implementation was provided to the container as the interface
	err := container.Invoke(func(iface TestInterface) {
		assert.Equal(t, "test value", iface.TestMethod(), "The provided interface should call the implementation's method")
	})
	assert.NoError(t, err, "Invoking the container with the provided interface should not error")
}

func TestImplWithNamedOption(t *testing.T) {
	// Create a test implementation
	testImpl := &TestImplementation{Value: "test value"}

	// Create a build process with the test implementation and both Impl and Named options
	buildProc := Object(testImpl, Impl[TestInterface](Named("test-name")))

	// Create a container and apply the build process
	container := ezapp.Construct(buildProc)

	// Define a struct with a field that has a name tag
	type NamedDep struct {
		dig.In
		NamedInterface TestInterface `name:"test-name"`
	}

	// Verify that the implementation was provided to the container as the named interface
	err := container.Invoke(func(deps NamedDep) {
		assert.Equal(t, "test value", deps.NamedInterface.TestMethod(), "The provided named interface should call the implementation's method")
	})
	assert.NoError(t, err, "Invoking the container with the provided named interface should not error")
}
