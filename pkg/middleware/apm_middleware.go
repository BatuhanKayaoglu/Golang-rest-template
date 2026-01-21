package middleware

import (
	"golang-rest-api-template/pkg/config"

	"github.com/gin-gonic/gin"
	"go.elastic.co/apm/module/apmgin/v2"
)

// APM returns a middleware that traces requests with Elastic APM.
// It automatically captures:
// - Transaction name (HTTP method + route)
// - Response status code
// - Request duration (latency)
// - Request/Response headers
// - Errors
func APM() gin.HandlerFunc {
	cfg := config.Get()
	if !cfg.APM.Active {
		// Return a no-op middleware if APM is disabled
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return apmgin.Middleware(nil)
}
