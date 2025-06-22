package config

import (
	"fmt"
	"reflect"

	"github.com/Netflix/go-env"
)

// LoadVar creates and populates a configuration struct of type CFG using environment variables.
// It validates that CFG is a struct type, creates a new instance, and populates its fields
// using the Netflix env var library based on struct tags.
// Returns an error if CFG is not a struct type or if there's an error populating the struct.
func LoadVar[CFG any]() (CFG, error) {
	var config CFG
	
	// Validate that CFG is a struct
	configType := reflect.TypeOf(config)
	if configType.Kind() != reflect.Struct {
		return config, fmt.Errorf("config type must be a struct, got %v", configType.Kind())
	}
	
	// Create a new instance of CFG
	// (Already done with var config CFG)
	
	// Use Netflix env var library to populate the struct
	_, err := env.UnmarshalFromEnviron(&config)
	if err != nil {
		return config, fmt.Errorf("failed to load configuration from environment: %w", err)
	}
	
	return config, nil
}