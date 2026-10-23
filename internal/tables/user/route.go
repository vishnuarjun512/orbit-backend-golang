package user

import (
	"orbit-backend-golang/internal/middleware"
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
		auth.POST("/refresh", handler.Refresh)
	}

	user := router.Group("/user")
	user.Use(middleware.AuthMiddleware(jwtService))
	{
		user.GET("/directory/:publicID", handler.LookupProfile)
		user.GET("/profile", handler.GetProfile)
		user.PATCH("/profile", handler.UpdateProfile)
		user.GET("/settings", handler.GetSettings)
		user.PATCH("/settings", handler.Settings)
	}
}
