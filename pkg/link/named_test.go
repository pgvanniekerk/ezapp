package link

import (
	"testing"

	"github.com/pgvanniekerk/ezapp/pkg/ezapp"
	"github.com/stretchr/testify/assert"
	"go.uber.org/dig"
)

func TestNamed(t *testing.T) {
	// Create a named Option
	namedOpt := Named("test-name")

	// Create a slice of dig.ProvideOption
	opts := make([]dig.ProvideOption, 0)

	// Apply the named Option to the slice
	namedOpt(&opts)

	// Test that the Option was applied correctly by using it in a build process
	testObj := "test string"
	buildProc := Object(testObj, func(o *[]dig.ProvideOption) {
		for _, opt := range opts {
			*o = append(*o, opt)
		}
	})

	// Create a container and apply the build process
	container := ezapp.Construct(buildProc)

	// Define a struct with a field that has a name tag
	type NamedDep struct {
		dig.In
		NamedString string `name:"test-name"`
	}

	// Verify that the object was provided to the container with the correct name
	err := container.Invoke(func(deps NamedDep) {
		assert.Equal(t, testObj, deps.NamedString, "The provided named object should match the original object")
	})
	assert.NoError(t, err, "Invoking the container with the provided named type should not error")
}

func TestNamedDirectly(t *testing.T) {
	// Test object
	testObj := "test string"

	// Create a build process with the test object and a named Option
	buildProc := Object(testObj, Named("test-name"))

	// Create a container and apply the build process
	container := ezapp.Construct(buildProc)

	// Define a struct with a field that has a name tag
	type NamedDep struct {
		dig.In
		NamedString string `name:"test-name"`
	}

	// Verify that the object was provided to the container with the correct name
	err := container.Invoke(func(deps NamedDep) {
		assert.Equal(t, testObj, deps.NamedString, "The provided named object should match the original object")
	})
	assert.NoError(t, err, "Invoking the container with the provided named type should not error")
}
