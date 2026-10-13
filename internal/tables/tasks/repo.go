package tasks

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{db: db}
}

func (repo *TaskRepository) CreateTaskRepo(ctx context.Context, workspaceID, projectID, userID, parentTaskID, priority, status, title, description string, start_date, due_date, completed_at string) (string, error) {
	query := `
		INSERT INTO tasks
		(workspace_id, project_id, created_by, parent_task_id, priority, status, title, description, start_date, due_date, completed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);
		RETURNING task_id;
	`

	var task_id string
	err := repo.db.QueryRow(ctx, query, workspaceID, projectID, userID, parentTaskID, priority, status, title, description, start_date, due_date, completed_at).Scan(&task_id)
	if err != nil {
		return "", fmt.Errorf("Creation Task Error: %w", err)
	}
	return task_id, err
}

func (repo *TaskRepository) GetTaskByID(ctx context.Context, task_id string) (*Task, error) {
	query := `
		SELECT
			task_id, workspace_id, project_id, created_by, parent_task_id, priority, status, title, description, start_date, due_date, completed_at, created_at, updated_at
		FROM tasks
		WHERE task_id = $1
	`

	task := &Task{}

	err := repo.db.QueryRow(ctx, query, task_id).Scan(&task.TaskID, &task.ProjectID, &task.WorkspaceID, &task.Created_By, &task.Title, &task.Description, &task.Status, &task.Priority, &task.StartDate, &task.DueDate, &task.CompletedAt, &task.Created_At, &task.Updated_At)
	if err != nil {
		return nil, fmt.Errorf("Get Task Error: %w", err)
	}
	return task, err
}

func (repo *TaskRepository) GetTasksByProjectID(ctx context.Context, project_id string) (*[]Task, error) {
	query := `
		SELECT
			task_id, workspace_id, project_id, created_by, parent_task_id, priority, status, title, description, start_date, due_date, completed_at, created_at, updated_at
		FROM tasks
		WHERE project_id = $1
	`

	tasks := make([]Task, 0)
	rows, err := repo.db.Query(ctx, query, project_id)
	if err != nil {
		return nil, fmt.Errorf("List Tasks with Project ID: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		task, err := scanTasks(rows)
		if err != nil {
			return nil, fmt.Errorf("Scan Task: %w", err)
		}
		tasks = append(tasks, *task)
	}

	return &tasks, err
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanTasks(row rowScanner) (*Task, error) {
	task := &Task{}
	err := row.Scan(&task.TaskID, &task.ProjectID, &task.WorkspaceID, &task.Created_By, &task.Title, &task.Description, &task.StartDate, &task.DueDate, &task.CompletedAt, &task.Created_At, &task.Updated_At)
	return task, err
}

func (r *TaskRepository) DeleteWorkspace(ctx context.Context, id, userID string) error {
	const query = `DELETE FROM tasks WHERE taks_id = $1 AND created_by = $2`
	result, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("Delete Task Error: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrTaskDeletion
	}
	return nil
}
