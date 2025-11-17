package contract

// ContractService handles business logic for contract operations.
type ContractService struct {
	repo *ContractRepository
}

// NewContractService creates a new ContractService instance.
func NewContractService(repo *ContractRepository) *ContractService {
	return &ContractService{repo: repo}
}

// GetAll retrieves all contracts.
func (s *ContractService) GetAll() ([]Contract, error) {
	return s.repo.GetAll()
}

// GetByID retrieves a contract by its ID.
func (s *ContractService) GetByID(id int) (*Contract, error) {
	return s.repo.GetByID(id)
}

// Create creates a new contract.
func (s *ContractService) Create(c *Contract) error {
	return s.repo.Create(c)
}

// Update updates an existing contract.
func (s *ContractService) Update(c *Contract) error {
	return s.repo.Update(c)
}

// Delete deletes a contract by its ID.
func (s *ContractService) Delete(id int) error {
	return s.repo.Delete(id)
}
