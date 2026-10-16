package connections

import (
	"orbit-backend-golang/internal/middleware"
	"orbit-backend-golang/internal/security"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectionRoutes(api *gin.RouterGroup, db *pgxpool.Pool, jwtService *security.JWTService) {

	connection := api.Group("/connection")

	connection.Use(middleware.AuthMiddleware(jwtService))
	connection.POST("/")
}
