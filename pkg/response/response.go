// Package response provides standardized API response structures.
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse defines the standard structure for all API responses.
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError defines the structure for error responses.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse sends a successful response.
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse sends an error response.
func ErrorResponse(c *gin.Context, statusCode int, code, message, details string) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// BadRequestError sends a 400 Bad Request error.
func BadRequestError(c *gin.Context, message, details string) {
	ErrorResponse(c, http.StatusBadRequest, "BAD_REQUEST", message, details)
}

// UnauthorizedError sends a 401 Unauthorized error.
func UnauthorizedError(c *gin.Context, message, details string) {
	ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", message, details)
}

// ForbiddenError sends a 403 Forbidden error.
func ForbiddenError(c *gin.Context, message, details string) {
	ErrorResponse(c, http.StatusForbidden, "FORBIDDEN", message, details)
}

// NotFoundError sends a 404 Not Found error.
func NotFoundError(c *gin.Context, message, details string) {
	ErrorResponse(c, http.StatusNotFound, "NOT_FOUND", message, details)
}

// ConflictError sends a 409 Conflict error.
func ConflictError(c *gin.Context, message, details string) {
	ErrorResponse(c, http.StatusConflict, "CONFLICT", message, details)
}

// InternalServerError sends a 500 Internal Server Error.
func InternalServerError(c *gin.Context, message, details string) {
	ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", message, details)
}

// ValidationError sends a validation error response.
func ValidationError(c *gin.Context, details string) {
	ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", details)
}
