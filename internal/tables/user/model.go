package user

import "time"

type User struct {
	UserID       string    `json:"user_id"`
	PublicID     string    `json:"public_id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Verified     bool      `json:"verified"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
