package webserver

import (
	"context"

	"github.com/gin-gonic/gin"
)

// Adaptador para gin.Context que implementa a interface WebContext
type GinContextAdapter struct {
	ctx *gin.Context
}

func (g *GinContextAdapter) JSON(code int, obj any) {
	g.ctx.JSON(code, obj)
}

func (g *GinContextAdapter) BindJSON(obj any) error {
	return g.ctx.BindJSON(obj)
}

func (g *GinContextAdapter) Param(key string) string {
	return g.ctx.Param(key)
}

func (g *GinContextAdapter) Query(key string) string {
	return g.ctx.Query(key)
}

func (g *GinContextAdapter) GetHeader(key string) string {
	return g.ctx.GetHeader(key)
}

func (g *GinContextAdapter) SetHeader(key, value string) {
	g.ctx.Header(key, value)
}

func (g *GinContextAdapter) GetContext() context.Context {
	return g.ctx.Request.Context()
}
