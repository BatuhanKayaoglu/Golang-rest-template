package response

import (
	"net/http"

	"golang-rest-api-template/pkg/apm"
	apperrors "golang-rest-api-template/pkg/errors"

	"github.com/gin-gonic/gin"
)

// APIResponse is the unified response structure for all API endpoints.
// Both success and error responses follow this pattern.
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo contains error details in the response.
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ============================================
// Success Responses
// ============================================

// OK sends a 200 OK response with data.
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
	})
}

// Created sends a 201 Created response with data.
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    data,
	})
}

// NoContent sends a 204 No Content response.
func NoContent(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}

// ============================================
// Error Responses (Idiomatic Go - requires return after call)
// ============================================

// Error sends an error response using AppError.
// Always follow with a return statement.
//
// Usage:
//
//	if err != nil {
//	    response.Error(c, apperrors.NotFound("Book not found"))
//	    return
//	}
func Error(c *gin.Context, err *apperrors.AppError) {
	// Send error to APM (ELK)
	apm.CaptureError(c.Request.Context(), err)

	c.AbortWithStatusJSON(err.StatusCode, APIResponse{
		Success: false,
		Error: &ErrorInfo{
			Code:    err.Code,
			Message: err.Message,
		},
	})
}

// ============================================
// Convenience Error Functions (requires return after call)
// ============================================

// BadRequest sends a 400 Bad Request error.
func BadRequest(c *gin.Context, message string) {
	Error(c, apperrors.BadRequest(message))
}

// Unauthorized sends a 401 Unauthorized error.
func Unauthorized(c *gin.Context, message string) {
	Error(c, apperrors.Unauthorized(message))
}

// Forbidden sends a 403 Forbidden error.
func Forbidden(c *gin.Context, message string) {
	Error(c, apperrors.Forbidden(message))
}

// NotFound sends a 404 Not Found error.
func NotFound(c *gin.Context, message string) {
	Error(c, apperrors.NotFound(message))
}

// Conflict sends a 409 Conflict error.
func Conflict(c *gin.Context, message string) {
	Error(c, apperrors.Conflict(message))
}

// ValidationError sends a 422 Unprocessable Entity error.
func ValidationError(c *gin.Context, message string) {
	Error(c, apperrors.ValidationError(message))
}

// InternalServerError sends a 500 Internal Server Error.
func InternalServerError(c *gin.Context, message string) {
	Error(c, apperrors.InternalError(message, nil))
}

// ============================================
// Business Error Functions (requires return after call)
// ============================================

// BusinessRuleViolation sends a 400 error for business rule violations.
func BusinessRuleViolation(c *gin.Context, message string) {
	Error(c, apperrors.BusinessRuleViolation(message))
}

// InsufficientBalance sends a 400 error for insufficient balance.
func InsufficientBalance(c *gin.Context, message string) {
	Error(c, apperrors.InsufficientBalance(message))
}

// DuplicateEntry sends a 409 error for duplicate entries.
func DuplicateEntry(c *gin.Context, message string) {
	Error(c, apperrors.DuplicateEntry(message))
}

// RateLimitExceeded sends a 429 error for rate limiting.
func RateLimitExceeded(c *gin.Context, message string) {
	Error(c, apperrors.RateLimitExceeded(message))
}

// ServiceUnavailable sends a 503 error for service failures.
func ServiceUnavailable(c *gin.Context, message string) {
	Error(c, apperrors.ServiceUnavailable(message))
}
