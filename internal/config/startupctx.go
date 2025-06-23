package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"
)

// StartupCtx creates a context with a timeout specified by the EZAPP_STARTUP_TIMEOUT
// environment variable (in seconds). If the variable is not set, it defaults to 15 seconds.
// If the variable contains an invalid value, it returns an error.
func StartupCtx() (context.Context, error) {
	startupTimeoutStr := os.Getenv("EZAPP_STARTUP_TIMEOUT")

	// Default timeout is 15 seconds
	startupTimeoutSec := 15

	// Parse startup timeout
	if startupTimeoutStr != "" {
		var err error
		startupTimeoutSec, err = strconv.Atoi(startupTimeoutStr)
		if err != nil {
			return nil, fmt.Errorf("invalid EZAPP_STARTUP_TIMEOUT value: %s - must be an integer representing seconds", startupTimeoutStr)
		}
	}

	// Create a context with the startup timeout
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(startupTimeoutSec)*time.Second)

	return ctx, nil
}
