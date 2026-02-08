package logger

import (
	"context"
	"sync"
)

// Logger interface for structured logging across the application.
// All log methods accept a context as the first parameter for trace propagation,
// followed by a message and optional custom fields.
type Logger interface {
	// Debug logs a debug-level message with optional custom fields
	Debug(ctx context.Context, message string, customFields ...CustomFields)

	// Info logs an info-level message with optional custom fields
	Info(ctx context.Context, message string, customFields ...CustomFields)

	// Warn logs a warning-level message with optional custom fields
	Warn(ctx context.Context, message string, customFields ...CustomFields)

	// Error logs an error-level message with optional custom fields
	Error(ctx context.Context, message string, customFields ...CustomFields)

	// With creates a new logger instance with additional context fields
	// that will be included in all subsequent log entries
	With(fields CustomFields) Logger
}

// CustomFields represents additional structured data to be included in log entries.
// These fields will appear under the "custom" key in the JSON output.
type CustomFields map[string]any

// Global logger instance
var (
	globalLogger Logger
	mu           sync.RWMutex
)

// SetGlobalLogger sets the global logger instance.
// This should be called once during application initialization (e.g., in dependencies.InitDependencies).
func SetGlobalLogger(logger Logger) {
	mu.Lock()
	defer mu.Unlock()
	globalLogger = logger
}

// getLogger returns the global logger instance.
// If no logger has been set, it panics (fail-fast during development).
func getLogger() Logger {
	mu.RLock()
	defer mu.RUnlock()
	if globalLogger == nil {
		panic("logger not initialized: call logger.SetGlobalLogger first")
	}
	return globalLogger
}

// Debug logs a debug-level message using the global logger.
func Debug(ctx context.Context, message string, customFields ...CustomFields) {
	getLogger().Debug(ctx, message, customFields...)
}

// Info logs an info-level message using the global logger.
func Info(ctx context.Context, message string, customFields ...CustomFields) {
	getLogger().Info(ctx, message, customFields...)
}

// Warn logs a warning-level message using the global logger.
func Warn(ctx context.Context, message string, customFields ...CustomFields) {
	getLogger().Warn(ctx, message, customFields...)
}

// Error logs an error-level message using the global logger.
func Error(ctx context.Context, message string, customFields ...CustomFields) {
	getLogger().Error(ctx, message, customFields...)
}

// With creates a new logger instance with additional context fields.
func With(fields CustomFields) Logger {
	return getLogger().With(fields)
}
