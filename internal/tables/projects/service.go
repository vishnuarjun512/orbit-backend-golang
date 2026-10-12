package projects

import (
	"context"
	"fmt"
	"time"
)

type ProjectService struct {
	repo *ProjectRepository
}

func NewProjectService(repo *ProjectRepository) *ProjectService {
	return &ProjectService{
		repo: repo,
	}
}

func (s *ProjectService) CreateProjectService(ctx context.Context, userID string, req ProjectCreateRequest) (*Project, error) {
	var dueDate *time.Time

	if req.DueDate != nil && *req.DueDate != "" {
		parsedDate, err := time.Parse("2006-01-02", *req.DueDate)
		if err != nil {
			return nil, fmt.Errorf("invalid due date: expected YYYY-MM-DD")
		}

		dueDate = &parsedDate
	}

	workspaceID, err := s.repo.CreateProjectProjectRepo(ctx, userID, req.WorkspaceID, req.Name, *req.Description, *req.ProjectColor, dueDate)
	if err != nil {
		return nil, err
	}

	return s.repo.GetProjectByIDRepo(ctx, workspaceID)
}

func (s *ProjectService) GetProjectByIDService(ctx context.Context, projectID string) (*Project, error) {
	return s.repo.GetProjectByIDRepo(ctx, projectID)
}

func (s *ProjectService) GetProjectsByWorkspaceIDService(ctx context.Context, workspaceID string, page, limit int) (*[]Project, error) {
	offset := (page - 1) * limit
	return s.repo.GetProjectsByWorkspaceIDRepo(ctx, workspaceID, page, offset)
}

func (s *ProjectService) DeleteProjectByIDService(ctx context.Context, projectID string) error {
	return s.repo.DeleteProjectByIDRepo(ctx, projectID)
}
