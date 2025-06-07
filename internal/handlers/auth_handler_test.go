package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/yeferson59/template-gin-api/internal/models"
)

// setupTestDB creates an in-memory SQLite database for testing.
func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	_ = db.AutoMigrate(&models.User{})
	return db
}

// setupRouter configures a Gin router for testing.
func setupRouter(db *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/register", Register(db))
	r.POST("/login", Login(db))
	return r
}

func TestRegisterAndLogin(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)
	os.Setenv("JWT_SECRET", "testsecret")

	// Test data
	user := map[string]string{
		"username": "testuser",
		"email":    "testuser@example.com",
		"password": "testpass123",
	}
	body, _ := json.Marshal(user)

	// Registration test
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d, body: %s", w.Code, w.Body.String())
	}

	// Verify that the email was saved correctly after registration
	var userRecord models.User
	if err := db.Where("username = ?", "testuser").First(&userRecord).Error; err != nil {
		t.Fatalf("user not found after registration: %v", err)
	}
	if userRecord.Email != "testuser@example.com" {
		t.Fatalf("user email not saved correctly: got %s, want %s", userRecord.Email, "testuser@example.com")
	}

	// Login test (also sending the email field to match the binding)
	loginData := map[string]string{
		"username": "testuser",
		"email":    "testuser@example.com",
		"password": "testpass123",
	}
	loginBody, _ := json.Marshal(loginData)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse login response: %v", err)
	}
	if resp.Token == "" {
		t.Fatalf("expected a JWT token, got empty string")
	}
}
