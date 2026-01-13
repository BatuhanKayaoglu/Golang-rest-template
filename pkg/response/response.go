package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Success(c *gin.Context, status int, data interface{}) {
	c.JSON(status, APIResponse{
		Success: true,
		Data:    data,
	})
}

func Error(c *gin.Context, status int, code, message string) {
	c.JSON(status, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	})
}

func OK(c *gin.Context, data interface{}) {
	Success(c, http.StatusOK, data)
}

func Created(c *gin.Context, data interface{}) {
	Success(c, http.StatusCreated, data)
}

func NoContent(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}

func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, "BAD_REQUEST", message)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, "UNAUTHORIZED", message)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, "NOT_FOUND", message)
}

func InternalServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", message)
}
