package workspace

import (
	"context"
	"errors"
)

type WorkSpaceService struct {
	repository *WorkspaceRepository
}

func NewWorkSpaceService(repo *WorkspaceRepository) *WorkSpaceService {
	return &WorkSpaceService{
		repository: repo,
	}
}

var (
	ErrInvalidInputs = errors.New("Invalid Inputs")
)

func (s *WorkSpaceService) createWorkspaceService(ctx context.Context, userID string, req WorkSpaceCreateRequest, slug string) {

	s.repository.CreateWorkspaceRepository(ctx, req.Name, slug, userID)
}

func (s *WorkSpaceService) GetWorkSpacesService(ctx context.Context, userID string, limit int, page int) (*[]WorkSpace, error) {
	offset := (page - 1) * limit

	return s.repository.GetWorkSpaces(ctx, userID, limit, offset)
}
