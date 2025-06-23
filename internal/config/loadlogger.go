package config

import (
	"log/slog"
	"os"
	"strings"
)

// LoadLogger creates a slog logger with the log level specified by the EZAPP_LOG_LEVEL
// environment variable. If the variable is not set or invalid, the default log level is INFO.
func LoadLogger() *slog.Logger {

	// Get log level from environment variable
	logLevelStr := os.Getenv("EZAPP_LOG_LEVEL")
	logLevelStr = strings.ToUpper(logLevelStr)

	// Set default log level to INFO
	var logLevel slog.Level

	// Parse log level from environment variable
	switch logLevelStr {
	case "DEBUG":
		logLevel = slog.LevelDebug
	case "INFO":
		logLevel = slog.LevelInfo
	case "WARN":
		logLevel = slog.LevelWarn
	case "ERROR":
		logLevel = slog.LevelError
	default:
		// Default to INFO for invalid or empty values
		logLevel = slog.LevelInfo
	}

	// Create JSON handler with the configured level
	opts := &slog.HandlerOptions{
		Level: logLevel,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)

	// Create and return logger
	return slog.New(handler)
}
