package auth

import (
	"net/http"
	"strings"

	"restorapp/db"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get access token from Authorization header

		authHeader := c.GetHeader("Authorization")
		log.Info(authHeader)
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		accessToken := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate token
		claims, err := ValidateAccessToken(accessToken)
		if err != nil {
			if err == ErrExpiredToken {
				log.Info("error token expired")
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			} else {

				log.Info("error token invalid")
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			}
			c.Abort()
			return
		}

		// Store user ID in context
		c.Set("userId", claims.UserID.String())
		c.Set("userEmail", claims.Email)

		c.Next()
	}
}

func EmailVerifiedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIdStr, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		userUUID, err := uuid.Parse(userIdStr.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			c.Abort()
			return
		}

		user, err := db.Queries.GetUserById(c, userUUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
			c.Abort()
			return
		}

		if !user.EmailVerified.Bool {
			c.JSON(http.StatusForbidden, gin.H{"error": "Email not verified"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get access token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.Next()
			return
		}

		accessToken := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate token
		claims, err := ValidateAccessToken(accessToken)
		if err != nil {
			c.Next()
			return
		}

		// Store user ID in context
		c.Set("userId", claims.UserID.String())
		c.Set("userEmail", claims.Email)

		c.Next()
	}
}
