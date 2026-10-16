package home

import (
	"orbit-backend-golang/internal/tables/projects"
	"orbit-backend-golang/internal/tables/tasks"
	"orbit-backend-golang/internal/tables/workspace"
)

type HomeService struct {
	workRepository    *workspace.WorkspaceRepository
	taskRepository    *tasks.TaskRepository
	projectRepository *projects.ProjectRepository
}

func NewHomeService(workspaceRepo *workspace.WorkspaceRepository, taskRepo *tasks.TaskRepository, projectRepo *projects.ProjectRepository) *HomeService {
	return &HomeService{
		workRepository:    workspaceRepo,
		taskRepository:    taskRepo,
		projectRepository: projectRepo,
	}
}
