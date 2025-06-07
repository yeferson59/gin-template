// Package main es el entrypoint de la aplicación.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/yeferson59/template-gin-api/internal/config"
	"github.com/yeferson59/template-gin-api/internal/database"
	"github.com/yeferson59/template-gin-api/internal/middlewares"
	"github.com/yeferson59/template-gin-api/internal/models"
	"github.com/yeferson59/template-gin-api/internal/routes"
)

func main() {
	// Cargar variables de entorno desde .env si existe
	_ = godotenv.Load()

	// Cargar configuración de la app
	config.LoadConfig()
	cfg := config.Cfg

	// Inicializar base de datos (SQLite por defecto, pero configurable)
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Error al conectar la base de datos: %v", err)
		return
	}
	// Migración automática del modelo User
	if err := db.AutoMigrate(&models.User{}); err != nil {
		database.CloseDB(db)
		log.Fatalf("Error en migración de User: %v", err)
		return
	}

	// Inicializar Gin
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	// Middlewares globales
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(middlewares.CORS())

	// Registrar rutas
	routes.RegisterAPIRoutes(router, db, cfg)

	// Iniciar servidor
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Servidor escuchando en %s", addr)
	if err := router.Run(addr); err != nil {
		database.CloseDB(db)
		log.Fatalf("No se pudo iniciar el servidor: %v", err)
		return
	}

	database.CloseDB(db)
}
