package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"
)

// ShutdownCtx creates a context with a timeout specified by the EZAPP_SHUTDOWN_TIMEOUT
// environment variable (in seconds). If the variable is not set, it defaults to 15 seconds.
// If the variable contains an invalid value, it returns an error.
//
// This context is intended to be used for cleanup operations during application shutdown.
// It is a non-cancellable context that will only expire after the specified timeout.
func ShutdownCtx() (context.Context, error) {
	shutdownTimeoutStr := os.Getenv("EZAPP_SHUTDOWN_TIMEOUT")

	// Default timeout is 15 seconds
	shutdownTimeoutSec := 15

	// Parse shutdown timeout
	if shutdownTimeoutStr != "" {
		var err error
		shutdownTimeoutSec, err = strconv.Atoi(shutdownTimeoutStr)
		if err != nil {
			return nil, fmt.Errorf("invalid EZAPP_SHUTDOWN_TIMEOUT value: %s - must be an integer representing seconds", shutdownTimeoutStr)
		}
	}

	// Create a context with the shutdown timeout
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(shutdownTimeoutSec)*time.Second)

	return ctx, nil
}
