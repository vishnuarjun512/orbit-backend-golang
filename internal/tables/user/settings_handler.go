package user

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Settings(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "Unauthorized"})
		return
	}
	var req PreferencesUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Invalid settings"})
		return
	}

	if req.Timezone != nil {
		zone := strings.TrimSpace(*req.Timezone)
		if _, err := time.LoadLocation(zone); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Invalid timezone"})
			return
		}
		req.Timezone = &zone
	}

	if req.WeekStartsOn != nil && *req.WeekStartsOn != 0 && *req.WeekStartsOn != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "week_starts_on must be 0 or 1"})
		return
	}

	if req.Theme != nil && *req.Theme != "light" && *req.Theme != "dark" && *req.Theme != "system" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Invalid theme"})
		return
	}

	p, err := h.service.repo.UpdatePreferences(c.Request.Context(), userID.(string), req)
	if err != nil {
		if err == ErrInvalidPreferences {
			c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Active workspace must be an active workspace member"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Something went wrong"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "preferences": p})
}

func (h *Handler) GetSettings(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "Unauthorized"})
		return
	}
	p, err := h.service.repo.GetPreferences(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Something went wrong"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "preferences": p})
}

func (h *Handler) GetProfile(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "Unauthorized"})
		return
	}
	p, err := h.service.repo.GetProfile(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Something went wrong"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "profile": p})
}

func (h *Handler) LookupProfile(c *gin.Context) {
	profile, err := h.service.repo.GetProfileByPublicID(c.Request.Context(), c.Param("publicID"))
	if err != nil {
		if err == ErrUserNotRegistered {
			c.JSON(http.StatusNotFound, gin.H{"error": true, "message": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Something went wrong"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "profile": profile})
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "Unauthorized"})
		return
	}
	var req ProfileUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Invalid profile"})
		return
	}
	if req.DisplayName != nil && strings.TrimSpace(*req.DisplayName) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "display_name cannot be empty"})
		return
	}
	if req.Bio != nil && len([]rune(*req.Bio)) > 500 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "bio must be 500 characters or fewer"})
		return
	}
	if req.StatusMessage != nil && len([]rune(*req.StatusMessage)) > 160 {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "status_message must be 160 characters or fewer"})
		return
	}
	if req.Availability != nil && *req.Availability != "available" && *req.Availability != "away" && *req.Availability != "do_not_disturb" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Invalid availability"})
		return
	}
	p, err := h.service.repo.UpdateProfile(c.Request.Context(), userID.(string), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Something went wrong"})
		return
	}
	full, err := h.service.repo.GetProfile(c.Request.Context(), userID.(string))
	if err == nil {
		p.PublicID = full.PublicID
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "profile": p})
}
