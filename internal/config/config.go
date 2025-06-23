// Package config handles the global configuration of the application.
package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config contains the global configuration for the application.
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	JWT      JWTConfig      `json:"jwt"`
	Logging  LoggingConfig  `json:"logging"`
	Security SecurityConfig `json:"security"`
}

// ServerConfig contains server-related configuration.
type ServerConfig struct {
	AppName      string        `json:"app_name"`
	Port         string        `json:"port"`
	Environment  string        `json:"environment"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	MaxBodySize  int64         `json:"max_body_size"`
}

// DatabaseConfig contains database-related configuration.
type DatabaseConfig struct {
	Driver          string        `json:"driver"`
	DSN             string        `json:"dsn"`
	MaxOpenConns    int           `json:"max_open_conns"`
	MaxIdleConns    int           `json:"max_idle_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
}

// JWTConfig contains JWT-related configuration.
type JWTConfig struct {
	Secret         string        `json:"secret"`
	ExpirationTime time.Duration `json:"expiration_time"`
	RefreshTime    time.Duration `json:"refresh_time"`
	Issuer         string        `json:"issuer"`
}

// LoggingConfig contains logging-related configuration.
type LoggingConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
}

// SecurityConfig contains security-related configuration.
type SecurityConfig struct {
	RateLimitRPS   float64 `json:"rate_limit_rps"`
	RateLimitBurst int     `json:"rate_limit_burst"`
	AuthRateLimit  int     `json:"auth_rate_limit"`
	CORSEnabled    bool    `json:"cors_enabled"`
	CORSOrigins    string  `json:"cors_origins"`
}

// Cfg is the loaded global configuration instance.
var Cfg *Config

// LoadConfig loads configuration from environment variables.
func LoadConfig() {
	// Load .env if it exists
	_ = godotenv.Load()

	Cfg = &Config{
		Server: ServerConfig{
			AppName:      getEnv("APP_NAME", "GinAPI"),
			Port:         getEnv("PORT", "8080"),
			Environment:  getEnv("APP_ENV", "development"),
			ReadTimeout:  getDurationEnv("READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getDurationEnv("WRITE_TIMEOUT", 10*time.Second),
			MaxBodySize:  getInt64Env("MAX_BODY_SIZE", 32<<20), // 32MB
		},
		Database: DatabaseConfig{
			Driver:          getEnv("DB_DRIVER", "sqlite"),
			DSN:             getEnv("DB_DSN", "./data/app.db"),
			MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", time.Hour),
		},
		JWT: JWTConfig{
			Secret:         getEnv("JWT_SECRET", "supersecretkey"),
			ExpirationTime: getDurationEnv("JWT_EXP_MINUTES", 60*time.Minute),
			RefreshTime:    getDurationEnv("JWT_REFRESH_MINUTES", 24*time.Hour),
			Issuer:         getEnv("JWT_ISSUER", "gin-api"),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "text"),
		},
		Security: SecurityConfig{
			RateLimitRPS:   getFloat64Env("RATE_LIMIT_RPS", 10.0),
			RateLimitBurst: getIntEnv("RATE_LIMIT_BURST", 20),
			AuthRateLimit:  getIntEnv("AUTH_RATE_LIMIT", 5),
			CORSEnabled:    getBoolEnv("CORS_ENABLED", true),
			CORSOrigins:    getEnv("CORS_ORIGINS", "*"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}

func getIntEnv(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

func getInt64Env(key string, fallback int64) int64 {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intVal
		}
	}
	return fallback
}

func getFloat64Env(key string, fallback float64) float64 {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal
		}
	}
	return fallback
}

func getBoolEnv(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return fallback
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		if durationVal, err := time.ParseDuration(value); err == nil {
			return durationVal
		}
	}
	return fallback
}

// MustLoad loads the configuration and terminates execution if any critical variable is missing.
func MustLoad() {
	LoadConfig()
	if Cfg.JWT.Secret == "" || Cfg.JWT.Secret == "supersecretkey" {
		log.Fatal("JWT_SECRET must be set and not use default value")
	}
	if Cfg.Database.DSN == "" {
		log.Fatal("DB_DSN must be set")
	}
}

// IsDevelopment returns true if the application is running in development mode.
func IsDevelopment() bool {
	return Cfg.Server.Environment == "development"
}

// IsProduction returns true if the application is running in production mode.
func IsProduction() bool {
	return Cfg.Server.Environment == "production"
}

// IsTest returns true if the application is running in test mode.
func IsTest() bool {
	return Cfg.Server.Environment == "test"
}
