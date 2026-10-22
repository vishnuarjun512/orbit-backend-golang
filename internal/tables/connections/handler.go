package connections

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ConnectionHandler struct {
	service *ConnectionService
}

func NewConnectionHandler(service *ConnectionService) *ConnectionHandler {
	return &ConnectionHandler{service: service}
}

func (h *ConnectionHandler) CreateConnectionRequest(c *gin.Context) {
	userID := GetUserID(c)

	receiverID := GetReceiverID(c)

	connection_request, err := h.service.CreateNewConnectionRequestService(c.Request.Context(), userID, receiverID)
	if err != nil {
		c.JSON(http.StatusGone, gin.H{"error": true, "message": "Creation Connection Failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"error": false, "message": "Created Connection", "connection": connection_request})
}

func (h *ConnectionHandler) GetReceivedRequests(c *gin.Context) {
	userID := GetUserID(c)

	connection_requests, err := h.service.GetReceivedRequestsService(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusGone, gin.H{"error": true, "message": "Fetch Received Reuqests Failed"})
		return
	}

	c.JSON(http.StatusFound, gin.H{"error": false, "message": "Fetch Received Reuqests", "connection": connection_requests})
}

func (h *ConnectionHandler) GetSentRequests(c *gin.Context) {
	userID := GetUserID(c)

	connection_requests, err := h.service.GetSentRequestsService(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusGone, gin.H{"error": true, "message": "Fetch Sent Reuqests Failed"})
		return
	}

	c.JSON(http.StatusFound, gin.H{"error": false, "message": "Fetch Sent Reuqests", "connection": connection_requests})
}

func (h *ConnectionHandler) AcceptConnectionHandler(c *gin.Context) {
	senderID, receiverID := GetSenderAndReceiverID(c)

	connection, err := h.service.AcceptConnectionRequestService(c.Request.Context(), senderID, receiverID)
	if err != nil {
		fmt.Printf("Accept Connection Request Error : %v", err)
		c.JSON(http.StatusGone, gin.H{"error": true, "message": "Accept Connection Request Failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"error": false, "message": "Connection Request Accepted", "connection": connection})
}

func (h *ConnectionHandler) RejectConnectionHandler(c *gin.Context) {
	senderID, receiverID := GetSenderAndReceiverID(c)

	err := h.service.RejectConnectionRequestService(c.Request.Context(), senderID, receiverID)
	if err != nil {
		fmt.Printf("Reject Connection Request Error : %v", err)
		c.JSON(http.StatusGone, gin.H{"error": true, "message": "Reject Connection Request Failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Reject Request"})
}

func (h *ConnectionHandler) DeleteConnectionHandler(c *gin.Context) {
	senderID, receiverID := GetSenderAndReceiverID(c)

	err := h.service.DeleteConnectionService(c.Request.Context(), senderID, receiverID)
	if err != nil {
		fmt.Printf("Delete Connection Error : %v", err)
		c.JSON(http.StatusGone, gin.H{"error": true, "message": "Delete Connection Failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Delete Connection"})
}

func (h *ConnectionHandler) GetAllConnections(c *gin.Context) {
	userID := GetUserID(c)

	connections, err := h.service.GetAllConnectionsService(c.Request.Context(), userID)

	if err != nil {
		fmt.Printf("Get connections failed: %v\n", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "Failed to fetch connections",
		})
		return
	}

	c.JSON(http.StatusFound, gin.H{"error": false, "message": "Fetch Connections", "connections": connections})

}

func GetSenderAndReceiverID(c *gin.Context) (string, string) {
	userID := c.GetString("userID")
	if userID == "" {
		fmt.Println("SenderID Empty")
		c.JSON(http.StatusExpectationFailed, gin.H{"error": true, "message": "SenderID Empty"})
		return "", ""
	}

	receiverID := c.Param("receiverID")
	if receiverID == "" {
		fmt.Println("ReceiverID Empty")
		c.JSON(http.StatusExpectationFailed, gin.H{"error": true, "message": "ReceiverID Empty"})
		return "", ""
	}
	return userID, receiverID
}

func GetUserID(c *gin.Context) string {
	userID := c.GetString("userID")
	if userID == "" {
		fmt.Println("UserID Empty")
		c.JSON(http.StatusExpectationFailed, gin.H{"error": true, "message": "UserID Empty"})
		return ""
	}
	return userID
}

func GetReceiverID(c *gin.Context) string {
	receiverID := c.Param("receiverID")
	if receiverID == "" {
		fmt.Println("ReceiverID Empty")
		c.JSON(http.StatusExpectationFailed, gin.H{"error": true, "message": "ReceiverID Empty"})
		return ""
	}
	return receiverID
}
