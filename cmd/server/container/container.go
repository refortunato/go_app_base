package container

import (
	"database/sql"

	"github.com/refortunato/go_app_base/configs"
	exampleInfra "github.com/refortunato/go_app_base/internal/example/infra"
	healthInfra "github.com/refortunato/go_app_base/internal/health/infra"
	"github.com/refortunato/go_app_base/internal/shared/logger"
)

// Container holds all application dependencies
// This is the Composition Root of the application
type Container struct {
	// Modules
	ExampleModule *exampleInfra.ExampleModule
	HealthModule  *healthInfra.HealthModule

	// Logger (shared utility)
	Logger logger.Logger
}

// New creates and wires all application dependencies
// This is the only place where dependencies are composed
func New(db *sql.DB, cfg *configs.Conf) (*Container, error) {
	// Logger
	log := logger.NewSlogLogger(cfg.ImageName, cfg.ImageVersion)
	logger.SetGlobalLogger(log)
	logger.Info("Logger initialized successfully")

	// Initialize modules (each module wires its own dependencies)
	exampleModule := exampleInfra.NewExampleModule(db)
	healthModule := healthInfra.NewHealthModule(db)

	return &Container{
		ExampleModule: exampleModule,
		HealthModule:  healthModule,
		Logger:        log,
	}, nil
}
