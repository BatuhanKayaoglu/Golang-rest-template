package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError represents a custom application error following Go idioms.
// Use this for business logic errors that should map to specific HTTP responses.
type AppError struct {
	Code       string `json:"code"`    // Error code (e.g., "NOT_FOUND", "BAD_REQUEST")
	Message    string `json:"message"` // User-facing error message
	StatusCode int    `json:"-"`       // HTTP status code
	Err        error  `json:"-"`       // Original error (optional, for logging/debugging)
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error for errors.Is/As support.
func (e *AppError) Unwrap() error {
	return e.Err
}

// Is checks if target error matches this AppError by code.
func (e *AppError) Is(target error) bool {
	var appErr *AppError
	if errors.As(target, &appErr) {
		return e.Code == appErr.Code
	}
	return false
}

// Wrap wraps an existing error with AppError context.
// Example: return errors.BadRequest("invalid input").Wrap(err)
func (e *AppError) Wrap(err error) *AppError {
	return &AppError{
		Code:       e.Code,
		Message:    e.Message,
		StatusCode: e.StatusCode,
		Err:        err,
	}
}

// WithMessage returns a copy of the error with a custom message.
func (e *AppError) WithMessage(message string) *AppError {
	return &AppError{
		Code:       e.Code,
		Message:    message,
		StatusCode: e.StatusCode,
		Err:        e.Err,
	}
}

// ============================================
// Error Constructors (Idiomatic Go Style)
// ============================================

// New creates a custom AppError.
func New(statusCode int, code, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

// BadRequest returns a 400 Bad Request error.
func BadRequest(message string) *AppError {
	return &AppError{
		Code:       "BAD_REQUEST",
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

// Unauthorized returns a 401 Unauthorized error.
func Unauthorized(message string) *AppError {
	return &AppError{
		Code:       "UNAUTHORIZED",
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

// Forbidden returns a 403 Forbidden error.
func Forbidden(message string) *AppError {
	return &AppError{
		Code:       "FORBIDDEN",
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
}

// NotFound returns a 404 Not Found error.
func NotFound(message string) *AppError {
	return &AppError{
		Code:       "NOT_FOUND",
		Message:    message,
		StatusCode: http.StatusNotFound,
	}
}

// Conflict returns a 409 Conflict error.
func Conflict(message string) *AppError {
	return &AppError{
		Code:       "CONFLICT",
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

// ValidationError returns a 422 Unprocessable Entity error.
func ValidationError(message string) *AppError {
	return &AppError{
		Code:       "VALIDATION_ERROR",
		Message:    message,
		StatusCode: http.StatusUnprocessableEntity,
	}
}

// InternalError returns a 500 Internal Server Error.
func InternalError(message string, err error) *AppError {
	return &AppError{
		Code:       "INTERNAL_ERROR",
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}

// ============================================
// Business Domain Errors
// ============================================

// BusinessRuleViolation returns a 400 error for business rule violations.
func BusinessRuleViolation(message string) *AppError {
	return &AppError{
		Code:       "BUSINESS_RULE_VIOLATION",
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

// InsufficientBalance returns a 400 error for insufficient funds.
func InsufficientBalance(message string) *AppError {
	return &AppError{
		Code:       "INSUFFICIENT_BALANCE",
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

// DuplicateEntry returns a 409 error for duplicate resources.
func DuplicateEntry(message string) *AppError {
	return &AppError{
		Code:       "DUPLICATE_ENTRY",
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

// ExpiredToken returns a 401 error for expired tokens.
func ExpiredToken(message string) *AppError {
	return &AppError{
		Code:       "TOKEN_EXPIRED",
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

// RateLimitExceeded returns a 429 error for rate limiting.
func RateLimitExceeded(message string) *AppError {
	return &AppError{
		Code:       "RATE_LIMIT_EXCEEDED",
		Message:    message,
		StatusCode: http.StatusTooManyRequests,
	}
}

// ServiceUnavailable returns a 503 error for external service failures.
func ServiceUnavailable(message string) *AppError {
	return &AppError{
		Code:       "SERVICE_UNAVAILABLE",
		Message:    message,
		StatusCode: http.StatusServiceUnavailable,
	}
}

// ============================================
// Helper Functions
// ============================================

// IsAppError checks if an error is an AppError.
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// AsAppError converts an error to AppError if possible.
func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

// GetStatusCode returns the HTTP status code for an error.
// Returns 500 for non-AppError errors.
func GetStatusCode(err error) int {
	if appErr, ok := AsAppError(err); ok {
		return appErr.StatusCode
	}
	return http.StatusInternalServerError
}
