package repository

import (
	"innotech/internal/models"
	"github.com/jmoiron/sqlx"
)

type ProjectRepository struct {
	db *sqlx.DB
}

func NewProjectRepository(db *sqlx.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) GetAll() ([]models.Project, error) {
	var items []models.Project
	err := r.db.Select(&items, "SELECT * FROM projects ORDER BY id")
	return items, err
}

func (r *ProjectRepository) GetByID(id int) (*models.Project, error) {
	var item models.Project
	err := r.db.Get(&item, "SELECT * FROM projects WHERE id=$1", id)
	return &item, err
}

func (r *ProjectRepository) Create(p *models.Project) error {
	return r.db.QueryRow(
		`INSERT INTO projects (name, description) VALUES ($1, $2) RETURNING id`,
		p.Name, p.Description,
	).Scan(&p.ID)
}

func (r *ProjectRepository) Update(p *models.Project) error {
	_, err := r.db.Exec(
		`UPDATE projects SET name=$1, description=$2 WHERE id=$3`,
		p.Name, p.Description, p.ID,
	)
	return err
}

func (r *ProjectRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM projects WHERE id=$1`, id)
	return err
}
