package repository

import (
	"innotech/internal/models"
	"github.com/jmoiron/sqlx"
)

type ModuleRepository struct {
	db *sqlx.DB
}

func NewModuleRepository(db *sqlx.DB) *ModuleRepository {
	return &ModuleRepository{db: db}
}

func (r *ModuleRepository) GetAll() ([]models.Module, error) {
	var items []models.Module
	err := r.db.Select(&items, "SELECT * FROM modules ORDER BY id")
	return items, err
}

func (r *ModuleRepository) GetByID(id int) (*models.Module, error) {
	var item models.Module
	err := r.db.Get(&item, "SELECT * FROM modules WHERE id=$1", id)
	return &item, err
}

func (r *ModuleRepository) Create(m *models.Module) error {
	return r.db.QueryRow(
		`INSERT INTO modules (project_id, name, description) VALUES ($1, $2, $3) RETURNING id`,
		m.ProjectID, m.Name, m.Description,
	).Scan(&m.ID)
}

func (r *ModuleRepository) Update(m *models.Module) error {
	_, err := r.db.Exec(
		`UPDATE modules SET name=$1, description=$2 WHERE id=$3`,
		m.Name, m.Description, m.ID,
	)
	return err
}

func (r *ModuleRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM modules WHERE id=$1`, id)
	return err
}
