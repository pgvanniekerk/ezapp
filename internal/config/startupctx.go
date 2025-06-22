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
//
// Additionally, it stores a shutdown timeout duration in the context value under the key
// "shutdownCtx". The shutdown timeout is controlled by the EZAPP_SHUTDOWN_TIMEOUT
// environment variable (in seconds), defaulting to 15 seconds if not set.
func StartupCtx() (context.Context, error) {
	startupTimeoutStr := os.Getenv("EZAPP_STARTUP_TIMEOUT")
	shutdownTimeoutStr := os.Getenv("EZAPP_SHUTDOWN_TIMEOUT")
	
	// Default timeouts are 15 seconds each
	startupTimeoutSec := 15
	shutdownTimeoutSec := 15
	
	// Parse startup timeout
	if startupTimeoutStr != "" {
		var err error
		startupTimeoutSec, err = strconv.Atoi(startupTimeoutStr)
		if err != nil {
			return nil, fmt.Errorf("invalid EZAPP_STARTUP_TIMEOUT value: %s - must be an integer representing seconds", startupTimeoutStr)
		}
	}
	
	// Parse shutdown timeout
	if shutdownTimeoutStr != "" {
		var err error
		shutdownTimeoutSec, err = strconv.Atoi(shutdownTimeoutStr)
		if err != nil {
			return nil, fmt.Errorf("invalid EZAPP_SHUTDOWN_TIMEOUT value: %s - must be an integer representing seconds", shutdownTimeoutStr)
		}
	}
	
	// Create a context with the startup timeout
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(startupTimeoutSec)*time.Second)
	
	// Add shutdown timeout duration as a context value
	shutdownDuration := time.Duration(shutdownTimeoutSec) * time.Second
	ctx = context.WithValue(ctx, "shutdownTimeoutDuration", shutdownDuration)
	
	return ctx, nil
}

// GetShutdownTimeout extracts the shutdown timeout duration from the context.
// If not found, it returns the default of 15 seconds.
func GetShutdownTimeout(ctx context.Context) time.Duration {
	if timeout, ok := ctx.Value("shutdownTimeoutDuration").(time.Duration); ok {
		return timeout
	}
	return 15 * time.Second
}
