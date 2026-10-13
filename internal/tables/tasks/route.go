package tasks

import (
	"orbit-backend-golang/internal/security"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TaskRoutes(api *gin.RouterGroup, db *pgxpool.Pool, jwtService *security.JWTService) {
	repository := NewTaskRepository(db)
	service := NewTaskService(repository)
	handler := NewTaskHandler(service)

	tasks := api.Group("/tasks")

	tasks.POST("", handler.CreateTaskHandler)
	tasks.DELETE("/:taskID")
	tasks.GET("/:taskID")
	tasks.GET("/workspace/:workspaceID ")
}
