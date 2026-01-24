package web

import (
	"github.com/gin-gonic/gin"
	"github.com/refortunato/go_app_base/internal/example/infra"
	"github.com/refortunato/go_app_base/internal/shared/web/context"
)

// RegisterRoutes registers all routes for the example module
func RegisterRoutes(router *gin.Engine, module *infra.ExampleModule) {
	router.GET("/examples/:id", func(ctx *gin.Context) {
		module.ExampleController.GetExample(context.NewGinContextAdapter(ctx))
	})
}
