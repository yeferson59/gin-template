// Package handlers contains HTTP controllers for authentication and other modules.
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/yeferson59/gin-template/internal/auth"
	"github.com/yeferson59/gin-template/internal/models"
	"github.com/yeferson59/gin-template/internal/validators"
	"github.com/yeferson59/gin-template/pkg/logger"
	"github.com/yeferson59/gin-template/pkg/response"
)

// AuthResponse represents the structure for the token response.
type AuthResponse struct {
	Token string          `json:"token"`
	User  *UserSafeResponse `json:"user"`
}

// UserSafeResponse represents user data safe for API responses.
type UserSafeResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// Register handles user registration.
func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req validators.AuthRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.WithField("error", err.Error()).Warn("Invalid JSON data for registration")
			response.BadRequestError(c, "Invalid request data", err.Error())
			return
		}

		// Validate the request data
		if err := validators.ValidateUserRegistration(&req); err != nil {
			logger.WithField("error", err.Error()).Warn("Validation failed for registration")
			response.ValidationError(c, err.Error())
			return
		}

		// Check if the user already exists by username or email
		var existing models.User
		if err := db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existing).Error; err == nil {
			logger.WithFields(map[string]interface{}{
				"username": req.Username,
				"email":    req.Email,
			}).Warn("Attempt to register with existing username or email")
			response.ConflictError(c, "User already exists", "Username or email already exists")
			return
		}

		// Hash the password
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			logger.WithField("error", err.Error()).Error("Failed to hash password")
			response.InternalServerError(c, "Error processing password", "Failed to secure password")
			return
		}

		user := models.User{
			Username: req.Username,
			Email:    req.Email,
			Password: string(hashed),
		}

		if err := db.Create(&user).Error; err != nil {
			logger.WithField("error", err.Error()).Error("Failed to create user in database")
			response.InternalServerError(c, "Could not create user", "Database error occurred")
			return
		}

		logger.WithFields(map[string]interface{}{
			"user_id":  user.ID,
			"username": user.Username,
			"email":    user.Email,
		}).Info("User registered successfully")

		userResponse := &UserSafeResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		}

		response.SuccessResponse(c, http.StatusCreated, "User registered successfully", userResponse)
	}
}

// Login handles user login.
func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req validators.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.WithField("error", err.Error()).Warn("Invalid JSON data for login")
			response.BadRequestError(c, "Invalid request data", err.Error())
			return
		}

		// Validate the request data
		if err := validators.ValidateUserLogin(&req); err != nil {
			logger.WithField("error", err.Error()).Warn("Validation failed for login")
			response.ValidationError(c, err.Error())
			return
		}

		var user models.User
		if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
			logger.WithField("username", req.Username).Warn("Login attempt with non-existent username")
			response.UnauthorizedError(c, "Invalid credentials", "Username or password is incorrect")
			return
		}

		// Verify password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			logger.WithFields(map[string]interface{}{
				"username": req.Username,
				"user_id":  user.ID,
			}).Warn("Login attempt with incorrect password")
			response.UnauthorizedError(c, "Invalid credentials", "Username or password is incorrect")
			return
		}

		// Generate JWT token using the centralized function
		token, err := auth.GenerateJWT(user.ID, user.Email)
		if err != nil {
			logger.WithField("error", err.Error()).Error("Failed to generate JWT token")
			response.InternalServerError(c, "Authentication failed", "Could not generate access token")
			return
		}

		logger.WithFields(map[string]interface{}{
			"user_id":  user.ID,
			"username": user.Username,
		}).Info("User logged in successfully")

		userResponse := &UserSafeResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		}

		authResponse := AuthResponse{
			Token: token,
			User:  userResponse,
		}

		response.SuccessResponse(c, http.StatusOK, "Login successful", authResponse)
	}
}
