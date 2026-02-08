package observability

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// MeterProvider wraps the OpenTelemetry meter provider
type MeterProvider struct {
	provider *sdkmetric.MeterProvider
}

// NewMeterProvider initializes a new OpenTelemetry meter provider
// If observability is disabled, returns a noop provider
// Uses non-blocking batch processing to avoid I/O overhead
func NewMeterProvider(cfg ConfigProvider) (*MeterProvider, error) {
	if !cfg.GetOtelEnabled() {
		log.Println("OpenTelemetry metrics is disabled")
		return &MeterProvider{
			provider: sdkmetric.NewMeterProvider(),
		}, nil
	}

	// Create OTLP HTTP exporter for metrics with compression
	exporter, err := otlpmetrichttp.New(
		context.Background(),
		otlpmetrichttp.WithEndpoint(cfg.GetJaegerEndpoint()),
		otlpmetrichttp.WithInsecure(),
		otlpmetrichttp.WithCompression(otlpmetrichttp.GzipCompression),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP metric exporter: %w", err)
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

	// Get metric export interval (default 10 seconds for lower overhead)
	exportInterval := cfg.GetOtelMetricExportInterval()
	if exportInterval == 0 {
		exportInterval = 10
	}

	// Get export timeout (default 30 seconds)
	exportTimeout := cfg.GetOtelExportTimeout()
	if exportTimeout == 0 {
		exportTimeout = 30
	}

	// Create periodic reader with optimized non-blocking batch processing
	// PeriodicReader exports metrics in background goroutine without blocking application
	reader := sdkmetric.NewPeriodicReader(
		exporter,
		sdkmetric.WithInterval(time.Duration(exportInterval)*time.Second),
		sdkmetric.WithTimeout(time.Duration(exportTimeout)*time.Second),
	)

	// Create meter provider with async reader
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(reader),
	)

	// Set global meter provider
	otel.SetMeterProvider(mp)

	log.Printf("OpenTelemetry metrics initialized: service=%s, endpoint=%s, interval=%ds",
		cfg.GetOtelServiceName(), cfg.GetJaegerEndpoint(), exportInterval)

	return &MeterProvider{
		provider: mp,
	}, nil
}

// Meter returns a named meter
func (mp *MeterProvider) Meter(name string) metric.Meter {
	return mp.provider.Meter(name)
}

// Shutdown gracefully shuts down the meter provider
// This ensures all metrics are flushed before application exit
func (mp *MeterProvider) Shutdown(ctx context.Context) error {
	if mp.provider == nil {
		return nil
	}

	log.Println("Shutting down OpenTelemetry meter provider...")
	if err := mp.provider.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown meter provider: %w", err)
	}

	log.Println("OpenTelemetry meter provider shut down successfully")
	return nil
}

// GetProvider returns the underlying SDK meter provider
func (mp *MeterProvider) GetProvider() *sdkmetric.MeterProvider {
	return mp.provider
}
