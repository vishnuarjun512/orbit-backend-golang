package home

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HomeHandler struct {
	service *HomeService
}

func NewHomeHanlder(service *HomeService) *HomeHandler {
	return &HomeHandler{service: service}
}

func (h *HomeHandler) GetHomeDataHandler(c *gin.Context) {

	userID := c.GetString("userID")
	if userID == "" {
		fmt.Printf("UserID Empty")
		c.JSON(http.StatusExpectationFailed, gin.H{"error": true, "message": "UserID Empty"})
		return
	}

	workspaces, err := h.service.workRepository.GetWorkSpaces(c.Request.Context(), userID, 5, 0)
	if err != nil {
		fmt.Printf("Workspaces Fetch Error : %v", err)
		c.JSON(http.StatusGone, gin.H{"error": true, "message": "Workspaces Fetch failed"})
		return
	}

	projects, err := h.service.projectRepository.GetProjectsByWorkspaceIDRepo(c.Request.Context(), workspaces[1].WorkSpaceID, 5, 0)
	if err != nil {
		fmt.Printf("Projects Fetch Error : %v", err)
		c.JSON(http.StatusGone, gin.H{"error": true, "message": "Projects Fetch failed"})
		return
	}

	tasks, err := h.service.taskRepository.GetTasksByWorkspaceID(c.Request.Context(), workspaces[0].WorkSpaceID)
	if err != nil {
		fmt.Printf("Task Fetch Error : %v", err)
		c.JSON(http.StatusGone, gin.H{"error": true, "message": "Task Fetch failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Home Data Sent", "tasks": tasks, "workspaces": workspaces, "projects": projects})
}
