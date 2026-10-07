package workspace

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkspaceRepository struct {
	db *pgxpool.Pool
}

func NewWorkSpaceRepository(db *pgxpool.Pool) *WorkspaceRepository {
	return &WorkspaceRepository{
		db: db,
	}
}

func (r *WorkspaceRepository) CreateWorkspaceRepository(ctx context.Context, name string, slug string, created_by string) (string, error) {
	query := `
		INSERT INTO workspaces
		(name, slug, created_by)
		VALUES ($1, $2, $3)	
		RETURNING workspace_id
	`

	var workspace_id string

	err := r.db.QueryRow(ctx, query, name, slug, created_by).Scan(&workspace_id)

	if err != nil {
		return "", err
	}

	return workspace_id, nil
}

func (r *WorkspaceRepository) GetWorkSpaceByID(ctx context.Context, id string) (*WorkSpace, error) {
	query := `
		SELECT
			workspace_id,
			name,
			slug,
			created_by,
			created_at,
			updated_at
		FROM workspace
		WHERE workspace_id = $1
		LIMIT 1
	`

	workspace := WorkSpace{}
	err := r.db.QueryRow(ctx, query, id).Scan(&workspace.WorkSpaceID, &workspace.Name, &workspace.Slug, &workspace.CreatedBy, &workspace.CreatedAt, &workspace.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &workspace, nil
}

func (r *WorkspaceRepository) GetWorkSpaces(ctx context.Context, user_id string, limit int, offset int) (*[]WorkSpace, error) {
	query := `
		SELECT
			workspace_id,
			name,
			slug,
			created_by,
			created_at,
			updated_at
		FROM workspace
		WHERE created_by = $1
		LIMIT $2 OFFSET $3
	`

	result, err := r.db.Query(ctx, query, user_id, limit, offset)
	if err != nil {
		return nil, err
	}

	workspaces := []WorkSpace{}
	for result.Next() {
		var workspace WorkSpace
		if err := result.Scan(&workspace.WorkSpaceID, &workspace.Name, &workspace.Slug, &workspace.CreatedBy, &workspace.CreatedAt, &workspace.UpdatedAt); err != nil {
			return nil, err
		}
		workspaces = append(workspaces, workspace)
	}

	return &workspaces, nil
}
