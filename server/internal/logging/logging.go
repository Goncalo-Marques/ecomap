package logging

import (
	"fmt"
	"log/slog"
)

var Logger *slog.Logger

// init initializes the logger with the default configurations.
func init() {
	// TODO: confirm this runs without being called
	fmt.Println("TODO: test remove me")
	Logger = slog.Default()
}

// Init initializes the logger with the provided handler.
func Init(handler slog.Handler) {
	Logger = slog.New(handler)
}
