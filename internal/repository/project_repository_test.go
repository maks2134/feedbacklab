package repository

import (
	"innotech/internal/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestProjectRepository_GetAll(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := NewProjectRepository(sqlxDB)
	rows := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow(1, "Test", "Desc")

	mock.ExpectQuery(`SELECT \* FROM projects ORDER BY id`).WillReturnRows(rows)

	items, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, "Test", items[0].Name)
}

func TestProjectRepository_GetByID(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := NewProjectRepository(sqlxDB)
	row := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow(1, "Test", "Desc")

	mock.ExpectQuery(`SELECT \* FROM projects WHERE id=\$1`).
		WithArgs(1).
		WillReturnRows(row)

	item, err := repo.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "Test", item.Name)
}

func TestProjectRepository_Create(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := NewProjectRepository(sqlxDB)
	p := &models.Project{Name: "Test", Description: "Desc"}

	mock.ExpectQuery(`INSERT INTO projects \(name, description\) VALUES \(\$1, \$2\) RETURNING id`).
		WithArgs("Test", "Desc").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err := repo.Create(p)
	assert.NoError(t, err)
	assert.Equal(t, 1, p.ID)
}

func TestProjectRepository_Update(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := NewProjectRepository(sqlxDB)
	p := &models.Project{ID: 1, Name: "New", Description: "Updated"}

	mock.ExpectExec(`UPDATE projects SET name=\$1, description=\$2 WHERE id=\$3`).
		WithArgs("New", "Updated", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Update(p)
	assert.NoError(t, err)
}

func TestProjectRepository_Delete(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := NewProjectRepository(sqlxDB)
	mock.ExpectExec(`DELETE FROM projects WHERE id=\$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Delete(1)
	assert.NoError(t, err)
}
