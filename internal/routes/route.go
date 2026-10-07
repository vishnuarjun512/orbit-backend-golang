package routes

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"orbit-backend-golang/internal/security"
	"orbit-backend-golang/internal/tables/user"
	"orbit-backend-golang/internal/tables/workspace"
)

func MainRouter(db *pgxpool.Pool, jwtService *security.JWTService) *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
		},
		AllowCredentials: true,
	}))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "Success",
			"message": "API is running",
		})
	})

	api := router.Group("/api")

	user.UserRoutes(api, db, jwtService)
	workspace.WorkspaceRoutes(api, db, jwtService)

	return router
}
