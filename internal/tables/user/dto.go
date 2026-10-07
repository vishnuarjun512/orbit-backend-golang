package user

import "time"

// RegisterRequest represents the data required to register a new user.
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=3,max=15"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=3,max=15"`
}

// RegisterResponse represents the public user data returned after registration.
type LoginResponse struct {
	User        SafeUser `json:"user"`
	AccessToken string   `json:"access_token"`
}

type SafeUser struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Username  *string   `json:"username"`
	FullName  *string   `json:"full_name"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
