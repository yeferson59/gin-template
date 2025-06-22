// Package database handles the connection and operations with the database.
package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/yeferson59/gin-template/internal/config"
)

// InitDB initializes the database connection using GORM.
// Supports SQLite, PostgreSQL, and MySQL depending on configuration.
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	driver := cfg.Database.Driver
	dsn := cfg.Database.DSN

	// For SQLite, ensure the directory exists
	if strings.ToLower(driver) == "sqlite" {
		if err := ensureDirectoryExists(dsn); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	var db *gorm.DB
	var err error

	switch strings.ToLower(driver) {
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	case "postgres", "postgresql":
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	case "mysql":
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	// Test connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get DB instance: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("could not ping the database: %w", err)
	}

	log.Printf("Connected to the database using %s", driver)
	return db, nil
}

// CloseDB closes the database connection.
func CloseDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Error getting DB instance to close: %v", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		log.Printf("Error closing the database: %v", err)
	}
}

// ensureDirectoryExists creates the directory for the database file if it doesn't exist.
func ensureDirectoryExists(dbPath string) error {
	dir := filepath.Dir(dbPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		log.Printf("Created database directory: %s", dir)
	}
	return nil
}
