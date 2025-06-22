// Package middlewares provides custom middleware for the API.
package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yeferson59/gin-template/internal/auth"
	"github.com/yeferson59/gin-template/internal/models"
	"github.com/yeferson59/gin-template/pkg/logger"
	"github.com/yeferson59/gin-template/pkg/response"
	"gorm.io/gorm"
)

// AuthRequired is a middleware that validates the JWT and checks if the user exists in the database.
func AuthRequired(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.WithField("ip", c.ClientIP()).Warn("Request to protected endpoint without authorization header")
			response.UnauthorizedError(c, "Authorization required", "Authorization header is missing")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			logger.WithField("auth_header", authHeader).Warn("Invalid authorization header format")
			response.UnauthorizedError(c, "Invalid authorization format", "Authorization header must be in format 'Bearer <token>'")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := auth.ValidateJWT(tokenString)
		if err != nil {
			logger.WithField("error", err.Error()).Warn("Invalid or expired JWT token")
			response.UnauthorizedError(c, "Invalid or expired token", err.Error())
			c.Abort()
			return
		}

		// Check if the user exists in the database
		var user models.User
		if err := db.First(&user, claims.UserID).Error; err != nil {
			logger.WithFields(map[string]interface{}{
				"user_id": claims.UserID,
				"error":   err.Error(),
			}).Warn("JWT token refers to non-existent user")
			response.UnauthorizedError(c, "Invalid token", "User associated with token not found")
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", user.ID)
		c.Set("user", user)
		c.Set("email", user.Email)
		c.Set("username", user.Username)

		logger.WithFields(map[string]interface{}{
			"user_id":  user.ID,
			"username": user.Username,
			"endpoint": c.Request.URL.Path,
		}).Debug("User authenticated successfully")

		c.Next()
	}
}

// ProtectedHandler is an example of a JWT-protected endpoint.
func ProtectedHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		email, _ := c.Get("email")
		username, _ := c.Get("username")

		data := gin.H{
			"user_id":  userID,
			"email":    email,
			"username": username,
			"message":  "You have successfully accessed a protected resource",
		}

		response.SuccessResponse(c, http.StatusOK, "Access granted to protected resource", data)
	}
}
