package webserver

import (
	"github.com/gin-gonic/gin"
	"github.com/refortunato/go_app_base/internal/infra/dependencies"
)

func Handler(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		dependencies.HealthController.HealthCheck(&GinContextAdapter{ctx: c})
	})
	router.GET("/examples/:id", func(c *gin.Context) {
		dependencies.ExampleController.GetExample(&GinContextAdapter{ctx: c})
	})
}
