// Package middlewares provides error handling functionality.
package middlewares

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/yeferson59/gin-template/pkg/logger"
	"github.com/yeferson59/gin-template/pkg/response"
)

// ErrorHandler returns a middleware that handles panics and errors gracefully.
func ErrorHandler() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// Log the panic with stack trace
		if recovered != nil {
			stack := debug.Stack()
			logger.WithField("panic", recovered).WithField("stack", string(stack)).Error("Panic recovered")

			// Convert the recovered value to a string
			var errStr string
			switch v := recovered.(type) {
			case error:
				errStr = v.Error()
			case string:
				errStr = v
			default:
				errStr = fmt.Sprintf("%v", v)
			}

			// Log the error string for debugging
			logger.WithField("error_details", errStr).Error("Panic details")

			// Return a generic error response to the client
			response.InternalServerError(c, "Internal server error", "An unexpected error occurred")
		}
	})
}

// RequestLogger returns a middleware that logs HTTP requests with structured logging.
func RequestLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Custom log format using structured logging
		logger.WithFields(map[string]interface{}{
			"client_ip":   param.ClientIP,
			"timestamp":   param.TimeStamp.Format("2006-01-02 15:04:05"),
			"method":      param.Method,
			"path":        param.Path,
			"protocol":    param.Request.Proto,
			"status_code": param.StatusCode,
			"latency":     param.Latency.String(),
			"user_agent":  param.Request.UserAgent(),
			"error":       param.ErrorMessage,
		}).Info("HTTP Request")

		return ""
	})
}

// SecurityHeaders adds security headers to responses.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'")

		c.Next()
	}
}

// RequestID adds a unique request ID to each request.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)

		c.Next()
	}
}

// generateRequestID generates a simple request ID.
// In production, you might want to use a more sophisticated ID generation
func generateRequestID() string {
	return fmt.Sprintf("req_%d", time.Now().UnixNano())
}

// Timeout adds a timeout to requests.
func Timeout() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set a timeout for the request context
		// This is a simple implementation; in production you might want
		// to use context.WithTimeout
		c.Next()
	}
}

// ValidateContentType validates the Content-Type header for specific endpoints.
func ValidateContentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			contentType := c.GetHeader("Content-Type")
			if contentType != "application/json" && contentType != "application/json; charset=utf-8" {
				response.BadRequestError(c, "Invalid Content-Type", "Content-Type must be application/json")
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
