// Package routes contiene el registro de rutas main de la API.
package routes

import (
	"github.com/yeferson59/gin-template/internal/config"
	"github.com/yeferson59/gin-template/internal/handlers"
	"github.com/yeferson59/gin-template/internal/middlewares"
	"github.com/yeferson59/gin-template/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterAPIRoutes registra las rutas main de la API.
func RegisterAPIRoutes(router *gin.Engine, db *gorm.DB, cfg *config.Config) {
	// Health check endpoints (no rate limiting for monitoring)
	health := router.Group("/health")
	{
		health.GET("/", handlers.HealthCheck(db))
		health.GET("/live", handlers.LivenessCheck())
		health.GET("/ready", handlers.ReadinessCheck(db))
	}

	// API routes with rate limiting
	api := router.Group("/api")
	api.Use(middlewares.RateLimit())
	api.Use(middlewares.ValidateContentType())
	{
		// Authentication endpoints with stricter rate limiting
		auth := api.Group("/auth")
		auth.Use(middlewares.AuthRateLimit())
		{
			auth.POST("/register", handlers.Register(db))
			auth.POST("/login", handlers.Login(db))
		}

		// Legacy endpoints (for backward compatibility)
		api.POST("/register", middlewares.AuthRateLimit(), handlers.Register(db))
		api.POST("/login", middlewares.AuthRateLimit(), handlers.Login(db))

		// Protected endpoints
		protected := api.Group("/protected")
		protected.Use(middlewares.AuthRequired(db))
		{
			protected.GET("/", middlewares.ProtectedHandler())
			protected.GET("/profile", getUserProfile())
		}

		// User endpoints
		users := api.Group("/users")
		users.Use(middlewares.AuthRequired(db))
		{
			users.GET("/me", getUserProfile())
			// Add more user endpoints as needed
		}
	}
}

// getUserProfile returns the current user's profile
func getUserProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		email, _ := c.Get("email")
		username, _ := c.Get("username")

		profile := gin.H{
			"id":       userID,
			"username": username,
			"email":    email,
		}

		response.SuccessResponse(c, 200, "User profile retrieved successfully", profile)
	}
}
