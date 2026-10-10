package workspacemembers

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkspaceMemberRepository struct {
	db *pgxpool.Pool
}

func NewWorkspaceMemberRepository(db *pgxpool.Pool) *WorkspaceMemberRepository {
	return &WorkspaceMemberRepository{db: db}
}

func (r *WorkspaceMemberRepository) CreateWorkspaceMemberRepository(ctx context.Context, workspace_id string, user_id string, role string, status string) {
	memberQuery := `
		INSERT INTO workspace_members (workspace_id, user_id, role, status)
		VALUES ($1, $2, $3, $4)
		RETURNING workspace_member_id;
	`

	r.db.Query(ctx, memberQuery, workspace_id, user_id, role, status)
}
