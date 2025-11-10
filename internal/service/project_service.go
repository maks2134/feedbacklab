package service

import (
	"innotech/internal/models"
	"innotech/internal/repository"
)

type ProjectService struct {
	repo *repository.ProjectRepository
}

func NewProjectService(repo *repository.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) GetAll() ([]models.Project, error) {
	return s.repo.GetAll()
}

func (s *ProjectService) GetByID(id int) (*models.Project, error) {
	return s.repo.GetByID(id)
}

func (s *ProjectService) Create(p *models.Project) error {
	return s.repo.Create(p)
}

func (s *ProjectService) Update(p *models.Project) error {
	return s.repo.Update(p)
}

func (s *ProjectService) Delete(id int) error {
	return s.repo.Delete(id)
}
