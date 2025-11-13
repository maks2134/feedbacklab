package contract

type ContractService struct {
	repo *ContractRepository
}

func NewContractService(repo *ContractRepository) *ContractService {
	return &ContractService{repo: repo}
}

func (s *ContractService) GetAll() ([]Contract, error) {
	return s.repo.GetAll()
}

func (s *ContractService) GetByID(id int) (*Contract, error) {
	return s.repo.GetByID(id)
}

func (s *ContractService) Create(c *Contract) error {
	return s.repo.Create(c)
}

func (s *ContractService) Update(c *Contract) error {
	return s.repo.Update(c)
}

func (s *ContractService) Delete(id int) error {
	return s.repo.Delete(id)
}
