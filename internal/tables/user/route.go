package user

import (
    "github.com/gin-gonic/gin"
    "github.com/jackc/pgx/v5/pgxpool"
)

func UserRoutes(router *gin.RouterGroup, db *pgxpool.Pool) {
    repository := NewRepository(db)
    service := NewService(repository)
    handler := NewHandler(service)

    auth := router.Group("/auth")
    {
        auth.POST("/register", handler.Register)
    }
}