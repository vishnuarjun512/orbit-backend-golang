package middleware

import (
	"net/http"
	"orbit-backend-golang/internal/security"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtService *security.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Get JWT from HttpOnly cookie
		token, err := c.Cookie("orbit_access_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		// Validate JWT and extract user ID
		userID, err := jwtService.ValidateAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

		// Store authenticated user ID in Gin context
		c.Set("userID", userID)

		// Continue to the actual handler
		c.Next()
	}
}
