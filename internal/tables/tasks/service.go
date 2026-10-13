package tasks

import (
	"context"
	"errors"
)

type TaskService struct {
	repository *TaskRepository
}

func NewTaskService(repository *TaskRepository) *TaskService {
	return &TaskService{repository: repository}
}

var (
	ErrInvalidTaskInput = errors.New("Invalid Task input")
	ErrTaskNotFound     = errors.New("Task not found")
	ErrTaskCreation     = errors.New("Task Creation Error")
	ErrTaskSlugConflict = errors.New("Task slug already exists")
	ErrTaskDeletion     = errors.New("Delete Task Failed")
)

func (s *TaskService) CreateTaskService(ctx context.Context, userID string, req *TaskCreateRequest) (*Task, error) {
	task_id, err := s.repository.CreateTaskRepo(ctx, req.WorkspaceID, req.ProjectID, userID, req.ParentTaskID, req.Priority, req.Status, req.Title, *req.Description, *req.StartDate, *req.DueDate, *req.CompletedAt)
	if err != nil {
		return nil, ErrTaskCreation
	}

	task, err := s.repository.GetTaskByID(ctx, task_id)
	return task, err
}
