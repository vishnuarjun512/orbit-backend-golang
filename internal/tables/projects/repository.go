package projects

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectRepository struct {
	db *pgxpool.Pool
}

func NewProjectRepository(db *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (repo *ProjectRepository) CreateProjectProjectRepo(ctx context.Context, userID string, workspaceID string, name string, description string, project_color string, due_date string) (string, error) {
	query := `
		INSERT INTO projects (workspace_id, created_by, name, description, project_color, due_date)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING project_id;	
	`

	var project_id string

	err := repo.db.QueryRow(ctx, query, workspaceID, userID, name, description, project_color, due_date).Scan(&project_id)

	return project_id, err
}

func (r *ProjectRepository) GetProjectByIDRepo(ctx context.Context, project_id string) (*Project, error) {
	query := `
		SELECT * FROM projects
		WHERE project_id = $1
		LIMIT 1;
	`

	project := &Project{}
	err := r.db.QueryRow(ctx, query, project_id).Scan(&project.ProjectID, &project.WorkspaceID, &project.Name, &project.Description, &project.DueDate, &project.ProjectColor, &project.Created_At)
	return project, err
}
