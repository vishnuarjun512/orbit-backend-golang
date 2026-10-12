package projects

import "time"

type ProjectCreateRequest struct {
	WorkspaceID  string  `json:"workspace_id"`
	Name         string  `json:"name" binding:"required,min=3,max=100"`
	Description  *string `json:"description"`
	ProjectColor *string `json:"project_color"`
	DueDate      *string `json:"due_date"`
}

type Project struct {
	ProjectID    string    `json:"project_id"`
	WorkspaceID  string    `json:"workspace_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	ProjectColor string    `json:"project_color"`
	DueDate      time.Time `json:"due_date"`
	Created_By   string    `json:"created_by"`
	Created_At   time.Time `json:"created_at"`
	Updated_At   time.Time `json:"updated_at"`
}
