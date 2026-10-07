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
	workspaceRepo := NewWorkSpaceRepository(db)
	workspaceService := NewWorkSpaceService(workspaceRepo)
	workspaceHandler := NewWorkSpaceHandler(workspaceService)

	workspace := api.Group("/workspace")

	workspace.Use(
		middleware.AuthMiddleware(jwtService),
	)

	workspace.POST("/", workspaceHandler.createWorkspace)
}
