package workspace

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

var workspaceIDPattern = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

type WorkspaceHandler struct {
	service *WorkSpaceService
}

func NewWorkSpaceHandler(service *WorkSpaceService) *WorkspaceHandler {
	return &WorkspaceHandler{service: service}
}

func (h *WorkspaceHandler) createWorkspace(c *gin.Context) {
	var req WorkSpaceCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Invalid workspace details"})
		return
	}

	workspace, err := h.service.CreateWorkspace(c.Request.Context(), c.GetString("userID"), req)
	if err != nil {
		h.respondWorkspaceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"error":     false,
		"workspace": workspace,
	})
}

func (h *WorkspaceHandler) getWorkspaces(c *gin.Context) {
	page, limit, ok := workspacePagination(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "page must be positive and limit must be between 1 and 100"})
		return
	}

	workspaces, err := h.service.GetWorkSpacesService(
		c.Request.Context(),
		c.GetString("userID"),
		page,
		limit,
	)

	if err != nil {
		h.respondWorkspaceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error":      false,
		"workspaces": workspaces,
		"pagination": gin.H{"page": page, "limit": limit},
	})
}

func (h *WorkspaceHandler) getWorkspace(c *gin.Context) {
	id, ok := workspaceID(c)
	if !ok {
		return
	}
	workspace, err := h.service.GetWorkspace(c.Request.Context(), id, c.GetString("userID"))
	if err != nil {
		h.respondWorkspaceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "workspace": workspace})
}

func (h *WorkspaceHandler) updateWorkspace(c *gin.Context) {
	id, ok := workspaceID(c)
	if !ok {
		return
	}
	var req WorkSpaceUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "Invalid workspace details"})
		return
	}
	workspace, err := h.service.UpdateWorkspace(c.Request.Context(), id, c.GetString("userID"), req)
	if err != nil {
		h.respondWorkspaceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"error": false, "workspace": workspace})
}

func (h *WorkspaceHandler) deleteWorkspace(c *gin.Context) {
	id, ok := workspaceID(c)
	if !ok {
		return
	}
	if err := h.service.DeleteWorkspace(c.Request.Context(), id, c.GetString("userID")); err != nil {
		h.respondWorkspaceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func workspaceID(c *gin.Context) (string, bool) {
	id := c.Param("workspaceID")
	if !workspaceIDPattern.MatchString(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": "workspaceID must be a UUID"})
		return "", false
	}
	return id, true
}

func workspacePagination(c *gin.Context) (int, int, bool) {
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

func (h *WorkspaceHandler) respondWorkspaceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalidWorkspaceInput):
		c.JSON(http.StatusBadRequest, gin.H{"error": true, "message": err.Error()})
	case errors.Is(err, ErrWorkspaceNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": true, "message": "Workspace not found"})
	case errors.Is(err, ErrWorkspaceSlugConflict):
		c.JSON(http.StatusConflict, gin.H{"error": true, "message": "Workspace slug already exists"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": true, "message": "Something went wrong"})
	}
}
