package config

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoadLogger creates a zap logger with the log level specified by the EZAPP_LOG_LEVEL
// environment variable. If the variable is not set or invalid, the default log level is INFO.
func LoadLogger() *zap.Logger {
	// Get log level from environment variable
	logLevelStr := os.Getenv("EZAPP_LOG_LEVEL")
	logLevelStr = strings.ToUpper(logLevelStr)

	// Set default log level to INFO
	logLevel := zapcore.InfoLevel

	// Parse log level from environment variable
	switch logLevelStr {
	case "DEBUG":
		logLevel = zapcore.DebugLevel
	case "INFO":
		logLevel = zapcore.InfoLevel
	case "WARN":
		logLevel = zapcore.WarnLevel
	case "ERROR":
		logLevel = zapcore.ErrorLevel
	case "DPANIC":
		logLevel = zapcore.DPanicLevel
	case "PANIC":
		logLevel = zapcore.PanicLevel
	case "FATAL":
		logLevel = zapcore.FatalLevel
	}

	// Create logger configuration
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(logLevel),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
	}

	// Build and return logger
	logger, err := config.Build()
	if err != nil {
		// If there's an error building the logger, fall back to a default logger
		return zap.NewExample()
	}

	return logger
}
