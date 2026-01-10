package webserver

import "context"

// Server represents a service that can be started and gracefully shut down
type Server interface {
	Start() error
	Shutdown(ctx context.Context) error
}
