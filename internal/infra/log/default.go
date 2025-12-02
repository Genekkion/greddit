package log

import (
	"log/slog"
	"os"
)

var (
	// defaultLoggerCleanups are functions to be called when the default logger is cleaned up.
	defaultLoggerCleanups []func()

	// defaultLogger is the "global" logger used when no logger is provided.
	defaultLogger = func() *slog.Logger {
		stdOutLogger := newStdOutLogger()

		handlers := []slog.Handler{
			NewHandler(os.Stdout, nil),
		}

		logPath, exists := os.LookupEnv("LOG_FILE")
		if exists {
			file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				stdOutLogger.Error("Unable to open log file",
					"path", logPath,
					"error", err,
				)
				os.Exit(1)
			}

			handlers = append(handlers, NewHandler(file, nil))
			defaultLoggerCleanups = append(defaultLoggerCleanups, func() {
				err := file.Close()
				if err != nil {
					stdOutLogger.Error("Unable to close log file",
						"path", logPath,
						"error", err,
					)
					os.Exit(2)
				}
			})
		}

		return NewLogger(handlers...)
	}()
)

// CleanupDefaultLogger cleans up the default logger.
func CleanupDefaultLogger() {
	for _, f := range defaultLoggerCleanups {
		f()
	}
}

// newStdOutLogger creates a new logger that writes to stdout.
func newStdOutLogger() *slog.Logger {
	return NewLogger(NewHandler(os.Stdout, nil))
}

// SetDefaultLogger sets the default logger.
func SetDefaultLogger(logger *slog.Logger) {
	if logger == nil {
		panic("Default logger cannot be nil")
	}
	defaultLogger = logger
}

// GetDefaultLogger returns the default logger.
func GetDefaultLogger() *slog.Logger {
	return defaultLogger
}
