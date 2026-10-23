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
	User         SafeUser `json:"user"`
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
}

type SafeUser struct {
	UserID    string    `json:"user_id"`
	PublicID  string    `json:"public_id"`
	Email     string    `json:"email"`
	Verified  bool      `json:"verified"`
	CreatedAt time.Time `json:"created_at"`
}

// UserProfile contains user-editable identity details. Authentication secrets
// and account credentials belong to User and are never part of this response.
type UserProfile struct {
	UserID        string    `json:"user_id"`
	PublicID      string    `json:"public_id"`
	DisplayName   string    `json:"display_name"`
	RoleTitle     *string   `json:"role_title"`
	Pronouns      *string   `json:"pronouns"`
	Bio           *string   `json:"bio"`
	AvatarURL     *string   `json:"avatar_url"`
	Availability  string    `json:"availability"`
	StatusMessage *string   `json:"status_message"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// ProfileUpdateRequest is a partial update: fields omitted by the client stay
// unchanged. Optional text fields currently cannot be cleared with JSON null.
type ProfileUpdateRequest struct {
	DisplayName   *string `json:"display_name"`
	RoleTitle     *string `json:"role_title"`
	Pronouns      *string `json:"pronouns"`
	Bio           *string `json:"bio"`
	AvatarURL     *string `json:"avatar_url"`
	Availability  *string `json:"availability"`
	StatusMessage *string `json:"status_message"`
}

// UserPreferences stores one user's cross-workspace defaults for the settings
// screen. ActiveWorkspaceID is accepted only for active workspace members.
type UserPreferences struct {
	UserID             string    `json:"user_id"`
	ActiveWorkspaceID  *string   `json:"active_workspace_id"`
	EmailNotifications bool      `json:"email_notifications"`
	TaskUpdates        bool      `json:"task_updates"`
	Mentions           bool      `json:"mentions"`
	WeeklyDigest       bool      `json:"weekly_digest"`
	Timezone           string    `json:"timezone"`
	WeekStartsOn       int16     `json:"week_starts_on"`
	Theme              *string   `json:"theme"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// PreferencesUpdateRequest is a partial user settings update. A blank
// ActiveWorkspaceID explicitly clears the selected workspace.
type PreferencesUpdateRequest struct {
	ActiveWorkspaceID  *string `json:"active_workspace_id"`
	EmailNotifications *bool   `json:"email_notifications"`
	TaskUpdates        *bool   `json:"task_updates"`
	Mentions           *bool   `json:"mentions"`
	WeeklyDigest       *bool   `json:"weekly_digest"`
	Timezone           *string `json:"timezone"`
	WeekStartsOn       *int16  `json:"week_starts_on"`
	Theme              *string `json:"theme"`
}
