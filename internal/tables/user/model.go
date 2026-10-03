package user

import "time"

type User struct {
    UserID         string     `json:"user_id"`
    Email          string     `json:"email"`
    Username       string     `json:"username"`
    FullName       string     `json:"full_name"`
    PasswordHash   string     `json:"-"`
    AvatarURL      *string    `json:"avatar_url,omitempty"`
    EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
    Status         string     `json:"status"`
    CreatedAt      time.Time  `json:"created_at"`
    UpdatedAt      time.Time  `json:"updated_at"`
}