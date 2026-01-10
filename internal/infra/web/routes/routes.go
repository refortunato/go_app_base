package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/refortunato/go_app_base/cmd/server/container"
	"github.com/refortunato/go_app_base/internal/shared/web/context"
)

// RegisterRoutes registers all application routes
// It receives the dependency container to access controllers
func RegisterRoutes(c *container.Container) func(*gin.Engine) {
	return func(router *gin.Engine) {
		router.GET("/health", func(ctx *gin.Context) {
			c.HealthController.HealthCheck(context.NewGinContextAdapter(ctx))
		})
		router.GET("/examples/:id", func(ctx *gin.Context) {
			c.ExampleController.GetExample(context.NewGinContextAdapter(ctx))
		})
	}
}
