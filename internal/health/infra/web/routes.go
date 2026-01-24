package web

import (
	"github.com/gin-gonic/gin"
	"github.com/refortunato/go_app_base/internal/health/infra/web/controllers"
	"github.com/refortunato/go_app_base/internal/shared/web/context"
)

// RegisterRoutes registers all routes for the health module
func RegisterRoutes(router *gin.Engine, controller *controllers.HealthController) {
	router.GET("/health", func(ctx *gin.Context) {
		controller.HealthCheck(context.NewGinContextAdapter(ctx))
	})
}
