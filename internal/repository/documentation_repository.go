package repository

import (
	"innotech/internal/models"
	"github.com/jmoiron/sqlx"
)

type DocumentationRepository struct {
	db *sqlx.DB
}

func NewDocumentationRepository(db *sqlx.DB) *DocumentationRepository {
	return &DocumentationRepository{db: db}
}

func (r *DocumentationRepository) GetAll() ([]models.Documentation, error) {
	var items []models.Documentation
	err := r.db.Select(&items, "SELECT * FROM documentations ORDER BY id")
	return items, err
}

func (r *DocumentationRepository) GetByID(id int) (*models.Documentation, error) {
	var item models.Documentation
	err := r.db.Get(&item, "SELECT * FROM documentations WHERE id=$1", id)
	return &item, err
}

func (r *DocumentationRepository) Create(d *models.Documentation) error {
	return r.db.QueryRow(
		`INSERT INTO documentations (project_id, title, content, file_url)
         VALUES ($1, $2, $3, $4) RETURNING id`,
		d.ProjectID, d.Title, d.Content, d.FileURL,
	).Scan(&d.ID)
}

func (r *DocumentationRepository) Update(d *models.Documentation) error {
	_, err := r.db.Exec(
		`UPDATE documentations SET title=$1, content=$2, file_url=$3 WHERE id=$4`,
		d.Title, d.Content, d.FileURL, d.ID,
	)
	return err
}

func (r *DocumentationRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM documentations WHERE id=$1`, id)
	return err
}
