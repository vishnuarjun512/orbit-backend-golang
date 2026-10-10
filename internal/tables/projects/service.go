package projects

import "context"

type ProjectService struct {
	repo *ProjectRepository
}

func NewProjectService(repo *ProjectRepository) *ProjectService {
	return &ProjectService{
		repo: repo,
	}
}

func (s *ProjectService) CreateProjectService(ctx context.Context, userID string, req ProjectCreateRequest) (*Project, error) {
	workspaceID, err := s.repo.CreateProjectProjectRepo(ctx, userID, req.WorkspaceID, req.Name, *req.Description, *req.ProjectColor, req.DueDate)
	if err != nil {
		return nil, err
	}

	return s.repo.GetProjectByIDRepo(ctx, workspaceID)
}

func (s *ProjectService) GetProjectByIDService(ctx context.Context, workspaceID string) (*Project, error) {
	return s.repo.GetProjectByIDRepo(ctx, workspaceID)
}
