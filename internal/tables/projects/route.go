package projects

import (
	"orbit-backend-golang/internal/middleware"
	"orbit-backend-golang/internal/security"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ProjectRoutes(api *gin.RouterGroup, db *pgxpool.Pool, jwt *security.JWTService) {
	repository := NewProjectRepository(db)
	service := NewProjectService(repository)
	handler := NewProjectHandler(service)

	projects := api.Group("/projects")

	projects.Use(middleware.AuthMiddleware(jwt))

	projects.POST("", handler.createProject)
	projects.GET("/:workspaceID", handler.getProjectsByWorkspaceID)
	// projects.GET("/:projectID", handler.getProjectById)
	projects.DELETE("/:projectID", handler.DeleteProject)
}
