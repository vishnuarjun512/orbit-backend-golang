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

func (h *ConnectionHandler) CreateNewConnectionHandler(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		fmt.Printf("No senderID Empty")
		c.JSON(http.StatusExpectationFailed, gin.H{"error": true, "message": "SenderID Empty"})
		return
	}

	receiverID := c.Param("receiverID")
	if receiverID == "" {
		fmt.Printf("ReceiverID Empty")
		c.JSON(http.StatusExpectationFailed, gin.H{"error": true, "message": "ReceiverID Empty"})
		return
	}

	connection_request, err := h.service.CreateNewConnectionRequestService(c.Request.Context(), userID, receiverID)
	if err != nil {
		fmt.Printf("Create Connection Error : %v", err)
		c.JSON(http.StatusGone, gin.H{"error": true, "message": "Creation Connection Failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"error": false, "message": "Created Connection", "connection": connection_request})
}

func (h *ConnectionHandler) GetConnections(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		fmt.Printf("No senderID Empty")
		c.JSON(http.StatusExpectationFailed, gin.H{"error": true, "message": "SenderID Empty"})
		return
	}

	receiverID := c.Param("receiverID")
	if receiverID == "" {
		fmt.Printf("ReceiverID Empty")
		c.JSON(http.StatusExpectationFailed, gin.H{"error": true, "message": "ReceiverID Empty"})
		return
	}

	connection_request, err := h.service.CreateNewConnectionRequestService(c.Request.Context(), userID, receiverID)
	if err != nil {
		fmt.Printf("Create Connection Error : %v", err)
		c.JSON(http.StatusGone, gin.H{"error": true, "message": "Creation Connection Failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"error": false, "message": "Created Connection", "connection": connection_request})
}
