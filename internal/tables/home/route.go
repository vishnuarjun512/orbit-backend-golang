package home

import (
	"orbit-backend-golang/internal/middleware"
	"orbit-backend-golang/internal/security"

	"orbit-backend-golang/internal/tables/projects"
	"orbit-backend-golang/internal/tables/tasks"
	"orbit-backend-golang/internal/tables/workspace"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func HomeRoutes(api *gin.RouterGroup, db *pgxpool.Pool, jwtService *security.JWTService) {
	workspaceRepo := workspace.NewWorkSpaceRepository(db)
	tasksRepo := tasks.NewTaskRepository(db)
	projectsRepo := projects.NewProjectRepository(db)

	service := NewHomeService(workspaceRepo, tasksRepo, projectsRepo)
	handler := NewHomeHanlder(service)

	home := api.Group("/home")
	home.Use(middleware.AuthMiddleware(jwtService))

	home.GET("", handler.GetHomeDataHandler)

}
