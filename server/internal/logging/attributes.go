package logging

import "log/slog"

// Error returns a logging attribute for an error.
func Error(v error) slog.Attr {
	errMessage := "nil"
	if v != nil {
		errMessage = v.Error()
	}

	return slog.String("error", errMessage)
}
