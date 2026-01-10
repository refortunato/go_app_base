package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GinServer wraps http.Server for graceful shutdown
type GinServer struct {
	httpServer *http.Server
}

// Shutdown gracefully shuts down the server
func (s *GinServer) Shutdown(ctx context.Context) error {
	fmt.Println("Shutting down HTTP server...")
	return s.httpServer.Shutdown(ctx)
}

// Start starts the server and blocks until it's stopped
func (s *GinServer) Start() error {
	fmt.Printf("Starting HTTP server on %s\n", s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// NewGinServer creates a new GinServer with the provided router and port
func NewGinServer(router *gin.Engine, port string) *GinServer {
	if port == "" {
		port = "8080"
	}

	httpServer := &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	return &GinServer{
		httpServer: httpServer,
	}
}
