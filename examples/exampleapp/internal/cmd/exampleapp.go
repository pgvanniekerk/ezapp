package main

import (
	"context"

	"github.com/pgvanniekerk/ezapp/examples/exampleapp/internal/config"
	"github.com/pgvanniekerk/ezapp/examples/exampleapp/internal/wire"
	"github.com/pgvanniekerk/ezapp/pkg/ezapp"
)

// This example application demonstrates how to use the ezapp framework to build
// a simple application with a database connection. It shows the recommended
// structure and patterns for applications built with ezapp.
//
// The application follows these key patterns:
// 1. Define a Config struct with fields for all application configuration
// 2. Create a BuildFunc that takes a context and config, and returns an EzApp
// 3. Use the wire package to create the application with its dependencies
// 4. Use ezapp.Run to start the application

// Config holds all configuration for the application.
// This struct demonstrates different ways to configure fields using struct tags:
// - Default values
// - Required fields
// - Validation rules
//
// In a real application, this struct would be populated from environment variables
// or other configuration sources before being passed to the BuildFunc.
type Config struct {
	// Database configuration for connecting to PostgreSQL
	DB config.DBConf

	// Example of an environment variable with a default value.
	// If LOG_LEVEL is not set, "info" will be used.
	LogLevel string `envvar:"LOG_LEVEL" default:"info"`

	// Example of a required environment variable.
	// The application will fail to start if APP_NAME is not set.
	AppName string `envvar:"APP_NAME" required:"true"`

	// Example of an environment variable with validation.
	// The value must be between 1024 and 65535 (valid port range).
	Port int `envvar:"PORT" default:"8080" validate:"min=1024,max=65535"`
}

// BuildFunc is the application builder function that creates and configures
// the application with all its dependencies.
//
// This function implements the ezapp.Builder interface and is responsible for:
// 1. Creating all application components (database connections, services, etc.)
// 2. Wiring them together using the wire package
// 3. Returning an ezapp.EzApp instance that can be run
//
// Parameters:
//   - startupCtx: A context that can be used for initialization operations
//   - cfg: The application configuration
//
// Returns:
//   - An ezapp.EzApp instance that can be run
//   - An error if the application could not be created
func BuildFunc(startupCtx context.Context, cfg Config) (ezapp.EzApp, error) {
	// Wire up the application with all dependencies
	return wire.Wire(startupCtx, cfg.DB)
}

// main is the entry point for the application.
// It uses ezapp.Run to start the application with the BuildFunc.
//
// ezapp.Run will:
// 1. Create a background context
// 2. Create an empty Config (in a real application, you would load this from environment variables)
// 3. Call BuildFunc to create the application
// 4. Run the application
func main() {
	// Run the application with the BuildFunc
	ezapp.Run(BuildFunc)
}
