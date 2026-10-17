package connections

import (
	"context"
	"fmt"

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
		RETURNING connection_request_id;
	`
	var connection_req_id string

	err := repo.db.QueryRow(ctx, query, senderID, receiverID).Scan(&connection_req_id)

	if err != nil {
		fmt.Printf("Create Connection Request Error: %v", err)
		return "", err
	}

	return connection_req_id, nil
}

func (repo *ConnectionRepository) CreateNewConnection(ctx context.Context, senderID, receiverID string) (string, error) {
	query1 := `
		INSERT INTO connections
		(sender_id, receiver_id)
		VALUES ($1, $2)
		RETURNING connection_id;
	`

	var connection_id string

	err := repo.db.QueryRow(ctx, query1, senderID, receiverID).Scan(&connection_id)

	if err != nil {
		fmt.Printf("Create Connection Error: %v", err)
		return "", err
	}

	return connection_id, nil
}

func (repo *ConnectionRepository) GetConnectionRequestByID(ctx context.Context, connection_request_id string) (*Connection_Request, error) {
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

func (repo *ConnectionRepository) GetConnectionByID(ctx context.Context, connection_request_id string) (*Connection, error) {
	query := `
		SELECT 
			connection_id,
			sender_id,
			receiver_id,
			created_at
		FROM connections
		WHERE connection_id=$1
	`

	connection := &Connection{}
	err := repo.db.QueryRow(ctx, query, connection_request_id).Scan(&connection.Connection_ID, &connection.SenderID, &connection.ReceiverID, &connection.Created_At)

	return connection, err
}

func (repo *ConnectionRepository) GetConnectionRequest(ctx context.Context, sender_id, receiver_id string) (*Connection_Request, error) {
	query := `
		SELECT 
			connection_request_id,
			sender_id,
			receiver_id,
			created_at
		FROM connection_requests
		WHERE sender_id=$1 OR receiver_id=$2;
	`

	connection_request := &Connection_Request{}
	err := repo.db.QueryRow(ctx, query, sender_id, receiver_id).Scan(&connection_request.Connection_Request_ID, &connection_request.SenderID, &connection_request.ReceiverID, &connection_request.Created_At)
	return connection_request, err
}

func (repo *ConnectionRepository) GetAllConnectionRequest(ctx context.Context, userID string) (*Connection_Request, error) {
	query := `
		SELECT 
			connection_request_id,
			sender_id,
			receiver_id,
			created_at
		FROM connection_requests
		WHERE receiver_id=$1;
	`

	connection_request := &Connection_Request{}
	err := repo.db.QueryRow(ctx, query, userID).Scan(&connection_request.Connection_Request_ID, &connection_request.SenderID, &connection_request.ReceiverID, &connection_request.Created_At)
	return connection_request, err
}

func (repo *ConnectionRepository) GetConnection(ctx context.Context, sender_id, receiver_id string) (*Connection, error) {
	query := `
		SELECT 
			connection_id,
			sender_id,
			receiver_id,
			created_at
		FROM connections
		WHERE sender_id=$1 AND receiver_id=$2;
	`

	connection_request := &Connection{}
	err := repo.db.QueryRow(ctx, query, sender_id, receiver_id).Scan(&connection_request.Connection_ID, &connection_request.SenderID, &connection_request.ReceiverID, &connection_request.Created_At)
	return connection_request, err
}

func (repo *ConnectionRepository) DeleteConnectionRequest(ctx context.Context, sender_id, receiver_id string) error {
	query := `
		DELETE FROM connection_requests
		WHERE sender_id=$1 AND receiver_id=$2;
	`

	result, err := repo.db.Exec(ctx, query, sender_id, receiver_id)

	if err != nil {
		return fmt.Errorf("Delete Connection Request: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrConnectionRequestNotFound
	}

	return nil
}

func (repo *ConnectionRepository) DeleteConnection(ctx context.Context, sender_id, receiver_id string) error {
	query := `
		DELETE FROM connections
		WHERE sender_id=$1 AND receiver_id=$2;
	`

	result, err := repo.db.Exec(ctx, query, sender_id, receiver_id)

	if err != nil {
		return fmt.Errorf("Delete Connection: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrConnectionNotFound
	}

	return nil
}

func (repo *ConnectionRepository) DeleteConnectionRequestByID(ctx context.Context, connection_req_id string) error {
	query := `
		DELETE FROM connection_requests
		WHERE connection_req_id=$1;
	`

	result, err := repo.db.Exec(ctx, query, connection_req_id)

	if err != nil {
		return fmt.Errorf("Delete Connection Request by ID error: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrConnectionRequestNotFound
	}

	return nil
}

func (repo *ConnectionRepository) DeleteConnectionByID(ctx context.Context, connection_id string) error {
	query := `
		DELETE FROM connections
		WHERE connection_id=$1;
	`

	result, err := repo.db.Exec(ctx, query, connection_id)

	if err != nil {
		return fmt.Errorf("Delete Connection by ID error: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrConnectionNotFound
	}

	return nil
}
