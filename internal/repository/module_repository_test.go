package repository

import (
	"innotech/internal/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestModuleRepository_CRUD(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewModuleRepository(sqlxDB)

	// GetAll
	rows := sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(1, "Module", "Desc")
	mock.ExpectQuery(`SELECT \* FROM modules ORDER BY id`).WillReturnRows(rows)
	all, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, all, 1)

	// GetByID
	row := sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(1, "Module", "Desc")
	mock.ExpectQuery(`SELECT \* FROM modules WHERE id=\$1`).WithArgs(1).WillReturnRows(row)
	got, err := repo.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "Module", got.Name)

	// Create
	p := &models.Module{Name: "New", Description: "D"}
	mock.ExpectQuery(`INSERT INTO modules \(name, description\) VALUES \(\$1, \$2\) RETURNING id`).
		WithArgs("New", "D").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	err = repo.Create(p)
	assert.NoError(t, err)

	// Update
	mock.ExpectExec(`UPDATE modules SET name=\$1, description=\$2 WHERE id=\$3`).
		WithArgs("New", "D", 1).WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.Update(p)
	assert.NoError(t, err)

	// Delete
	mock.ExpectExec(`DELETE FROM modules WHERE id=\$1`).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.Delete(1)
	assert.NoError(t, err)
}
