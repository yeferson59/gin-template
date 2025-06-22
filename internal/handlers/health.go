// Package handlers contains HTTP controllers for health checks and other modules.
package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/yeferson59/gin-template/pkg/logger"
	"github.com/yeferson59/gin-template/pkg/response"
)

// HealthCheckResponse represents the structure for health check responses.
type HealthCheckResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version,omitempty"`
	Services  map[string]string `json:"services"`
}

// HealthCheck provides a comprehensive health check endpoint.
func HealthCheck(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		healthResp := HealthCheckResponse{
			Status:    "ok",
			Timestamp: time.Now(),
			Version:   "1.0.0", // You can make this dynamic
			Services:  make(map[string]string),
		}

		// Check database connection
		if db != nil {
			sqlDB, err := db.DB()
			if err != nil {
				logger.WithField("error", err.Error()).Error("Failed to get database instance")
				healthResp.Status = "error"
				healthResp.Services["database"] = "error"
			} else if err := sqlDB.Ping(); err != nil {
				logger.WithField("error", err.Error()).Error("Database ping failed")
				healthResp.Status = "degraded"
				healthResp.Services["database"] = "error"
			} else {
				healthResp.Services["database"] = "ok"
			}
		} else {
			healthResp.Services["database"] = "not_configured"
		}

		// Add more service checks here as needed
		// For example: Redis, external APIs, etc.

		statusCode := http.StatusOK
		if healthResp.Status == "error" {
			statusCode = http.StatusServiceUnavailable
		} else if healthResp.Status == "degraded" {
			statusCode = http.StatusPartialContent
		}

		response.SuccessResponse(c, statusCode, "Health check completed", healthResp)
	}
}

// ReadinessCheck provides a readiness check endpoint for Kubernetes.
func ReadinessCheck(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if all critical services are ready
		if db != nil {
			sqlDB, err := db.DB()
			if err != nil || sqlDB.Ping() != nil {
				response.ErrorResponse(c, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "Service not ready", "Database connection failed")
				return
			}
		}

		response.SuccessResponse(c, http.StatusOK, "Service is ready", gin.H{
			"status":    "ready",
			"timestamp": time.Now(),
		})
	}
}

// LivenessCheck provides a liveness check endpoint for Kubernetes.
func LivenessCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simple liveness check - if this endpoint responds, the service is alive
		response.SuccessResponse(c, http.StatusOK, "Service is alive", gin.H{
			"status":    "alive",
			"timestamp": time.Now(),
		})
	}
}
