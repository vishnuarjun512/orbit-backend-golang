package tasks

import "time"

type TaskCreateRequest struct {
	ProjectID    string  `json:"project_id"`
	WorkspaceID  string  `json:"workspace_id"`
	ParentTaskID string  `json:"parent_task_id"`
	Title        string  `json:"title" binding:"required,min=3,max=100"`
	Description  *string `json:"description"`
	Status       string  `json:"status"`
	Priority     string  `json:"priority"`
	StartDate    *string `json:"start_date"`
	DueDate      *string `json:"due_date"`
	CompletedAt  *string `json:"completed_at"`
}

type Task struct {
	TaskID       string    `json:"task_id"`
	ProjectID    string    `json:"project_id"`
	WorkspaceID  string    `json:"workspace_id"`
	Created_By   string    `json:"created_by"`
	ParentTaskID string    `json:"parent_task_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Status       string    `json:"status"`
	Priority     string    `json:"priority"`
	StartDate    time.Time `json:"start_date"`
	DueDate      time.Time `json:"due_date"`
	CompletedAt  time.Time `json:"completed_at"`
	Created_At   time.Time `json:"created_at"`
	Updated_At   time.Time `json:"updated_at"`
}
