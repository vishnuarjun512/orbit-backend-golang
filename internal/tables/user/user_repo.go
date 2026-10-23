package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Begin create user: %w", err)
	}

	userQuery := `
		INSERT INTO users 
		(public_id, email, password_hash) 
		VALUES ('ORB-' || UPPER(ENCODE(gen_random_bytes(8), 'hex')), $1, $2) 
		RETURNING user_id, public_id
	`

	defer tx.Rollback(ctx)
	err = tx.QueryRow(ctx, userQuery, user.Email, user.PasswordHash).Scan(&user.UserID, &user.PublicID)

	if err != nil {
		// PostgreSQL error code 23505 = unique violation.
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrEmailAlreadyExists
		}
		return fmt.Errorf("Create user: %w", err)
	}

	name := user.Email
	if at := strings.IndexByte(name, '@'); at >= 0 {
		name = name[:at]
	}

	userProfileQuery := `
		INSERT INTO user_profiles 
		(user_id, display_name) 
		VALUES ($1 ,$2)
	`

	if _, err := tx.Exec(ctx, userProfileQuery, user.UserID, name); err != nil {
		return fmt.Errorf("Create user profile: %w", err)
	}

	userPreferenceQuery := `
		INSERT INTO user_preferences 
		(user_id) 
		VALUES ($1)
	`

	if _, err := tx.Exec(ctx, userPreferenceQuery, user.UserID); err != nil {
		return fmt.Errorf("create user preferences: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit create user: %w", err)
	}

	return nil
}

// GetByEmail retrieves a user using their email address.
func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT
			user_id,
			public_id,
			email,
			password_hash,
			verified,
			status,
			created_at
		FROM users
		WHERE LOWER(email) = LOWER($1)
		LIMIT 1
	`

	user := &User{}

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.UserID,
		&user.PublicID,
		&user.Email,
		&user.PasswordHash,
		&user.Verified,
		&user.Status,
		&user.CreatedAt,
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
			public_id,
			email,
			password_hash,
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
		&user.PublicID,
		&user.Email,
		&user.PasswordHash,
		&user.Verified,
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
