package projects

type ProjectCreateRequest struct {
	WorkspaceID  string  `json:"name"`
	Name         string  `json:"name" binding:"required,min=3,max=100"`
	Description  *string `json:"description"`
	ProjectColor *string `json:"project_color"`
	DueDate      string  `json:"due_date"`
}

type Project struct {
	ProjectID    string `json:"project_id"`
	WorkspaceID  string `json:"workspace_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	ProjectColor string `json:"project_color"`
	DueDate      string `json:"due_date"`
	Created_At   string `json:"created_at"`
	Updated_At   string `json:"updated_at"`
}
