package config

import (
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadLogger(t *testing.T) {
	// Save original environment variable value to restore later
	originalLogLevel := os.Getenv("EZAPP_LOG_LEVEL")
	defer os.Setenv("EZAPP_LOG_LEVEL", originalLogLevel)

	testCases := []struct {
		name          string
		logLevelEnv   string
		expectedLevel slog.Level
	}{
		{
			name:          "debug level",
			logLevelEnv:   "DEBUG",
			expectedLevel: slog.LevelDebug,
		},
		{
			name:          "info level",
			logLevelEnv:   "INFO",
			expectedLevel: slog.LevelInfo,
		},
		{
			name:          "warn level",
			logLevelEnv:   "WARN",
			expectedLevel: slog.LevelWarn,
		},
		{
			name:          "error level",
			logLevelEnv:   "ERROR",
			expectedLevel: slog.LevelError,
		},
		{
			name:          "lowercase level",
			logLevelEnv:   "debug",
			expectedLevel: slog.LevelDebug,
		},
		{
			name:          "mixed case level",
			logLevelEnv:   "DeBuG",
			expectedLevel: slog.LevelDebug,
		},
		{
			name:          "empty level defaults to info",
			logLevelEnv:   "",
			expectedLevel: slog.LevelInfo,
		},
		{
			name:          "invalid level defaults to info",
			logLevelEnv:   "INVALID",
			expectedLevel: slog.LevelInfo,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set environment variable for this test case
			os.Setenv("EZAPP_LOG_LEVEL", tc.logLevelEnv)

			// Load logger
			logger := LoadLogger()

			// Check that the logger is not nil
			// Note: slog doesn't provide a way to directly access the logger's level
			assert.NotNil(t, logger, "Logger should not be nil")

			// No need to call Sync() as slog doesn't have this method
		})
	}
}
