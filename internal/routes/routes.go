// Package routes contiene el registro de rutas main de la API.
package routes

import (
	"github.com/yeferson59/template-gin-api/internal/config"
	"github.com/yeferson59/template-gin-api/internal/handlers"
	"github.com/yeferson59/template-gin-api/internal/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterAPIRoutes registra las rutas main de la API.
func RegisterAPIRoutes(router *gin.Engine, db *gorm.DB, _ *config.Config) {
	api := router.Group("/api")
	{
		api.POST("/register", handlers.Register(db))
		api.POST("/login", handlers.Login(db))

		// Endpoint protegido con JWT y verificación real de usuario
		protected := api.Group("/protected")
		protected.Use(middlewares.AuthRequired(db))
		protected.GET("", middlewares.ProtectedHandler())
	}
}
