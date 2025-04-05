package conf

import (
	"github.com/kelseyhightower/envconfig"
)

// LoadAppConf loads the application configuration from environment variables.
// This function uses the envconfig package to populate the AppConf struct
// with values from environment variables.
//
// The function adds the prefix "EZAPP_" to all environment variable names
// defined in the AppConf struct. For example, the ShutdownTimeout field
// will be populated from the EZAPP_SHUTDOWN_TIMEOUT environment variable.
//
// Returns:
//   - An AppConf struct populated with values from environment variables
//   - An error if the configuration could not be loaded
//
// Example:
//
//	conf, err := LoadAppConf()
//	if err != nil {
//	    log.Fatalf("Failed to load configuration: %v", err)
//	}
//	fmt.Printf("Shutdown timeout: %v\n", conf.ShutdownTimeout)
//
// Potential errors:
//   - Environment variables with invalid values (e.g., non-numeric values for numeric fields)
//   - Required fields that are not set (if any)
func LoadAppConf() (AppConf, error) {
	var conf AppConf
	err := envconfig.Process("EZAPP", &conf)
	return conf, err
}
