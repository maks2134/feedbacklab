package contract

import (
	"github.com/jmoiron/sqlx"
)

type ContractRepository struct {
	db *sqlx.DB
}

func NewContractRepository(db *sqlx.DB) *ContractRepository {
	return &ContractRepository{db: db}
}

func (r *ContractRepository) GetAll() ([]Contract, error) {
	var items []Contract
	err := r.db.Select(&items, "SELECT * FROM contracts ORDER BY id")
	return items, err
}

func (r *ContractRepository) GetByID(id int) (*Contract, error) {
	var item Contract
	err := r.db.Get(&item, "SELECT * FROM contracts WHERE id=$1", id)
	return &item, err
}

func (r *ContractRepository) Create(c *Contract) error {
	return r.db.QueryRow(
		`INSERT INTO contracts (project_id, client_name, start_date, end_date, description)
         VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		c.ProjectID, c.ClientName, c.StartDate, c.EndDate, c.Description,
	).Scan(&c.ID)
}

func (r *ContractRepository) Update(c *Contract) error {
	_, err := r.db.Exec(
		`UPDATE contracts SET client_name=$1, start_date=$2, end_date=$3, description=$4 WHERE id=$5`,
		c.ClientName, c.StartDate, c.EndDate, c.Description, c.ID,
	)
	return err
}

func (r *ContractRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM contracts WHERE id=$1`, id)
	return err
}
