package service

import (
	"innotech/internal/models"
	"innotech/internal/repository"
)

type DocumentationService struct {
	repo *repository.DocumentationRepository
}

func NewDocumentationService(repo *repository.DocumentationRepository) *DocumentationService {
	return &DocumentationService{repo: repo}
}

func (s *DocumentationService) GetAll() ([]models.Documentation, error) {
	return s.repo.GetAll()
}

func (s *DocumentationService) GetByID(id int) (*models.Documentation, error) {
	return s.repo.GetByID(id)
}

func (s *DocumentationService) Create(d *models.Documentation) error {
	return s.repo.Create(d)
}

func (s *DocumentationService) Update(d *models.Documentation) error {
	return s.repo.Update(d)
}

func (s *DocumentationService) Delete(id int) error {
	return s.repo.Delete(id)
}
