package observability

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

// ConfigProvider defines the interface for observability configuration
type ConfigProvider interface {
	GetOtelEnabled() bool
	GetOtelServiceName() string
	GetJaegerEndpoint() string
	GetEnvironment() string
	GetOtelBatchTimeout() int
	GetOtelMaxExportBatchSize() int
	GetOtelMaxQueueSize() int
	GetOtelExportTimeout() int
}

// TracerProvider wraps the OpenTelemetry tracer provider
type TracerProvider struct {
	provider *sdktrace.TracerProvider
}

// NewTracerProvider initializes a new OpenTelemetry tracer provider
// If observability is disabled, returns a noop provider
func NewTracerProvider(cfg ConfigProvider) (*TracerProvider, error) {
	if !cfg.GetOtelEnabled() {
		log.Println("OpenTelemetry tracing is disabled")
		return &TracerProvider{
			provider: sdktrace.NewTracerProvider(sdktrace.WithSampler(sdktrace.NeverSample())),
		}, nil
	}

	// Create OTLP HTTP exporter for Jaeger with optimized settings
	exporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint(cfg.GetJaegerEndpoint()),
		otlptracehttp.WithInsecure(),                                 // Use insecure for local development
		otlptracehttp.WithCompression(otlptracehttp.GzipCompression), // Compress payloads
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	// Create resource with service information
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(cfg.GetOtelServiceName()),
			semconv.ServiceVersion("1.0.0"),
			semconv.DeploymentEnvironment(cfg.GetEnvironment()),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create batch span processor with optimized settings for non-blocking I/O
	batchTimeout := cfg.GetOtelBatchTimeout()
	if batchTimeout == 0 {
		batchTimeout = 5 // Default: 5 seconds
	}
	maxExportBatchSize := cfg.GetOtelMaxExportBatchSize()
	if maxExportBatchSize == 0 {
		maxExportBatchSize = 512 // Default: 512 spans
	}
	maxQueueSize := cfg.GetOtelMaxQueueSize()
	if maxQueueSize == 0 {
		maxQueueSize = 2048 // Default: 2048 spans
	}
	exportTimeout := cfg.GetOtelExportTimeout()
	if exportTimeout == 0 {
		exportTimeout = 30 // Default: 30 seconds
	}

	batchProcessor := sdktrace.NewBatchSpanProcessor(
		exporter,
		sdktrace.WithBatchTimeout(time.Duration(batchTimeout)*time.Second),
		sdktrace.WithMaxExportBatchSize(maxExportBatchSize),
		sdktrace.WithMaxQueueSize(maxQueueSize),
		sdktrace.WithExportTimeout(time.Duration(exportTimeout)*time.Second),
	)

	// Create tracer provider with optimized batching
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batchProcessor),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // Sample all traces in development
	)

	// Set global tracer provider
	otel.SetTracerProvider(tp)

	// Set global propagator for context propagation (W3C Trace Context)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	log.Printf("OpenTelemetry tracing initialized: service=%s, endpoint=%s", cfg.GetOtelServiceName(), cfg.GetJaegerEndpoint())

	return &TracerProvider{
		provider: tp,
	}, nil
}

// Tracer returns a named tracer
func (tp *TracerProvider) Tracer(name string) trace.Tracer {
	return tp.provider.Tracer(name)
}

// Shutdown gracefully shuts down the tracer provider
// This ensures all spans are flushed before application exit
func (tp *TracerProvider) Shutdown(ctx context.Context) error {
	if tp.provider == nil {
		return nil
	}

	log.Println("Shutting down OpenTelemetry tracer provider...")
	if err := tp.provider.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown tracer provider: %w", err)
	}

	log.Println("OpenTelemetry tracer provider shut down successfully")
	return nil
}

// GetProvider returns the underlying SDK tracer provider
func (tp *TracerProvider) GetProvider() *sdktrace.TracerProvider {
	return tp.provider
}
