// Package handlers contains HTTP controllers for authentication and other modules.
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/yeferson59/template-gin-api/internal/auth"
	"github.com/yeferson59/template-gin-api/internal/models"
)

// AuthRequest represents the structure for registration and login requests.
type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse represents the structure for the token response.
type AuthResponse struct {
	Token string `json:"token"`
}

// Register handles user registration.
func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AuthRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
			return
		}

		// Check if the user already exists by username or email
		var existing models.User
		if err := db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
			return
		}

		// Hash the password
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing password"})
			return
		}

		user := models.User{
			Username: req.Username,
			Email:    req.Email,
			Password: string(hashed),
		}

		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
	}
}

// Login handles user login.
func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req AuthRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
			return
		}

		var user models.User
		if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		// Verify password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		// Generate JWT token using the centralized function
		token, err := auth.GenerateJWT(user.ID, user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}

		c.JSON(http.StatusOK, AuthResponse{Token: token})
	}
}
