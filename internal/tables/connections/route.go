package connections

import (
	"orbit-backend-golang/internal/middleware"
	"orbit-backend-golang/internal/security"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectionRoutes(api *gin.RouterGroup, db *pgxpool.Pool, jwtService *security.JWTService) {

	repository := NewConnectionRepository(db)
	service := NewConnectionService(repository)
	handler := NewConnectionHandler(service)

	connection := api.Group("/connection")

	connection.Use(middleware.AuthMiddleware(jwtService))
	connection.GET("/", handler.GetConnections)
	connection.GET("/:receiverID/request", handler.CreateConnectionRequest)

	connection.GET("/:receiverID/accept", handler.AcceptConnectionHandler)
	connection.GET("/:receiverID/reject", handler.RejectConnectionHandler)
	connection.GET("/:receiverID/delete", handler.DeleteConnectionHandler)

}
