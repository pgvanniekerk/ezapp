package link

import (
	"testing"

	"github.com/pgvanniekerk/ezapp/pkg/ezapp"
	"github.com/stretchr/testify/assert"
	"go.uber.org/dig"
)

func TestGrouped(t *testing.T) {
	// Create a grouped Option
	groupedOpt := Grouped("test-group")

	// Create a slice of dig.ProvideOption
	opts := make([]dig.ProvideOption, 0)

	// Apply the grouped Option to the slice
	groupedOpt(&opts)

	// Test that the Option was applied correctly by using it in a build process
	testObj := "test string"
	buildProc := Object(testObj, groupedOpt)

	// Create a container and apply the build process
	container := ezapp.Construct(buildProc)

	// Define a struct with a field that has a group tag
	type GroupedDep struct {
		dig.In
		GroupedStrings []string `group:"test-group"`
	}

	// Verify that the object was provided to the container with the correct group
	err := container.Invoke(func(deps GroupedDep) {
		assert.Len(t, deps.GroupedStrings, 1, "There should be one string in the group")
		assert.Equal(t, testObj, deps.GroupedStrings[0], "The provided grouped object should match the original object")
	})
	assert.NoError(t, err, "Invoking the container with the provided grouped type should not error")
}

func TestGroupedDirectly(t *testing.T) {
	// Test object
	testObj := "test string"

	// Create a build process with the test object and a grouped Option
	buildProc := Object(testObj, Grouped("test-group"))

	// Create a container and apply the build process
	container := ezapp.Construct(buildProc)

	// Define a struct with a field that has a group tag
	type GroupedDep struct {
		dig.In
		GroupedStrings []string `group:"test-group"`
	}

	// Verify that the object was provided to the container with the correct group
	err := container.Invoke(func(deps GroupedDep) {
		assert.Len(t, deps.GroupedStrings, 1, "There should be one string in the group")
		assert.Equal(t, testObj, deps.GroupedStrings[0], "The provided grouped object should match the original object")
	})
	assert.NoError(t, err, "Invoking the container with the provided grouped type should not error")
}
