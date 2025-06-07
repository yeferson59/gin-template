// Package config handles the global configuration of the application.
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config contains the global configuration for the application.
type Config struct {
	AppName     string
	Port        string
	DatabaseURL string
	JWTSecret   string
}

// Cfg is the loaded global configuration instance.
var Cfg *Config

// LoadConfig loads configuration from environment variables.
func LoadConfig() {
	// Load .env if it exists
	_ = godotenv.Load()

	Cfg = &Config{
		AppName:     getEnv("APP_NAME", "GinAPI"),
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "sqlite.db"),
		JWTSecret:   getEnv("JWT_SECRET", "supersecretkey"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}

// MustLoad loads the configuration and terminates execution if any critical variable is missing.
func MustLoad() {
	LoadConfig()
	if Cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET must be set")
	}
	if Cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL must be set")
	}
}
