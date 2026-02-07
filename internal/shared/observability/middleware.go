package observability

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// TracingMiddleware returns a Gin middleware that instruments HTTP requests with OpenTelemetry
// This middleware automatically:
// - Creates a span for each HTTP request
// - Propagates trace context (W3C Trace Context headers)
// - Captures HTTP method, path, status code, and errors
// - Adds span attributes for request metadata
func TracingMiddleware(serviceName string) gin.HandlerFunc {
	return otelgin.Middleware(serviceName)
}
