package service

import (
	"innotech/internal/models"
	"innotech/internal/repository"
)

type ContractService struct {
	repo *repository.ContractRepository
}

func NewContractService(repo *repository.ContractRepository) *ContractService {
	return &ContractService{repo: repo}
}

func (s *ContractService) GetAll() ([]models.Contract, error) {
	return s.repo.GetAll()
}

func (s *ContractService) GetByID(id int) (*models.Contract, error) {
	return s.repo.GetByID(id)
}

func (s *ContractService) Create(c *models.Contract) error {
	return s.repo.Create(c)
}

func (s *ContractService) Update(c *models.Contract) error {
	return s.repo.Update(c)
}

func (s *ContractService) Delete(id int) error {
	return s.repo.Delete(id)
}
