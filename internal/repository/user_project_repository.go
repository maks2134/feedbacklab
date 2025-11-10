package repository

import (
	"innotech/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserProjectRepository struct {
	db *sqlx.DB
}

func NewUserProjectRepository(db *sqlx.DB) *UserProjectRepository {
	return &UserProjectRepository{db: db}
}

func (r *UserProjectRepository) GetAll() ([]models.UserProject, error) {
	var items []models.UserProject
	err := r.db.Select(&items, "SELECT * FROM user_projects ORDER BY id")
	return items, err
}

func (r *UserProjectRepository) GetByID(id int) (*models.UserProject, error) {
	var item models.UserProject
	err := r.db.Get(&item, "SELECT * FROM user_projects WHERE id=$1", id)
	return &item, err
}

func (r *UserProjectRepository) Create(up *models.UserProject) error {
	return r.db.QueryRow(
		`INSERT INTO user_projects (user_id, project_id, role)
         VALUES ($1, $2, $3) RETURNING id`,
		up.UserID, up.ProjectID, up.Role,
	).Scan(&up.ID)
}

func (r *UserProjectRepository) Update(up *models.UserProject) error {
	_, err := r.db.Exec(
		`UPDATE user_projects SET user_id=$1, project_id=$2, role=$3 WHERE id=$4`,
		up.UserID, up.ProjectID, up.Role, up.ID,
	)
	return err
}

func (r *UserProjectRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM user_projects WHERE id=$1`, id)
	return err
}
