package service

import (
	"innotech/internal/models"
	"innotech/internal/repository"
)

type UserProjectService struct {
	repo *repository.UserProjectRepository
}

func NewUserProjectService(repo *repository.UserProjectRepository) *UserProjectService {
	return &UserProjectService{repo: repo}
}

func (s *UserProjectService) GetAll() ([]models.UserProject, error) {
	return s.repo.GetAll()
}

func (s *UserProjectService) GetByID(id int) (*models.UserProject, error) {
	return s.repo.GetByID(id)
}

func (s *UserProjectService) Create(up *models.UserProject) error {
	return s.repo.Create(up)
}

func (s *UserProjectService) Update(up *models.UserProject) error {
	return s.repo.Update(up)
}

func (s *UserProjectService) Delete(id int) error {
	return s.repo.Delete(id)
}
