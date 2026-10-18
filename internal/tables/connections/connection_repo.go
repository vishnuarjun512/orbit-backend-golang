package connections

import (
	"context"
	"fmt"
)

func (repo *ConnectionRepository) CreateNewConnection(ctx context.Context, senderID, receiverID string) (string, error) {
	query1 := `
		INSERT INTO connections
		(user_a_id, user_b_id)
		VALUES (
			LEAST($1::uuid, $2::uuid),
			GREATEST($1::uuid, $2::uuid)
		)
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

func (repo *ConnectionRepository) GetConnectionByID(ctx context.Context, connection_request_id string) (*Connection, error) {
	query := `
		SELECT 
			connection_id,
			user_a_id,
			user_b_id,
			created_at
		FROM connections
		WHERE connection_id=$1
	`

	connection := &Connection{}
	err := repo.db.QueryRow(ctx, query, connection_request_id).Scan(&connection.ConnectionID, &connection.UserAID, &connection.UserBID, &connection.CreatedAt)

	return connection, err
}

func (repo *ConnectionRepository) GetConnection(ctx context.Context, sender_id, receiver_id string) (*Connection, error) {
	query := `
		SELECT 
			connection_id,
			user_a_id,
			user_b_id,
			created_at
		FROM connections
		WHERE user_a_id = LEAST($1::uuid, $2::uuid)
  			AND user_b_id = GREATEST($1::uuid, $2::uuid);
	`

	connection_request := &Connection{}
	err := repo.db.QueryRow(ctx, query, sender_id, receiver_id).Scan(&connection_request.ConnectionID, &connection_request.UserAID, &connection_request.UserBID, &connection_request.CreatedAt)
	return connection_request, err
}

func (repo *ConnectionRepository) GetAllConnections(
	ctx context.Context,
	userID string,
) (*[]ConnectionResponse, error) {
	query := `
        SELECT 
			c.connection_id,
			u.user_id,
			u.username,
			u.email,
			u.avatar_url,
			c.created_at
        FROM connections c
		JOIN users u
		ON u.user_id = CASE
                WHEN c.user_a_id = $1::uuid THEN c.user_b_id
                ELSE c.user_a_id
            END
		WHERE c.user_a_id = $1::uuid
           OR c.user_b_id = $1::uuid
        ORDER BY created_at DESC;
    `

	rows, err := repo.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("Get connections: %w", err)
	}
	defer rows.Close()

	connections := make([]ConnectionResponse, 0)

	for rows.Next() {
		var connection ConnectionResponse

		if err := rows.Scan(
			&connection.ConnectionID,
			&connection.ConnectedUser.UserID,
			&connection.ConnectedUser.Username,
			&connection.ConnectedUser.Email,
			&connection.ConnectedUser.AvatarURL,
			&connection.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("Scan connection: %w", err)
		}

		connections = append(connections, connection)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate connections: %w", err)
	}

	return &connections, nil
}

func (repo *ConnectionRepository) DeleteConnection(ctx context.Context, sender_id, receiver_id string) error {
	query := `
		DELETE FROM connections
		WHERE user_a_id = LEAST($1::uuid, $2::uuid)
  			AND user_b_id = GREATEST($1::uuid, $2::uuid);
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
