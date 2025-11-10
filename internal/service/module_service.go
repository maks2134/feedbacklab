package service

import (
	"innotech/internal/models"
	"innotech/internal/repository"
)

type ModuleService struct {
	repo *repository.ModuleRepository
}

func NewModuleService(repo *repository.ModuleRepository) *ModuleService {
	return &ModuleService{repo: repo}
}

func (s *ModuleService) GetAll() ([]models.Module, error) {
	return s.repo.GetAll()
}

func (s *ModuleService) GetByID(id int) (*models.Module, error) {
	return s.repo.GetByID(id)
}

func (s *ModuleService) Create(m *models.Module) error {
	return s.repo.Create(m)
}

func (s *ModuleService) Update(m *models.Module) error {
	return s.repo.Update(m)
}

func (s *ModuleService) Delete(id int) error {
	return s.repo.Delete(id)
}
