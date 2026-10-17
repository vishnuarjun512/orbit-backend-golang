package connections

import "time"

type Connection_Request struct {
	Connection_Request_ID string    `json:"connection_request_id"`
	SenderID              string    `json:"sender_id"`
	ReceiverID            string    `json:"receiver_id"`
	Created_At            time.Time `json:"created_at"`
	Updated_At            time.Time `json:"updated_at"`
}

type Connection struct {
	Connection_ID string    `json:"connection_id"`
	SenderID      string    `json:"sender_id"`
	ReceiverID    string    `json:"receiver_id"`
	Created_At    time.Time `json:"created_at"`
	Updated_At    time.Time `json:"updated_at"`
}
