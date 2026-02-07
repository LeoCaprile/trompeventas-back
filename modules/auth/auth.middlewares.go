package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get access token from Authorization header
		authHeader := c.GetHeader("Authorization")
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
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			} else {
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
