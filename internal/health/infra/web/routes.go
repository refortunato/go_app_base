package web

import (
	"github.com/gin-gonic/gin"
	"github.com/refortunato/go_app_base/internal/health/infra"
	"github.com/refortunato/go_app_base/internal/shared/web/context"
)

// RegisterRoutes registers all routes for the health module
func RegisterRoutes(router *gin.Engine, module *infra.HealthModule) {
	router.GET("/health", func(ctx *gin.Context) {
		module.HealthController.HealthCheck(context.NewGinContextAdapter(ctx))
	})
}
