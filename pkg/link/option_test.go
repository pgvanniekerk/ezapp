package link

import (
	"testing"

	"github.com/pgvanniekerk/ezapp/pkg/ezapp"
	"github.com/stretchr/testify/assert"
	"go.uber.org/dig"
)

func TestOption(t *testing.T) {
	// Create a custom Option function
	customOpt := func(opts *[]dig.ProvideOption) {
		*opts = append(*opts, dig.Name("custom-name"))
	}

	// Test object
	testObj := "test string"

	// Create a build process with the test object and the custom Option
	buildProc := Object(testObj, customOpt)

	// Create a container and apply the build process
	container := ezapp.Construct(buildProc)

	// Define a struct with a field that has a name tag
	type NamedDep struct {
		dig.In
		NamedString string `name:"custom-name"`
	}

	// Verify that the object was provided to the container with the correct name
	err := container.Invoke(func(deps NamedDep) {
		assert.Equal(t, testObj, deps.NamedString, "The provided named object should match the original object")
	})
	assert.NoError(t, err, "Invoking the container with the provided named type should not error")
}

func TestOptionComposition(t *testing.T) {
	// Create two separate Option functions
	namedOpt := Named("composite-name")
	groupedOpt := Grouped("composite-group")

	// Test object
	testObj := "test string"

	// Create a build process with the test object and the named Option
	namedProc := Object(testObj, namedOpt)

	// Create a build process with the test object and the grouped Option
	groupedProc := Object(testObj, groupedOpt)

	// Create a container and apply both build processes
	container := ezapp.Construct(namedProc, groupedProc)

	// Define a struct with a field that has a name tag
	type NamedDep struct {
		dig.In
		NamedString string `name:"composite-name"`
	}

	// Verify that the object was provided to the container with the correct name
	err := container.Invoke(func(deps NamedDep) {
		assert.Equal(t, testObj, deps.NamedString, "The provided named object should match the original object")
	})
	assert.NoError(t, err, "Invoking the container with the provided named type should not error")

	// Define a struct with a field that has a group tag
	type GroupedDep struct {
		dig.In
		GroupedStrings []string `group:"composite-group"`
	}

	// Verify that the object was provided to the container with the correct group
	err = container.Invoke(func(deps GroupedDep) {
		assert.Len(t, deps.GroupedStrings, 1, "There should be one string in the group")
		assert.Equal(t, testObj, deps.GroupedStrings[0], "The provided grouped object should match the original object")
	})
	assert.NoError(t, err, "Invoking the container with the provided grouped type should not error")
}
