package workspace

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkspaceRepository struct {
	db *pgxpool.Pool
}

func NewWorkSpaceRepository(db *pgxpool.Pool) *WorkspaceRepository {
	return &WorkspaceRepository{db: db}
}

func (r *WorkspaceRepository) CreateWorkspaceRepository(
	ctx context.Context,
	name string,
	slug string,
	description *string,
	createdBy string,
) (*WorkSpace, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin create workspace transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	const createQuery = `
		INSERT INTO workspaces (name, slug, description, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING workspace_id, name, slug, description, created_by, created_at, updated_at
	`

	workspace, err := scanWorkspace(tx.QueryRow(ctx, createQuery, name, slug, description, createdBy))

	if err != nil {
		return nil, translateWorkspaceError("create workspace", err)
	}

	const memberQuery = `
		INSERT INTO workspace_members (workspace_id, user_id, role, status)
		VALUES ($1, $2, 'owner', 'active')
	`

	if _, err := tx.Exec(ctx, memberQuery, workspace.WorkSpaceID, createdBy); err != nil {
		return nil, fmt.Errorf("add workspace owner membership: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit create workspace transaction: %w", err)
	}

	return workspace, nil
}

func (r *WorkspaceRepository) GetWorkspaceByID(ctx context.Context, id, userID string) (*WorkSpace, error) {
	const query = `
		SELECT workspace_id, name, slug, description, created_by, created_at, updated_at
		FROM workspaces
		WHERE workspace_id = $1 AND created_by = $2
	`
	workspace, err := scanWorkspace(r.db.QueryRow(ctx, query, id, userID))
	if err != nil {
		return nil, translateWorkspaceError("get workspace", err)
	}
	return workspace, nil
}

func (r *WorkspaceRepository) GetWorkSpaces(ctx context.Context, userID string, limit, offset int) ([]WorkSpace, error) {
	const query = `
		SELECT workspace_id, name, slug, description, created_by, created_at, updated_at
		FROM workspaces
		WHERE created_by = $1
		ORDER BY created_at DESC, workspace_id
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list workspaces: %w", err)
	}
	defer rows.Close()

	workspaces := make([]WorkSpace, 0)
	for rows.Next() {
		workspace, err := scanWorkspace(rows)
		if err != nil {
			return nil, fmt.Errorf("scan workspace: %w", err)
		}
		workspaces = append(workspaces, *workspace)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate workspaces: %w", err)
	}
	return workspaces, nil
}

func (r *WorkspaceRepository) UpdateWorkspace(
	ctx context.Context,
	id, userID string,
	req WorkSpaceUpdateRequest,
) (*WorkSpace, error) {
	const query = `
		UPDATE workspaces
		SET name = COALESCE($3, name),
		    slug = COALESCE($4, slug),
		    description = COALESCE($5, description),
		    updated_at = NOW()
		WHERE workspace_id = $1 AND created_by = $2
		RETURNING workspace_id, name, slug, description, created_by, created_at, updated_at
	`
	workspace, err := scanWorkspace(r.db.QueryRow(ctx, query, id, userID, req.Name, req.Slug, req.Description))
	if err != nil {
		return nil, translateWorkspaceError("update workspace", err)
	}
	return workspace, nil
}

func (r *WorkspaceRepository) DeleteWorkspace(ctx context.Context, id, userID string) error {
	const query = `DELETE FROM workspaces WHERE workspace_id = $1 AND created_by = $2`
	result, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("delete workspace: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrWorkspaceNotFound
	}
	return nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanWorkspace(row rowScanner) (*WorkSpace, error) {
	workspace := &WorkSpace{}
	err := row.Scan(
		&workspace.WorkSpaceID,
		&workspace.Name,
		&workspace.Slug,
		&workspace.Description,
		&workspace.CreatedBy,
		&workspace.CreatedAt,
		&workspace.UpdatedAt,
	)
	return workspace, err
}

func translateWorkspaceError(operation string, err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrWorkspaceNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ErrWorkspaceSlugConflict
	}
	return fmt.Errorf("%s: %w", operation, err)
}
