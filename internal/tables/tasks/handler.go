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
