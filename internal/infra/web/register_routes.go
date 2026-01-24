package web

import (
	"github.com/gin-gonic/gin"
	"github.com/refortunato/go_app_base/cmd/server/container"
	exampleWeb "github.com/refortunato/go_app_base/internal/example/infra/web"
	healthWeb "github.com/refortunato/go_app_base/internal/health/infra/web"
)

// RegisterRoutes is the main route orchestrator
// It delegates route registration to each module
func RegisterRoutes(c *container.Container) func(*gin.Engine) {
	return func(router *gin.Engine) {
		// Register routes for each module
		healthWeb.RegisterRoutes(router, c.HealthModule.HealthController)
		exampleWeb.RegisterRoutes(router, c.ExampleModule.ExampleController)
	}
}
