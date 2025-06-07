// Package database handles the connection and operations with the database.
package database

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/yeferson59/template-gin-api/internal/config"
)

// InitDB initializes the database connection using GORM.
// Supports SQLite, PostgreSQL, and MySQL depending on configuration.
func InitDB(_ *config.Config) (*gorm.DB, error) {
	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		driver = "sqlite"
	}
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		// Default value for SQLite
		dsn = "./data/app.db"
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
