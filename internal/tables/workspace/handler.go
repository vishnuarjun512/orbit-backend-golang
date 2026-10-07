package workspace

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WorkspaceHandler struct {
	service *WorkSpaceService
}

func NewWorkSpaceHandler(service *WorkSpaceService) *WorkspaceHandler {
	return &WorkspaceHandler{
		service: service,
	}
}

func (h *WorkspaceHandler) createWorkspace(c *gin.Context) {

}

func (h *WorkspaceHandler) GetWorkspaces(c *gin.Context) {
	userID := c.GetString("userID")

	if userID == "" {
		fmt.Println("Unauthorized")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Authentication required",
		})
	}

	page := 1
	limit := 20

	if value := c.Query("page"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			page = parsed
		}
	}

	if value := c.Query("limit"); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			limit = parsed
		}
	}

	workspaces, err := h.service.GetWorkSpacesService(
		c.Request.Context(),
		userID,
		page,
		limit,
	)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusConflict, gin.H{
			"error":   true,
			"message": "Workspaces Fetch Error",
		})
	}

	c.JSON(http.StatusFound, gin.H{
		"error":      false,
		"message":    "Login successful",
		"workspaces": workspaces,
	})
}
