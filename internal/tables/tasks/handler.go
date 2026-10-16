package tasks

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	service *TaskService
}

func NewTaskHandler(service *TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) CreateTaskHandler(c *gin.Context) {
	var req TaskCreateRequest
	if err := c.ShouldBind(&req); err != nil {
		fmt.Printf("Invalid request fields : %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Invalid Request fields"})
	}

	task, err := h.service.CreateTaskService(c.Request.Context(), c.GetString("userID"), &req)
	if err != nil {
		fmt.Printf("Task Creation Error : %v", err)
		c.JSON(http.StatusCreated, gin.H{"error": true, "message": "Task Creation Failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"error": false, "message": "Task Created", "task": task})
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		fmt.Printf("UserID Empty")
		c.JSON(http.StatusExpectationFailed, gin.H{"error": true, "message": "UserID Empty"})
		return
	}

	workspaceID := c.Param("workspaceID")
	if workspaceID == "" {
		fmt.Printf("WorkspaceID Empty")
		c.JSON(http.StatusExpectationFailed, gin.H{"error": true, "message": "WorkspaceID Empty"})
		return
	}

	tasks, err := h.service.GetTasksByWorkspaceID(c.Request.Context(), userID)

	if err != nil {
		fmt.Printf("Task Fetch Error : %v", err)
		c.JSON(http.StatusGone, gin.H{"error": true, "message": "Task Fetch failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"error": false, "message": "Task Created", "tasks": tasks})
}

func (h *TaskHandler) GetTasksByProjectID(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		fmt.Printf("UserID Empty")
		c.JSON(http.StatusExpectationFailed, gin.H{"error": true, "message": "UserID Empty"})
		return
	}

	projectID := c.Param("projectID")
	if projectID == "" {
		fmt.Printf("ProjectID Empty")
		c.JSON(http.StatusExpectationFailed, gin.H{"error": true, "message": "ProjectID Empty"})
		return
	}

	tasks, err := h.service.GetTasksByWorkspaceID(c.Request.Context(), projectID)

	if err != nil {
		fmt.Printf("Task Fetch Error : %v", err)
		c.JSON(http.StatusGone, gin.H{"error": true, "message": "Task Fetch failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Tasks Fetched", "tasks": tasks})
}
