package web

import (
	"github.com/gin-gonic/gin"
	"github.com/refortunato/go_app_base/internal/example/infra/web/controllers"
	"github.com/refortunato/go_app_base/internal/shared/web/context"
)

// RegisterRoutes registers all routes for the example module
func RegisterRoutes(router *gin.Engine, controller *controllers.ExampleController) {
	router.GET("/examples/:id", func(ctx *gin.Context) {
		controller.GetExample(context.NewGinContextAdapter(ctx))
	})
}
