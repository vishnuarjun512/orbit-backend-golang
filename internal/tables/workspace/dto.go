package workspace

import "time"

type WorkSpaceCreateRequest struct {
	Name        string  `json:"name" binding:"required,min=3,max=100"`
	Slug        string  `json:"slug" binding:"omitempty,max=100"`
	Description *string `json:"description"`
}

type WorkSpaceUpdateRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=3,max=100"`
	Slug        *string `json:"slug" binding:"omitempty,max=100"`
	Description *string `json:"description"`
}

type WorkSpace struct {
	WorkSpaceID string    `json:"workspace_id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description *string   `json:"description,omitempty"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
