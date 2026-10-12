package projects

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectRepository struct {
	db *pgxpool.Pool
}

func NewProjectRepository(db *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (repo *ProjectRepository) CreateProjectProjectRepo(ctx context.Context, userID string, workspaceID string, name string, description string, project_color string, due_date *time.Time) (string, error) {
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
		SELECT
		 	project_id,
            workspace_id,
            created_by,
            name,
            description,
            due_date,
            project_color,
            created_at
		FROM projects
		WHERE project_id = $1
		LIMIT 1;
	`

	project := &Project{}
	err := r.db.QueryRow(ctx, query, project_id).Scan(&project.ProjectID, &project.WorkspaceID, &project.Created_By, &project.Name, &project.Description, &project.DueDate, &project.ProjectColor, &project.Created_At)
	return project, err
}

func (r *ProjectRepository) GetProjectsByWorkspaceIDRepo(ctx context.Context, workspace_id string, limit int, offset int) (*[]Project, error) {
	query := `
		SELECT
		 	project_id,
            workspace_id,
            created_by,
            name,
            description,
            due_date,
            project_color,
            created_at
		FROM projects
		WHERE workspace_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, workspace_id, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("List Projects: %w", err)
	}
	defer rows.Close()

	projects := make([]Project, 0)

	for rows.Next() {
		project, err := scanProjects(rows)
		if err != nil {
			return nil, fmt.Errorf("Scan Project: %w", err)
		}
		projects = append(projects, *project)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Iterate Project: %w", err)
	}

	return &projects, err
}

func (r *ProjectRepository) DeleteProjectByIDRepo(ctx context.Context, project_id string) error {
	query := `
		DELETE FROM projects
		WHERE project_id = $1
	`

	_, err := r.db.Exec(ctx, query, project_id)

	return err
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanProjects(row rowScanner) (*Project, error) {
	project := &Project{}
	err := row.Scan(&project.ProjectID, &project.WorkspaceID, &project.Created_By, &project.Name, &project.Description, &project.DueDate, &project.ProjectColor, &project.Created_At)
	return project, err
}
