package projects

import (
	"context"
	"innotech/internal/storage/postgres"

	"github.com/jmoiron/sqlx"
)

// Repository defines the interface for project data access operations.
type Repository interface {
	Create(ctx context.Context, p *postgres.Project) error
	GetByID(ctx context.Context, id int) (*postgres.Project, error)
	GetAll(ctx context.Context) ([]postgres.Project, error)
	Update(ctx context.Context, p *postgres.Project) error
	Delete(ctx context.Context, id int) error
}

type projectRepository struct {
	db *sqlx.DB
}

// NewRepository creates a new Repository instance.
func NewRepository(db *sqlx.DB) Repository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(ctx context.Context, p *postgres.Project) error {
	query := `
		INSERT INTO projects (name, description, gitlab_project_id, mattermost_team)
		VALUES (:name, :description, :gitlab_project_id, :mattermost_team)
		RETURNING id, date_created, date_updated
	`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	return stmt.GetContext(ctx, p, p)
}

func (r *projectRepository) GetByID(ctx context.Context, id int) (*postgres.Project, error) {
	var p postgres.Project
	err := r.db.GetContext(ctx, &p, "SELECT * FROM projects WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *projectRepository) GetAll(ctx context.Context) ([]postgres.Project, error) {
	var ps []postgres.Project
	err := r.db.SelectContext(ctx, &ps, "SELECT * FROM projects ORDER BY date_created DESC")
	return ps, err
}

func (r *projectRepository) Update(ctx context.Context, p *postgres.Project) error {
	query := `
		UPDATE projects
		SET name=:name, description=:description, gitlab_project_id=:gitlab_project_id, mattermost_team=:mattermost_team
		WHERE id=:id
		RETURNING date_updated
	`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	return stmt.GetContext(ctx, &p.DateUpdated, p)
}

func (r *projectRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM projects WHERE id=$1", id)
	return err
}
