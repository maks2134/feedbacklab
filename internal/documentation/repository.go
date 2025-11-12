package documentation

import (
    "github.com/jmoiron/sqlx"
)

type DocumentationRepository struct {
    db *sqlx.DB
}

func NewDocumentationRepository(db *sqlx.DB) *DocumentationRepository {
    return &DocumentationRepository{db: db}
}

func (r *DocumentationRepository) GetAll() ([]Documentation, error) {
    var items []Documentation
    err := r.db.Select(&items, "SELECT * FROM documentations ORDER BY id")
    return items, err
}

func (r *DocumentationRepository) GetByID(id int) (*Documentation, error) {
    var item Documentation
    err := r.db.Get(&item, "SELECT * FROM documentations WHERE id=$1", id)
    return &item, err
}

func (r *DocumentationRepository) Create(d *Documentation) error {
    return r.db.QueryRow(
        `INSERT INTO documentations (project_id, title, content, file_url)
         VALUES ($1, $2, $3, $4) RETURNING id`,
        d.ProjectID, d.Title, d.Content, d.FileURL,
    ).Scan(&d.ID)
}

func (r *DocumentationRepository) Update(d *Documentation) error {
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
package documentation

import (
	"github.com/jmoiron/sqlx"
)

type DocumentationRepository struct {
	db *sqlx.DB
}

func NewDocumentationRepository(db *sqlx.DB) *DocumentationRepository {
	return &DocumentationRepository{db: db}
}

func (r *DocumentationRepository) GetAll() ([]Documentation, error) {
	var items []Documentation
	err := r.db.Select(&items, "SELECT * FROM documentations ORDER BY id")
	return items, err
}

func (r *DocumentationRepository) GetByID(id int) (*Documentation, error) {
	var item Documentation
	err := r.db.Get(&item, "SELECT * FROM documentations WHERE id=$1", id)
	return &item, err
}

func (r *DocumentationRepository) Create(d *Documentation) error {
	return r.db.QueryRow(
		`INSERT INTO documentations (project_id, title, content, file_url)
         VALUES ($1, $2, $3, $4) RETURNING id`,
		d.ProjectID, d.Title, d.Content, d.FileURL,
	).Scan(&d.ID)
}

func (r *DocumentationRepository) Update(d *Documentation) error {
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
