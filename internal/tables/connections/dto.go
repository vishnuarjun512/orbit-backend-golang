package connections

import "time"

type Connection_Request struct {
	ConnectionRequestID string    `json:"connection_request_id"`
	SenderID            string    `json:"sender_id"`
	ReceiverID          string    `json:"receiver_id"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type Connection struct {
	ConnectionID string    `json:"connection_id"`
	UserAID      string    `json:"user_a_id"`
	UserBID      string    `json:"user_b_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ConnectedUserResponse struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

type ConnectionResponse struct {
	ConnectionID  string                `json:"connection_id"`
	ConnectedUser ConnectedUserResponse `json:"connected_user"`
	CreatedAt     time.Time             `json:"created_at"`
}
