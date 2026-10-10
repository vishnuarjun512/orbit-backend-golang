package workspace

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode"
)

var (
	ErrInvalidWorkspaceInput = errors.New("invalid workspace input")
	ErrWorkspaceNotFound     = errors.New("workspace not found")
	ErrWorkspaceSlugConflict = errors.New("workspace slug already exists")
)

type WorkSpaceService struct {
	repository *WorkspaceRepository
}

func NewWorkSpaceService(repo *WorkspaceRepository) *WorkSpaceService {
	return &WorkSpaceService{repository: repo}
}

func (s *WorkSpaceService) CreateWorkspace(
	ctx context.Context,
	userID string,
	req WorkSpaceCreateRequest,
) (*WorkSpace, error) {
	name := strings.TrimSpace(req.Name)
	if len([]rune(name)) < 3 || len([]rune(name)) > 100 {
		return nil, fmt.Errorf("%w: name must be between 3 and 100 characters", ErrInvalidWorkspaceInput)
	}

	slug := slugify(name)
	if slug == "" || len([]rune(slug)) > 100 {
		return nil, fmt.Errorf("%w: name must produce a slug containing 1 to 100 characters", ErrInvalidWorkspaceInput)
	}

	return s.repository.CreateWorkspaceRepository(ctx, name, slug, req.Description, userID)
}

func (s *WorkSpaceService) GetWorkspace(ctx context.Context, id, userID string) (*WorkSpace, error) {
	return s.repository.GetWorkspaceByID(ctx, id, userID)
}

func (s *WorkSpaceService) GetWorkSpacesService(
	ctx context.Context,
	userID string,
	page, limit int,
) ([]WorkSpace, error) {
	offset := (page - 1) * limit
	return s.repository.GetWorkSpaces(ctx, userID, limit, offset)
}

func (s *WorkSpaceService) UpdateWorkspace(
	ctx context.Context,
	id, userID string,
	req WorkSpaceUpdateRequest,
) (*WorkSpace, error) {
	if req.Name == nil && req.Slug == nil && req.Description == nil {
		return nil, fmt.Errorf("%w: provide at least one field to update", ErrInvalidWorkspaceInput)
	}
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if len([]rune(name)) < 3 || len([]rune(name)) > 100 {
			return nil, fmt.Errorf("%w: name must be between 3 and 100 characters", ErrInvalidWorkspaceInput)
		}
		req.Name = &name
	}
	if req.Slug != nil {
		slug := slugify(*req.Slug)
		if slug == "" || len([]rune(slug)) > 100 {
			return nil, fmt.Errorf("%w: slug must contain 1 to 100 characters", ErrInvalidWorkspaceInput)
		}
		req.Slug = &slug
	}
	return s.repository.UpdateWorkspace(ctx, id, userID, req)
}

func (s *WorkSpaceService) DeleteWorkspace(ctx context.Context, id, userID string) error {
	return s.repository.DeleteWorkspace(ctx, id, userID)
}

func slugify(value string) string {
	var slug strings.Builder
	hyphen := false
	for _, char := range strings.ToLower(strings.TrimSpace(value)) {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			slug.WriteRune(char)
			hyphen = false
			continue
		}
		if slug.Len() > 0 && !hyphen {
			slug.WriteByte('-')
			hyphen = true
		}
	}
	return strings.Trim(slug.String(), "-")
}
