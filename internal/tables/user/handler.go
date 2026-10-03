package user

import (
	"errors"
	"net/http"

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

	// Parse and validate incoming JSON.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid registration details",
		})
		return
	}

	// Execute registration business logic.
	user, err := h.service.Register(c.Request.Context(), req)

	if err != nil {
		switch {
		case errors.Is(err, ErrEmailAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{
				"status":  "error",
				"message": "Email is already registered",
			})

		case errors.Is(err, ErrUsernameAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{
				"status":  "error",
				"message": "Username is already taken",
			})

		case errors.Is(err, ErrInvalidInput):
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid registration details",
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Something went wrong",
			})
		}

		return
	}

	// Build a public response without exposing sensitive fields.
	response := RegisterResponse{
		UserID:    user.UserID,
		Email:     user.Email,
		Username:  user.Username,
		FullName:  user.FullName,
		AvatarURL: user.AvatarURL,
		CreatedAt: user.CreatedAt,
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Account created successfully",
		"data":    response,
	})
}
