// Package main es el entrypoint de la aplicación.
package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/yeferson59/gin-template/internal/config"
	"github.com/yeferson59/gin-template/internal/database"
	"github.com/yeferson59/gin-template/internal/middlewares"
	"github.com/yeferson59/gin-template/internal/models"
	"github.com/yeferson59/gin-template/internal/routes"
	"github.com/yeferson59/gin-template/pkg/logger"
)

func main() {
	// Parse command line flags
	healthCheck := flag.Bool("health-check", false, "Perform health check and exit")
	version := flag.Bool("version", false, "Show version and exit")
	flag.Parse()

	// Handle version flag
	if *version {
		fmt.Println("Gin Template API v1.0.0")
		os.Exit(0)
	}

	// Handle health check flag (for Docker HEALTHCHECK)
	if *healthCheck {
		performHealthCheck()
		return
	}

	// Load environment variables
	_ = godotenv.Load()

	// Initialize logger first
	logger.Init()

	// Load application configuration
	config.LoadConfig()
	cfg := config.Cfg

	logger.WithFields(map[string]interface{}{
		"app_name":    cfg.Server.AppName,
		"environment": cfg.Server.Environment,
		"port":        cfg.Server.Port,
		"db_driver":   cfg.Database.Driver,
	}).Info("Starting application with configuration")

	// Initialize database
	db, err := database.InitDB(cfg)
	if err != nil {
		logger.WithField("error", err.Error()).Fatal("Failed to connect to database")
		return
	}
	defer database.CloseDB(db)

	// Configure database connection pool
	sqlDB, err := db.DB()
	if err != nil {
		logger.WithField("error", err.Error()).Fatal("Failed to get database instance")
		return
	}
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	// Auto-migrate models
	if err := db.AutoMigrate(&models.User{}); err != nil {
		logger.WithField("error", err.Error()).Fatal("Failed to migrate User model")
		return
	}
	logger.Info("Database migrations completed successfully")

	// Set Gin mode based on environment
	if config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else if config.IsTest() {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Initialize Gin router
	router := gin.New()

	// Global middlewares
	router.Use(middlewares.ErrorHandler())
	router.Use(middlewares.RequestLogger())
	router.Use(middlewares.SecurityHeaders())
	router.Use(middlewares.RequestID())
	router.Use(middlewares.CORS())

	// Register routes
	routes.RegisterAPIRoutes(router, db, cfg)

	// Create HTTP server with timeouts
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		MaxHeaderBytes: int(cfg.Server.MaxBodySize),
	}

	// Start server in a goroutine
	go func() {
		logger.WithFields(map[string]interface{}{
			"addr":         server.Addr,
			"environment":  cfg.Server.Environment,
			"read_timeout": cfg.Server.ReadTimeout,
			"write_timeout": cfg.Server.WriteTimeout,
		}).Info("Starting HTTP server")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithField("error", err.Error()).Fatal("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.WithField("error", err.Error()).Error("Server forced to shutdown")
	} else {
		logger.Info("Server shutdown completed gracefully")
	}
}

// performHealthCheck performs a health check for Docker HEALTHCHECK
func performHealthCheck() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Get(fmt.Sprintf("http://localhost:%s/health/live", port))
	if err != nil {
		fmt.Printf("Health check failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Health check failed with status: %d\n", resp.StatusCode)
		os.Exit(1)
	}

	fmt.Println("Health check passed")
	os.Exit(0)
}
