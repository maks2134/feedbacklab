package repository

import (
	"innotech/internal/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestDocumentationRepository_CRUD(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewDocumentationRepository(sqlxDB)

	rows := sqlmock.NewRows([]string{"id", "title", "content"}).AddRow(1, "Doc", "Text")
	mock.ExpectQuery(`SELECT \* FROM documentations ORDER BY id`).WillReturnRows(rows)
	all, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, all, 1)

	row := sqlmock.NewRows([]string{"id", "title", "content"}).AddRow(1, "Doc", "Text")
	mock.ExpectQuery(`SELECT \* FROM documentations WHERE id=\$1`).WithArgs(1).WillReturnRows(row)
	got, err := repo.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "Doc", got.Title)

	d := &models.Documentation{Title: "New", Content: "Body"}
	mock.ExpectQuery(`INSERT INTO documentations \(title, content\) VALUES \(\$1, \$2\) RETURNING id`).
		WithArgs("New", "Body").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	err = repo.Create(d)
	assert.NoError(t, err)

	mock.ExpectExec(`UPDATE documentations SET title=\$1, content=\$2 WHERE id=\$3`).
		WithArgs("New", "Body", 1).WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.Update(d)
	assert.NoError(t, err)

	mock.ExpectExec(`DELETE FROM documentations WHERE id=\$1`).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.Delete(1)
	assert.NoError(t, err)
}
