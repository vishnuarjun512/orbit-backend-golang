package connections

import (
	"orbit-backend-golang/internal/middleware"
	"orbit-backend-golang/internal/security"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectionRoutes(
	api *gin.RouterGroup,
	db *pgxpool.Pool,
	jwtService *security.JWTService,
) {
	repository := NewConnectionRepository(db)
	service := NewConnectionService(repository)
	handler := NewConnectionHandler(service)

	connection := api.Group("/connections")
	connection.Use(middleware.AuthMiddleware(jwtService))

	connection.GET(
		"",
		handler.GetAllConnections,
	)

	connection.POST(
		"/requests/:receiverID",
		handler.CreateConnectionRequest,
	)

	connection.GET(
		"/requests/received",
		handler.GetReceivedRequests,
	)

	connection.GET(
		"/requests/sent",
		handler.GetSentRequests,
	)

	connection.POST(
		"/requests/:requestID/accept",
		handler.AcceptConnectionHandler,
	)

	connection.POST(
		"/requests/:requestID/reject",
		handler.RejectConnectionHandler,
	)

	connection.DELETE(
		"/:userID",
		handler.DeleteConnectionHandler,
	)
}
