package logger

import (
	"context"
	"log/slog"
	"os"
	"time"

	sharedLogger "github.com/refortunato/go_app_base/internal/shared/logger"
)

// SlogLogger is a concrete implementation of Logger interface using Go's log/slog package.
type SlogLogger struct {
	logger      *slog.Logger
	imageName   string
	imageVer    string
	baseAttrs   []slog.Attr
	contextData sharedLogger.CustomFields
}

// NewSlogLogger creates a new logger instance configured to output JSON to STDOUT.
// It includes imageName and imageVersion in all log entries.
func NewSlogLogger(imageName, imageVersion string) sharedLogger.Logger {
	// Create a custom JSON handler that writes to STDOUT
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Customize timestamp format to include microseconds
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					// Format: 2006-01-02T15:04:05.000000Z07:00
					return slog.String("timestamp", t.Format("2006-01-02T15:04:05.000000Z07:00"))
				}
			}
			return a
		},
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)

	return &SlogLogger{
		logger:      logger,
		imageName:   imageName,
		imageVer:    imageVersion,
		baseAttrs:   []slog.Attr{},
		contextData: make(sharedLogger.CustomFields),
	}
}

// Debug logs a debug-level message
func (l *SlogLogger) Debug(message string, customFields ...sharedLogger.CustomFields) {
	l.log(slog.LevelDebug, message, customFields...)
}

// Info logs an info-level message
func (l *SlogLogger) Info(message string, customFields ...sharedLogger.CustomFields) {
	l.log(slog.LevelInfo, message, customFields...)
}

// Warn logs a warning-level message
func (l *SlogLogger) Warn(message string, customFields ...sharedLogger.CustomFields) {
	l.log(slog.LevelWarn, message, customFields...)
}

// Error logs an error-level message
func (l *SlogLogger) Error(message string, customFields ...sharedLogger.CustomFields) {
	l.log(slog.LevelError, message, customFields...)
}

// With creates a new logger instance with additional context fields
func (l *SlogLogger) With(fields sharedLogger.CustomFields) sharedLogger.Logger {
	// Create a new logger with merged context
	newContextData := make(sharedLogger.CustomFields)
	for k, v := range l.contextData {
		newContextData[k] = v
	}
	for k, v := range fields {
		newContextData[k] = v
	}

	return &SlogLogger{
		logger:      l.logger,
		imageName:   l.imageName,
		imageVer:    l.imageVer,
		baseAttrs:   l.baseAttrs,
		contextData: newContextData,
	}
}

// log is the internal method that performs the actual logging
func (l *SlogLogger) log(level slog.Level, message string, customFields ...sharedLogger.CustomFields) {
	ctx := context.Background()

	// Build the list of attributes
	attrs := []any{}

	// Add imageName and imageVersion
	attrs = append(attrs, slog.String("imageName", l.imageName))
	attrs = append(attrs, slog.String("imageVersion", l.imageVer))

	// Merge contextData with custom fields
	mergedCustom := make(sharedLogger.CustomFields)
	for k, v := range l.contextData {
		mergedCustom[k] = v
	}
	for _, cf := range customFields {
		for k, v := range cf {
			mergedCustom[k] = v
		}
	}

	// Add custom fields as a nested group if present
	if len(mergedCustom) > 0 {
		customAttrs := make([]any, 0, len(mergedCustom))
		for k, v := range mergedCustom {
			customAttrs = append(customAttrs, slog.Any(k, v))
		}
		attrs = append(attrs, slog.Group("custom", customAttrs...))
	}

	// Log with the appropriate level
	l.logger.Log(ctx, level, message, attrs...)
}
