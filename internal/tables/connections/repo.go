package connections

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ConnectionRepository struct {
	db *pgxpool.Pool
}

func NewConnectionRepository(db *pgxpool.Pool) *ConnectionRepository {
	return &ConnectionRepository{
		db: db,
	}
}

func (repo *ConnectionRepository) CreateNewConnectionRequest(ctx context.Context, senderID, receiverID string) (string, error) {
	query := `
		INSERT INTO connection_requests
		(sender_id, receiver_id)
		VALUES ($1, $2)
		RETURNING id;
	`
	var connection_req_id string
	err := repo.db.QueryRow(ctx, query, senderID, receiverID).Scan(&connection_req_id)

	return connection_req_id, err
}

func (repo *ConnectionRepository) GetConnectionRequest(ctx context.Context, connection_request_id string) (*Connection_Request, error) {
	query := `
		SELECT 
			connection_request_id,
			sender_id,
			receiver_id,
			created_at
		FROM connection_requests
		WHERE connection_request_id=$1
	`

	connection_request := &Connection_Request{}
	err := repo.db.QueryRow(ctx, query, connection_request_id).Scan(&connection_request.Connection_Request_ID, &connection_request.SenderID, &connection_request.ReceiverID, &connection_request.Created_At)

	return connection_request, err
}
