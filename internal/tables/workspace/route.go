package workspace

import (
	"orbit-backend-golang/internal/middleware"
	"orbit-backend-golang/internal/security"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func WorkspaceRoutes(
	api *gin.RouterGroup,
	db *pgxpool.Pool,
	jwtService *security.JWTService,
) {
	repository := NewWorkSpaceRepository(db)
	service := NewWorkSpaceService(repository)
	handler := NewWorkSpaceHandler(service)

	workspaces := api.Group("/workspace")
	workspaces.Use(middleware.AuthMiddleware(jwtService))

	// POST /api/workspace/
	// Auth: orbit_access_token cookie or Authorization: Bearer <JWT>.
	// Body: {"name":"Team name"}; optional fields: "slug" (generated from name when omitted)
	// and "description" (string). No path parameters or query parameters.
	workspaces.POST("/", handler.createWorkspace)

	// GET /api/workspace/
	// Auth: orbit_access_token cookie or Authorization: Bearer <JWT>.
	// Body: none. Query: page (1-based, default 1), limit (1-100, default 20).
	workspaces.GET("/", handler.getWorkspaces)

	// GET /api/workspace/:workspaceID
	// Auth: orbit_access_token cookie or Authorization: Bearer <JWT>.
	// Body: none. Path: workspaceID (workspace UUID). Query: none.
	workspaces.GET("/:workspaceID", handler.getWorkspace)

	// PATCH /api/workspace/:workspaceID
	// Auth: orbit_access_token cookie or Authorization: Bearer <JWT>.
	// Body: at least one of "name", "slug", or "description"; send description as ""
	// to clear it. Path: workspaceID (workspace UUID). Query: none.
	workspaces.PATCH("/:workspaceID", handler.updateWorkspace)

	// DELETE /api/workspace/:workspaceID
	// Auth: orbit_access_token cookie or Authorization: Bearer <JWT>.
	// Body: none. Path: workspaceID (workspace UUID). Query: none.
	workspaces.DELETE("/:workspaceID", handler.deleteWorkspace)
}
