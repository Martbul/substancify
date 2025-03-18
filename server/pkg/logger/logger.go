package logger

import (
	"github.com/hashicorp/go-hclog"
)

// Logger is a globally accessible logger instance.
var Logger hclog.Logger

func Init() {
	Logger = hclog.New(&hclog.LoggerOptions{
		Output:     hclog.DefaultOutput,            // Standard output (can also log to a file, etc.)
		Level:      hclog.LevelFromString("DEBUG"), // Log level (DEBUG, INFO, WARN, ERROR, etc.)
		JSONFormat: false,                          // Enable JSON output for structured logging
	})
}

func GetLogger() hclog.Logger {
	return Logger
}
