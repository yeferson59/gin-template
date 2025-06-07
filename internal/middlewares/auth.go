package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yeferson59/template-gin-api/internal/auth"
	"github.com/yeferson59/template-gin-api/internal/models"
	"gorm.io/gorm"
)

// AuthRequired is a middleware that validates the JWT and checks if the user exists in the database.
func AuthRequired(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := auth.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Check if the user exists in the database
		var user models.User
		if err := db.First(&user, claims.UserID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		c.Set("user_id", user.ID)
		c.Set("email", user.Email)
		c.Next()
	}
}

// ProtectedHandler is an example of a JWT-protected endpoint.
func ProtectedHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		email, _ := c.Get("email")
		c.JSON(http.StatusOK, gin.H{
			"message": "Access granted to protected resource",
			"user_id": userID,
			"email":   email,
		})
	}
}
