package observability

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// MetricsMiddleware returns a Gin middleware that instruments HTTP requests with OpenTelemetry metrics
// All metric operations are non-blocking and use async aggregation
// appName is used as the metric prefix (e.g., "ms-registration" -> "ms_registration.http.server.request.count")
func MetricsMiddleware(serviceName, appName string) gin.HandlerFunc {
	meter := otel.Meter(serviceName)

	// Normalize app name for metric prefix (replace hyphens with underscores)
	metricPrefix := normalizeMetricPrefix(appName)

	// Initialize metrics with custom prefix (async, no blocking)
	requestCounter, _ := meter.Int64Counter(
		metricPrefix+".http.server.request.count",
		metric.WithDescription("Total number of HTTP requests"),
		metric.WithUnit("{request}"),
	)

	requestDuration, _ := meter.Float64Histogram(
		metricPrefix+".http.server.request.duration",
		metric.WithDescription("HTTP request duration"),
		metric.WithUnit("ms"),
	)

	activeRequests, _ := meter.Int64UpDownCounter(
		metricPrefix+".http.server.active_requests",
		metric.WithDescription("Number of active HTTP requests"),
		metric.WithUnit("{request}"),
	)

	requestSize, _ := meter.Int64Histogram(
		metricPrefix+".http.server.request.size",
		metric.WithDescription("HTTP request body size"),
		metric.WithUnit("By"),
	)

	responseSize, _ := meter.Int64Histogram(
		metricPrefix+".http.server.response.size",
		metric.WithDescription("HTTP response body size"),
		metric.WithUnit("By"),
	)

	return func(c *gin.Context) {
		start := time.Now()

		// Get route early (remains constant)
		route := c.FullPath()
		if route == "" {
			route = "unknown" // For 404s or unmapped routes
		}
		method := c.Request.Method

		// Increment active requests (non-blocking)
		activeRequests.Add(c.Request.Context(), 1,
			metric.WithAttributes(
				attribute.String("http.method", method),
				attribute.String("http.route", route),
			),
		)

		// Record request size (non-blocking)
		if c.Request.ContentLength > 0 {
			requestSize.Record(c.Request.Context(), c.Request.ContentLength,
				metric.WithAttributes(
					attribute.String("http.method", method),
					attribute.String("http.route", route),
				),
			)
		}

		// Process request
		c.Next()

		// Calculate duration
		duration := float64(time.Since(start).Milliseconds())
		statusCode := c.Writer.Status()

		// Common attributes with endpoint and status code
		attrs := []attribute.KeyValue{
			attribute.String("http.method", method),
			attribute.String("http.route", route),
			attribute.Int("http.status_code", statusCode),
		}

		// Record metrics (all non-blocking, async aggregation)
		requestCounter.Add(c.Request.Context(), 1, metric.WithAttributes(attrs...))
		requestDuration.Record(c.Request.Context(), duration, metric.WithAttributes(attrs...))

		// Record response size with status code
		responseSize.Record(c.Request.Context(), int64(c.Writer.Size()),
			metric.WithAttributes(attrs...),
		)

		// Decrement active requests (no status_code needed here as it tracks in-flight)
		activeRequests.Add(c.Request.Context(), -1,
			metric.WithAttributes(
				attribute.String("http.method", method),
				attribute.String("http.route", route),
			),
		)
	}
}

// CustomMetrics provides a simple interface for creating custom application metrics
// All operations are non-blocking
type CustomMetrics struct {
	meter metric.Meter
}

// NewCustomMetrics creates a new custom metrics instance
func NewCustomMetrics(serviceName string) *CustomMetrics {
	return &CustomMetrics{
		meter: otel.Meter(serviceName),
	}
}

// Counter creates a new counter metric (monotonic increment only)
func (cm *CustomMetrics) Counter(name, description, unit string) (metric.Int64Counter, error) {
	return cm.meter.Int64Counter(
		name,
		metric.WithDescription(description),
		metric.WithUnit(unit),
	)
}

// UpDownCounter creates a new up/down counter (can increment and decrement)
func (cm *CustomMetrics) UpDownCounter(name, description, unit string) (metric.Int64UpDownCounter, error) {
	return cm.meter.Int64UpDownCounter(
		name,
		metric.WithDescription(description),
		metric.WithUnit(unit),
	)
}

// Histogram creates a new histogram metric (distribution of values)
func (cm *CustomMetrics) Histogram(name, description, unit string) (metric.Float64Histogram, error) {
	return cm.meter.Float64Histogram(
		name,
		metric.WithDescription(description),
		metric.WithUnit(unit),
	)
}

// Gauge creates a new observable gauge (async callback for current value)
func (cm *CustomMetrics) Gauge(name, description, unit string, callback metric.Int64Callback) error {
	_, err := cm.meter.Int64ObservableGauge(
		name,
		metric.WithDescription(description),
		metric.WithUnit(unit),
		metric.WithInt64Callback(callback),
	)
	return err
}

// FloatGauge creates a new observable float gauge (async callback for current value)
func (cm *CustomMetrics) FloatGauge(name, description, unit string, callback metric.Float64Callback) error {
	_, err := cm.meter.Float64ObservableGauge(
		name,
		metric.WithDescription(description),
		metric.WithUnit(unit),
		metric.WithFloat64Callback(callback),
	)
	return err
}

// normalizeMetricPrefix converts app names to valid metric prefixes
// Examples: "ms-registration" -> "ms_registration", "go_app_base" -> "go_app_base"
func normalizeMetricPrefix(appName string) string {
	// Replace hyphens with underscores (Prometheus naming convention)
	prefix := ""
	for _, char := range appName {
		if char == '-' {
			prefix += "_"
		} else {
			prefix += string(char)
		}
	}
	return prefix
}
