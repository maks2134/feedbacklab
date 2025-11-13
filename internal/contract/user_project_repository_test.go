package repository

import (
	"innotech/internal/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestUserProjectRepository_CRUD(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewUserProjectRepository(sqlxDB)

	rows := sqlmock.NewRows([]string{"id", "user_id", "project_id"}).AddRow(1, "uuid-1", 2)
	mock.ExpectQuery(`SELECT \* FROM user_projects ORDER BY id`).WillReturnRows(rows)
	all, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, all, 1)

	row := sqlmock.NewRows([]string{"id", "user_id", "project_id"}).AddRow(1, "uuid-1", 2)
	mock.ExpectQuery(`SELECT \* FROM user_projects WHERE id=\$1`).WithArgs(1).WillReturnRows(row)
	got, err := repo.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "uuid-1", got.UserID)

	up := &models.UserProject{UserID: "uuid-1", ProjectID: 2}
	mock.ExpectQuery(`INSERT INTO user_projects \(user_id, project_id\) VALUES \(\$1, \$2\) RETURNING id`).
		WithArgs("uuid-1", 2).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	err = repo.Create(up)
	assert.NoError(t, err)

	mock.ExpectExec(`UPDATE user_projects SET user_id=\$1, project_id=\$2 WHERE id=\$3`).
		WithArgs("uuid-1", 2, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.Update(up)
	assert.NoError(t, err)

	mock.ExpectExec(`DELETE FROM user_projects WHERE id=\$1`).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	err = repo.Delete(1)
	assert.NoError(t, err)
}
