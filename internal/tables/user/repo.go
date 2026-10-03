package user

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

// Create inserts a new user into PostgreSQL.
func (r *Repository) Create(ctx context.Context, user *User) error {
	query := `
        INSERT INTO users (
            email,
            username,
            full_name,
            password_hash
        )
        VALUES ($1, $2, $3, $4)
        RETURNING
            user_id,
            avatar_url,
            status,
            created_at,
            updated_at
    `

	err := r.db.QueryRow(
		ctx,
		query,
		user.Email,
		user.Username,
		user.FullName,
		user.PasswordHash,
	).Scan(
		&user.UserID,
		&user.AvatarURL,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

// GetByEmail retrieves a user using their email address.
func (r *Repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
        SELECT
            user_id,
            email,
            username,
            full_name,
            password_hash,
            avatar_url,
            email_verified_at,
            status,
            created_at,
            updated_at
        FROM users
        WHERE LOWER(email) = LOWER($1)
        LIMIT 1
    `

	user := &User{}

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.UserID,
		&user.Email,
		&user.Username,
		&user.FullName,
		&user.PasswordHash,
		&user.AvatarURL,
		&user.EmailVerifiedAt,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return user, nil
}

// GetByUsername retrieves a user using their username.
func (r *Repository) GetByUsername(ctx context.Context, username string) (*User, error) {
	query := `
        SELECT
            user_id,
            email,
            username,
            full_name,
            password_hash,
            avatar_url,
            email_verified_at,
            status,
            created_at,
            updated_at
        FROM users
        WHERE LOWER(username) = LOWER($1)
        LIMIT 1
    `

	user := &User{}

	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.UserID,
		&user.Email,
		&user.Username,
		&user.FullName,
		&user.PasswordHash,
		&user.AvatarURL,
		&user.EmailVerifiedAt,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}

		return nil, fmt.Errorf("get user by username: %w", err)
	}

	return user, nil
}
