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
	err := repo.db.QueryRow(ctx, query, connection_request_id).Scan(&connection_request.ConnectionRequestID, &connection_request.SenderID, &connection_request.ReceiverID, &connection_request.CreatedAt)

	return connection_request, err
}

func (repo *ConnectionRepository) GetConnectionRequest(ctx context.Context, sender_id, receiver_id string) (*Connection_Request, error) {
	query := `
		SELECT 
			connection_request_id,
			sender_id,
			receiver_id,
			created_at
		FROM connection_requests
		WHERE sender_id=$1 AND receiver_id=$2;
	`

	connection_request := &Connection_Request{}
	err := repo.db.QueryRow(ctx, query, sender_id, receiver_id).Scan(&connection_request.ConnectionRequestID, &connection_request.SenderID, &connection_request.ReceiverID, &connection_request.CreatedAt)
	return connection_request, err
}

func (repo *ConnectionRepository) GetReceivedRequests(ctx context.Context, receiver_id string) (*[]Connection_Request, error) {
	query := `
		SELECT 
			connection_request_id,
			sender_id,
			receiver_id,
			created_at
		FROM connection_requests
		WHERE receiver_id=$1;
	`

	results, err := repo.db.Query(ctx, query, receiver_id)

	if err != nil {
		return nil, fmt.Errorf("List Received Requests: %w", err)
	}
	defer results.Close()

	connection_requests := make([]Connection_Request, 0)

	for results.Next() {
		connection, err := scanConnectionRequest(results)
		if err != nil {
			return nil, fmt.Errorf("Scan Connection Request : %w", err)
		}
		connection_requests = append(connection_requests, *connection)
	}
	if err := results.Err(); err != nil {
		return nil, fmt.Errorf("Iterate connection requests: %w", err)
	}

	return &connection_requests, err
}

func (repo *ConnectionRepository) GetSentRequests(ctx context.Context, sender_id string) (*[]Connection_Request, error) {
	query := `
		SELECT 
			connection_request_id,
			sender_id,
			receiver_id,
			created_at
		FROM connection_requests
		WHERE sender_id=$1;
	`

	results, err := repo.db.Query(ctx, query, sender_id)

	if err != nil {
		return nil, fmt.Errorf("List Sent Requests: %w", err)
	}
	defer results.Close()

	connection_requests := make([]Connection_Request, 0)

	for results.Next() {
		connection, err := scanConnectionRequest(results)
		if err != nil {
			return nil, fmt.Errorf("Scan Sent Connection Request : %w", err)
		}
		connection_requests = append(connection_requests, *connection)
	}

	if err := results.Err(); err != nil {
		return nil, fmt.Errorf("Iterate connection requests: %w", err)
	}

	return &connection_requests, err
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
	err := repo.db.QueryRow(ctx, query, userID).Scan(&connection_request.ConnectionRequestID, &connection_request.SenderID, &connection_request.ReceiverID, &connection_request.CreatedAt)
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

type rowScanner interface {
	Scan(dest ...any) error
}

func scanConnectionRequest(row rowScanner) (*Connection_Request, error) {
	connection_request := &Connection_Request{}
	err := row.Scan(&connection_request.ConnectionRequestID, &connection_request.SenderID, &connection_request.ReceiverID, &connection_request.CreatedAt)
	return connection_request, err
}
