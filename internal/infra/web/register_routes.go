package web

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/refortunato/go_app_base/cmd/server/container"
	exampleWeb "github.com/refortunato/go_app_base/internal/example/infra/web"
	healthWeb "github.com/refortunato/go_app_base/internal/health/infra/web"
	"github.com/refortunato/go_app_base/internal/simple_module"
)

// RegisterRoutes is the main route orchestrator
// It delegates route registration to each module
func RegisterRoutes(c *container.Container) func(*gin.Engine) {
	return func(router *gin.Engine) {
		// Swagger documentation
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		// Register routes for each module
		healthWeb.RegisterRoutes(router, c.HealthModule)
		exampleWeb.RegisterRoutes(router, c.ExampleModule)
		simple_module.RegisterRoutes(router, c.SimpleModule)
	}
}
