package container

import (
	"context"
	"database/sql"

	"github.com/refortunato/go_app_base/configs"
	exampleInfra "github.com/refortunato/go_app_base/internal/example/infra"
	healthInfra "github.com/refortunato/go_app_base/internal/health/infra"
	"github.com/refortunato/go_app_base/internal/shared/logger"
	"github.com/refortunato/go_app_base/internal/shared/observability"
	"github.com/refortunato/go_app_base/internal/simple_module"
)

// Container holds all application dependencies
// This is the Composition Root of the application
type Container struct {
	// Modules
	ExampleModule *exampleInfra.ExampleModule
	HealthModule  *healthInfra.HealthModule
	SimpleModule  *simple_module.SimpleModule

	// Shared infrastructure
	Logger         logger.Logger
	TracerProvider *observability.TracerProvider
	MeterProvider  *observability.MeterProvider
}

// New creates and wires all application dependencies
// This is the only place where dependencies are composed
func New(db *sql.DB, cfg *configs.Conf, tracerProvider *observability.TracerProvider, meterProvider *observability.MeterProvider) (*Container, error) {
	// Logger
	log := logger.NewSlogLogger(cfg.ImageName, cfg.ImageVersion)
	logger.SetGlobalLogger(log)

	// Use context.Background() for initialization logs (no HTTP request context)
	ctx := context.Background()
	logger.Info(ctx, "Logger initialized successfully")

	// Database tracing is handled at repository level via observability.TraceQuery/TraceExec helpers
	// See internal/shared/observability/db_helpers.go for implementation
	if cfg.OtelEnabled {
		logger.Info(ctx, "Database tracing enabled (via repository helpers)")
	}

	// Initialize modules (each module wires its own dependencies)
	exampleModule := exampleInfra.NewExampleModule(db)
	healthModule := healthInfra.NewHealthModule(db)
	simpleModule := simple_module.NewSimpleModule(db)

	return &Container{
		ExampleModule:  exampleModule,
		HealthModule:   healthModule,
		SimpleModule:   simpleModule,
		Logger:         log,
		TracerProvider: tracerProvider,
		MeterProvider:  meterProvider,
	}, nil
}
