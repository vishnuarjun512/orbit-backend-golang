package projects

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	service *ProjectService
}

func NewProjectHandler(service *ProjectService) *ProjectHandler {
	return &ProjectHandler{
		service: service,
	}
}

func (h *ProjectHandler) createProjectHandler(c *gin.Context) {
	var req ProjectCreateRequest
	userID := c.GetString("userID")
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Invalid Project Details"})
	}

	project, err := h.service.CreateProjectService(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Create Project Failed"})
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": false, "message": "Created Project Successfully", "project": project})
}

func (h *ProjectHandler) getProjectByIdHandler(c *gin.Context) {
	workspaceID := c.Param("workspaceID")
	if workspaceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Invalid Project ID"})
	}

	project, err := h.service.GetProjectByIDService(c.Request.Context(), workspaceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Get Project Failed"})
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": false, "message": "Get Project Successful", "project": project})
}
