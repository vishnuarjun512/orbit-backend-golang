package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"orbit-backend-golang/internal/tables/user"
)

func MainRouter(db *pgxpool.Pool) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "Success",
			"message": "API is running",
		})
	})

	api := router.Group("/api")

	user.UserRoutes(api, db)

	return router
}
