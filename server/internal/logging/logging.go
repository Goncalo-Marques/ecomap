package logging

import (
	"log/slog"
)

// Logger defines the application logger.
var Logger *slog.Logger

// init initializes the logger with the default configurations.
func init() {
	Logger = slog.Default()
}

// Init initializes the logger with the provided handler and base attributes.
func Init(handler slog.Handler, attrs ...any) {
	Logger = slog.New(handler)
	Logger = Logger.With(attrs...)
}
