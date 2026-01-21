package middleware

import (
	"errors"
	"fmt"
	"runtime/debug"

	"golang-rest-api-template/pkg/apm"
	"golang-rest-api-template/pkg/response"

	"github.com/gin-gonic/gin"
)

// Recovery returns a middleware that recovers from any panics and writes a 500 response.
// It also captures the panic in APM for monitoring.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// Get stack trace
				stack := string(debug.Stack())

				// Log the panic
				fmt.Printf("[PANIC RECOVERED] %v\n%s\n", r, stack)

				// Create error from panic
				var err error
				switch v := r.(type) {
				case error:
					err = v
				case string:
					err = errors.New(v)
				default:
					err = fmt.Errorf("panic: %v", v)
				}

				// Capture in APM
				apm.CaptureError(c.Request.Context(), err)

				// Return 500 response
				response.InternalServerError(c, "Internal server error")
				c.Abort()
			}
		}()

		c.Next()
	}
}
