package repository

import (
	"innotech/internal/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestContractRepository_CRUD(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewContractRepository(sqlxDB)

	rows := sqlmock.NewRows([]string{"id", "title", "description"}).AddRow(1, "Contract", "Desc")
	mock.ExpectQuery(`SELECT \* FROM contracts ORDER BY id`).WillReturnRows(rows)
	all, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, all, 1)

	row := sqlmock.NewRows([]string{"id", "title", "description"}).AddRow(1, "Contract", "Desc")
	mock.ExpectQuery(`SELECT \* FROM contracts WHERE id=\$1`).WithArgs(1).WillReturnRows(row)
	got, err := repo.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "Contract", got.Title)

	c := &models.Contract{Title: "New", Description: "D"}
	mock.ExpectQuery(`INSERT INTO contracts \(title, description\) VALUES \(\$1, \$2\) RETURNING id`).
		WithArgs("New", "D").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	err = repo.Create(c)
	assert.NoError(t, err)

	mock.ExpectExec(`UPDATE contracts SET title=\$1, description=\$2 WHERE id=\$3`).
		WithArgs("New", "D", 1).WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.Update(c)
	assert.NoError(t, err)

	mock.ExpectExec(`DELETE FROM contracts WHERE id=\$1`).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.Delete(1)
	assert.NoError(t, err)
}
