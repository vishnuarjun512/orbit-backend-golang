package projects

import (
	"fmt"
	"net/http"
	"strconv"

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

func (h *ProjectHandler) createProject(c *gin.Context) {
	var req ProjectCreateRequest
	userID := c.GetString("userID")
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Invalid Project Details"})
		return
	}

	project, err := h.service.CreateProjectService(c.Request.Context(), userID, req)
	if err != nil {
		fmt.Printf("Create Project Failed: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Create Project Failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Created Project Successfully", "project": project})
}

func (h *ProjectHandler) getProjectById(c *gin.Context) {
	projectID := c.Param("projectID")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Invalid Project ID"})
		return
	}

	project, err := h.service.GetProjectByIDService(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Get Project Failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Get Project Successful", "project": project})
}

func (h *ProjectHandler) getProjectsByWorkspaceID(c *gin.Context) {
	workspaceID := c.Param("workspaceID")
	if workspaceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Invalid Workspace ID"})
		return
	}

	page, limit, ok := projectPagination(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "page must be positive and limit must be between 1 and 100"})
		return
	}

	projects, err := h.service.GetProjectsByWorkspaceIDService(c.Request.Context(), workspaceID, page, limit)
	if err != nil {
		fmt.Printf("Get Projects By Workspace Failed: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Get Projects By Workspace Failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Get Projects By Workspace Successful", "projects": projects})
}

func projectPagination(c *gin.Context) (int, int, bool) {
	page, limit := 1, 20
	if value := c.Query("page"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 {
			return 0, 0, false
		}
		page = parsed
	}
	if value := c.Query("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 || parsed > 100 {
			return 0, 0, false
		}
		limit = parsed
	}
	maxInt := int(^uint(0) >> 1)
	if page-1 > maxInt/limit {
		return 0, 0, false
	}
	return page, limit, true
}

func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	projectID := c.Param("projectID")
	userID := c.GetString("userID")

	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Invalid Project ID"})
		return
	}

	project, err := h.service.GetProjectByIDService(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Get Project Failed"})
		return
	}

	if project.Created_By != userID {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Project cannot be deleted."})
		return
	}

	delErr := h.service.DeleteProjectByIDService(c.Request.Context(), projectID)

	if delErr != nil {
		print("Delete Error -> ", delErr)
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Delete Project Failed"})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": false, "message": "Delete Project Successful"})
}
