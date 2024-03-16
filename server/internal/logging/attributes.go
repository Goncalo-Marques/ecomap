package logging

import "log/slog"

// Error returns a logging attribute for an error.
func Error(v error) slog.Attr {
	return slog.String("error", v.Error())
}
