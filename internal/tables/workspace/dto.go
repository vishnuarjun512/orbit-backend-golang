package workspace

import "time"

// DTO represents the data required to register a new user.
type WorkSpaceCreateRequest struct {
	Name string `json:"name" binding:"required,min=3,max=15"`
}

// DTO represents the data required to register a new user.
type WorkSpace struct {
	WorkSpaceID string    `json:"workspace_id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
