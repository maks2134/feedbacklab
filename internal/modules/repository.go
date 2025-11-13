package modules

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Create(ctx context.Context, m *Module) error
	GetByID(ctx context.Context, id int) (*Module, error)
	GetAll(ctx context.Context) ([]Module, error)
	Update(ctx context.Context, m *Module) error
	Delete(ctx context.Context, id int) error
}

type moduleRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &moduleRepository{db: db}
}

func (r *moduleRepository) Create(ctx context.Context, m *Module) error {
	query := `
		INSERT INTO modules (project_id, name, description, responsible_user_id)
		VALUES (:project_id, :name, :description, :responsible_user_id)
		RETURNING id, date_created, date_updated
	`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	return stmt.GetContext(ctx, m, m)
}

func (r *moduleRepository) GetByID(ctx context.Context, id int) (*Module, error) {
	var m Module
	err := r.db.GetContext(ctx, &m, "SELECT * FROM modules WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *moduleRepository) GetAll(ctx context.Context) ([]Module, error) {
	var modules []Module
	err := r.db.SelectContext(ctx, &modules, "SELECT * FROM modules ORDER BY date_created DESC")
	return modules, err
}

func (r *moduleRepository) Update(ctx context.Context, m *Module) error {
	query := `
		UPDATE modules
		SET name=:name, description=:description, responsible_user_id=:responsible_user_id
		WHERE id=:id
		RETURNING date_updated
	`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	return stmt.GetContext(ctx, &m.DateUpdated, m)
}

func (r *moduleRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM modules WHERE id=$1", id)
	return err
}
