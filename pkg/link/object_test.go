package link

import (
	"testing"

	"github.com/pgvanniekerk/ezapp/pkg/ezapp"
	"github.com/stretchr/testify/assert"
	"go.uber.org/dig"
)

func TestObject(t *testing.T) {
	// Test object
	testObj := "test string"

	// Create a build process with the test object and nil Option
	buildProc := Object(testObj, nil)

	// Create a container and apply the build process
	container := ezapp.Construct(buildProc)

	// Verify that the object was provided to the container
	err := container.Invoke(func(s string) {
		assert.Equal(t, testObj, s, "The provided object should match the original object")
	})
	assert.NoError(t, err, "Invoking the container with the provided type should not error")
}

func TestObjectWithOptions(t *testing.T) {

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

func TestObjectWithMultipleOptions(t *testing.T) {
	// Test object
	testObj := "test string"

	// Create a build process with the test object and a named Option
	namedProc := Object(testObj, Named("test-name"))

	// Create a build process with the test object and a grouped Option
	groupedProc := Object(testObj, Grouped("test-group"))

	// Create a container and apply both build processes
	container := ezapp.Construct(namedProc, groupedProc)

	// Define a struct with fields that have name and group tags
	type NamedAndGroupedDep struct {
		dig.In
		NamedString    string   `name:"test-name"`
		GroupedStrings []string `group:"test-group"`
	}

	// Verify that the object was provided to the container with the correct name and group
	err := container.Invoke(func(deps NamedAndGroupedDep) {
		assert.Equal(t, testObj, deps.NamedString, "The provided named object should match the original object")
		assert.Len(t, deps.GroupedStrings, 1, "There should be one string in the group")
		assert.Equal(t, testObj, deps.GroupedStrings[0], "The provided grouped object should match the original object")
	})
	assert.NoError(t, err, "Invoking the container with the provided named and grouped type should not error")
}
