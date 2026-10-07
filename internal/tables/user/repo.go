package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
func (r *Repository) CreateUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (
			email,
			password_hash
		)
		VALUES ($1, $2)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		user.Email,
		user.PasswordHash,
	)

	if err != nil {
		// PostgreSQL error code 23505 = unique violation.
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrEmailAlreadyExists
		}

		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

// GetByEmail retrieves a user using their email address.
func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("Get user by Email: %w", err)
	}

	return user, nil
}

// GetByUsername retrieves a user using their username.
func (r *Repository) GetUserByID(ctx context.Context, id string) (*User, error) {
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
		WHERE id = $1
		LIMIT 1
	`

	user := &User{}

	err := r.db.QueryRow(ctx, query, id).Scan(
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("Get user by ID: %w", err)
	}

	return user, nil
}
