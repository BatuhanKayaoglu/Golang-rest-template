package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"golang-rest-api-template/pkg/apm"

	"github.com/gin-gonic/gin"
)

// ErrorHandler is a middleware that recovers from unexpected panics.
// For business errors, use response.XX() functions with explicit return.
//
// Usage in handlers:
//
//	func MyHandler(c *gin.Context) {
//	    if err != nil {
//	        response.NotFound(c, "Resource not found")
//	        return  // Always return after error response
//	    }
//	    response.OK(c, result)
//	}
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// Unexpected panic - log stack trace
				stack := string(debug.Stack())
				fmt.Printf("[PANIC RECOVERED] %v\n%s\n", r, stack)

				var err error
				switch v := r.(type) {
				case error:
					err = v
				default:
					err = fmt.Errorf("panic: %v", v)
				}

				// Send to APM
				apm.CaptureError(c.Request.Context(), err)

				// Return 500 response
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"error": gin.H{
						"code":    "INTERNAL_ERROR",
						"message": "Internal server error",
					},
				})
			}
		}()

		c.Next()
	}
}
