package user

import (
	"orbit-backend-golang/internal/security"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UserRoutes(router *gin.RouterGroup, db *pgxpool.Pool, jwtService *security.JWTService) {

	repository := NewRepository(db)
	service := NewService(repository, jwtService)
	handler := NewHandler(service)

	auth := router.Group("/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
	}
}
