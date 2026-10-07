package middleware

import (
	"net/http"
	"strings"

	"orbit-backend-golang/internal/security"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtService *security.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := ""
		if authorization := strings.TrimSpace(c.GetHeader("Authorization")); authorization != "" {
			parts := strings.Fields(authorization)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error":   true,
					"message": "Unauthorized",
				})
				c.Abort()
				return
			}
			token = parts[1]
		} else {
			cookieToken, err := c.Cookie("orbit_access_token")
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error":   true,
					"message": "Unauthorized",
				})
				c.Abort()
				return
			}
			token = cookieToken
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   true,
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}

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
