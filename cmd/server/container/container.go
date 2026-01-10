package container

import (
"database/sql"

"github.com/refortunato/go_app_base/internal/core/application/usecases"
"github.com/refortunato/go_app_base/internal/infra/config"
infraRepositories "github.com/refortunato/go_app_base/internal/infra/repositories"
"github.com/refortunato/go_app_base/internal/infra/web/controllers"
"github.com/refortunato/go_app_base/internal/shared/logger"
)

// Container holds all application dependencies
// This is the Composition Root of the application
type Container struct {
	// Controllers (delivery layer)
	ExampleController *controllers.ExampleController
	HealthController  *controllers.HealthController

	// Use Cases (application layer)
	GetExampleUseCase  *usecases.GetExampleUseCase
	HealthCheckUseCase *usecases.HealthCheckUseCase

	// Logger (shared utility)
	Logger logger.Logger
}

// New creates and wires all application dependencies
// This is the only place where dependencies are composed
func New(db *sql.DB, cfg *config.Conf) (*Container, error) {
	// Logger
	log := logger.NewSlogLogger(cfg.ImageName, cfg.ImageVersion)
	logger.SetGlobalLogger(log)
	logger.Info("Logger initialized successfully")

	// Repositories
	exampleRepository := infraRepositories.NewExampleMySQLRepository(db)
	healthRepository := infraRepositories.NewHealthMySQLRepository(db)

	// Use Cases
	getExampleUseCase := usecases.NewGetExampleUseCase(exampleRepository)
	healthCheckUseCase := usecases.NewHealthCheckUseCase(healthRepository)

	// Controllers
	exampleController := controllers.NewExampleController(*getExampleUseCase)
	healthController := controllers.NewHealthController(*healthCheckUseCase)

	return &Container{
		ExampleController:  exampleController,
		HealthController:   healthController,
		GetExampleUseCase:  getExampleUseCase,
		HealthCheckUseCase: healthCheckUseCase,
		Logger:             log,
	}, nil
}
