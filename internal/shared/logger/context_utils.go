package logger

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// ExtractTraceContext extracts trace and span IDs from context using OpenTelemetry.
// Returns empty strings if context is nil or span context is not valid.
func ExtractTraceContext(ctx context.Context) (traceID, spanID string) {
	if ctx == nil {
		return "", ""
	}

	spanCtx := trace.SpanContextFromContext(ctx)
	if !spanCtx.IsValid() {
		return "", ""
	}

	return spanCtx.TraceID().String(), spanCtx.SpanID().String()
}

// ExtractCustomContextFields extracts all relevant observability fields from context.
// This includes trace information from OpenTelemetry and any custom context values.
func ExtractCustomContextFields(ctx context.Context) CustomFields {
	fields := make(CustomFields)

	if ctx == nil {
		return fields
	}

	// Extract OpenTelemetry trace information
	traceID, spanID := ExtractTraceContext(ctx)
	if traceID != "" {
		fields["traceId"] = traceID
	}
	if spanID != "" {
		fields["spanId"] = spanID
	}

	// Future: Add extraction of custom context values
	// Example:
	// if userId := ctx.Value(userIDKey); userId != nil {
	//     fields["userId"] = userId
	// }
	//
	// if requestId := ctx.Value(requestIDKey); requestId != nil {
	//     fields["requestId"] = requestId
	// }

	return fields
}
