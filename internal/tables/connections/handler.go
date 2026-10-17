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
		c.JSON(http.StatusGone, gin.H{"error": true, "message": "Creation Connection Failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"error": false, "message": "Created Connection", "connection": connection_request})
}

func (h *ConnectionHandler) GetConnections(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		fmt.Println("SenderID Empty")
		c.JSON(http.StatusExpectationFailed, gin.H{"error": true, "message": "SenderID Empty"})
		return
	}

	receiverID := c.Param("receiverID")
	if receiverID == "" {
		fmt.Println("ReceiverID Empty")
		c.JSON(http.StatusExpectationFailed, gin.H{"error": true, "message": "ReceiverID Empty"})
		return
	}

	connection_request, err := h.service.GetConnectionRequestService(c.Request.Context(), userID, receiverID)

	if err != nil {
		fmt.Printf("Fetch Connection Error : %v", err)
		c.JSON(http.StatusGone, gin.H{"error": true, "message": "Fetch Connection Request Failed"})
		return
	}

	c.JSON(http.StatusFound, gin.H{"error": false, "message": "Fetch Connection Successs", "connection": connection_request})
}

func (h *ConnectionHandler) AcceptConnectionHandler(c *gin.Context) {
	senderID, receiverID := GetSenderAndReceiverID(c)

	connection, err := h.service.AcceptConnectionRequestService(c.Request.Context(), senderID, receiverID)
	if err != nil {
		fmt.Printf("Accept Connection Request Error : %v", err)
		c.JSON(http.StatusGone, gin.H{"error": true, "message": "Accept Connection Request Failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Connection Request Accepted", "connection": connection})
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
