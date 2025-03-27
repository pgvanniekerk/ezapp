package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Load loads environment variables into the provided struct.
// The struct should have `envconfig` tags to specify which environment variables to load.
// For example:
//
//	type Config struct {
//		DatabaseURL string `envconfig:"DATABASE_URL" required:"true"`
//		Port        int    `envconfig:"PORT" default:"8080"`
//	}
func Load[T any](prefix string, c *T) error {
	return envconfig.Process(prefix, c)
}