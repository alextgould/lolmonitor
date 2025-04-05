// make it easier to control log levels (DEBUG, INFO etc)
package logger

import (
	"log/slog"
	"os"
)

// Init sets the global logger with the given log level and optional JSON format.
func Init(level slog.Leveler, useJSON bool) {
	var handler slog.Handler
	if useJSON {
		handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	} else {
		handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	}

	slog.SetDefault(slog.New(handler))
}
