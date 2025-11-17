package contract

import (
	"github.com/jmoiron/sqlx"
)

// ContractRepository handles database operations for contracts.
type ContractRepository struct {
	db *sqlx.DB
}

// NewContractRepository creates a new ContractRepository instance.
func NewContractRepository(db *sqlx.DB) *ContractRepository {
	return &ContractRepository{db: db}
}

// GetAll retrieves all contracts from the database.
func (r *ContractRepository) GetAll() ([]Contract, error) {
	var items []Contract
	err := r.db.Select(&items, "SELECT * FROM contracts ORDER BY id")
	return items, err
}

// GetByID retrieves a contract by its ID.
func (r *ContractRepository) GetByID(id int) (*Contract, error) {
	var item Contract
	err := r.db.Get(&item, "SELECT * FROM contracts WHERE id=$1", id)
	return &item, err
}

// Create inserts a new contract into the database.
func (r *ContractRepository) Create(c *Contract) error {
	return r.db.QueryRow(
		`INSERT INTO contracts (project_id, client_name, start_date, end_date, description)
         VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		c.ProjectID, c.ClientName, c.StartDate, c.EndDate, c.Description,
	).Scan(&c.ID)
}

// Update modifies an existing contract in the database.
func (r *ContractRepository) Update(c *Contract) error {
	_, err := r.db.Exec(
		`UPDATE contracts SET client_name=$1, start_date=$2, end_date=$3, description=$4 WHERE id=$5`,
		c.ClientName, c.StartDate, c.EndDate, c.Description, c.ID,
	)
	return err
}

// Delete removes a contract from the database by its ID.
func (r *ContractRepository) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM contracts WHERE id=$1`, id)
	return err
}
