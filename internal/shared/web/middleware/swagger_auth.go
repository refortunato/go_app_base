package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// SwaggerBasicAuth middleware protects Swagger documentation with Basic Authentication
// In production environment, it requires valid credentials
// In development/staging, it can be optionally disabled via environment variables
func SwaggerBasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		env := os.Getenv("SERVER_APP_ENVIRONMENT")
		swaggerEnabled := os.Getenv("SERVER_APP_SWAGGER_ENABLED")

		// Se swagger está explicitamente desabilitado, bloqueia acesso
		if swaggerEnabled == "false" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Swagger documentation is disabled",
			})
			return
		}

		// Em produção ou staging, sempre exige autenticação
		if env == "production" || env == "staging" {
			username := os.Getenv("SERVER_APP_SWAGGER_USER")
			password := os.Getenv("SERVER_APP_SWAGGER_PASS")

			// Valida se credenciais estão configuradas
			if username == "" || password == "" {
				c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
					"error": "Swagger authentication not configured",
				})
				return
			}

			// Valida Basic Auth
			user, pass, hasAuth := c.Request.BasicAuth()

			if !hasAuth || user != username || pass != password {
				// Envia header para o navegador pedir credenciais
				c.Header("WWW-Authenticate", `Basic realm="Swagger Documentation - Restricted Access"`)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Authentication required",
				})
				return
			}
		}

		// Em desenvolvimento, permite acesso livre (ou com auth opcional)
		c.Next()
	}
}
