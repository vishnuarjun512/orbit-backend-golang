
package user

import "time"

// RegisterRequest represents the data required to register a new user.
type RegisterRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Username string `json:"username" binding:"required,min=3,max=50"`
    FullName string `json:"full_name" binding:"required,min=2,max=100"`
    Password string `json:"password" binding:"required,min=8"`
}

// RegisterResponse represents the public user data returned after registration.
type RegisterResponse struct {
    UserID    string    `json:"user_id"`
    Email     string    `json:"email"`
    Username  string    `json:"username"`
    FullName  string    `json:"full_name"`
    AvatarURL *string   `json:"avatar_url,omitempty"`
    CreatedAt time.Time `json:"created_at"`
}
