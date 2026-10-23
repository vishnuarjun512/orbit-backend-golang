package user

import (
	"errors"
	"fmt"
	"net/http"
	"orbit-backend-golang/internal/config"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid registration details",
		})
		return
	}

	err := h.service.Register(c.Request.Context(), req)

	if err != nil {
		switch {
		case errors.Is(err, ErrEmailAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{
				"error":   true,
				"message": "Email is already registered",
			})

		case errors.Is(err, ErrInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "Invalid registration details",
			})

		default:
			fmt.Printf("Error : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Something went wrong",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"error":   false,
		"message": "Registration successful",
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "Invalid Login Details",
		})
	}

	loginResponse, err := h.service.SignIn(c.Request.Context(), req)

	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidPassword):
			c.JSON(http.StatusConflict, gin.H{
				"error":   true,
				"message": "Invalid Credentials",
			})

		case errors.Is(err, ErrUserNotRegistered):
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "User not found",
			})

		case errors.Is(err, ErrTokenGenerationError):
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   true,
				"message": "Token generation failed",
			})

		default:
			fmt.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "Something went wrong",
			})
		}

		return
	}

	c.SetCookie(
		"orbit_access_token",
		loginResponse.AccessToken,
		config.LoadConfig().JWT_ACCESS_TOKEN_EXPIRY_IN_MINUTES*60, // 15 minutes
		"/",
		"",
		true, // Secure
		true, // HttpOnly
	)

	c.SetCookie(
		"orbit_refresh_token",
		loginResponse.RefreshToken,
		config.LoadConfig().JWT_REFRESH_TOKEN_EXPIRY_IN_MINUTES*60, // 15 minutes
		"/",
		"",
		true, // Secure
		true, // HttpOnly
	)

	c.JSON(http.StatusCreated, gin.H{
		"error":   false,
		"message": "Login successful",
		"user":    loginResponse.User,
	})

}

func (h *Handler) Refresh(c *gin.Context) {

	refreshToken, err := c.Cookie("orbit_refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Refresh Token Extraction Error",
		})
		return
	}

	token := refreshToken

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Refresh Token Empty",
		})
		return
	}

	userID, err := h.service.jwtService.ValidateAccessToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Refresh Token Validation Error",
		})
		return
	}

	accessToken, err := h.service.jwtService.CreateToken(userID, config.LoadConfig().JWT_ACCESS_TOKEN_EXPIRY_IN_MINUTES)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "New Access Token Creation Error",
		})
		return
	}

	c.SetCookie(
		"orbit_access_token",
		accessToken,
		config.LoadConfig().JWT_ACCESS_TOKEN_EXPIRY_IN_MINUTES*60, // 15 minutes
		"/",
		"",
		true, // Secure
		true, // HttpOnly
	)

	c.JSON(http.StatusCreated, gin.H{
		"error":   false,
		"message": "Refreshed Access token successful",
	})

}
