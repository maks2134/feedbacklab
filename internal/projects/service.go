package projects

import (
	"context"
	"innotech/internal/storage/postgres"
)

// Service defines the interface for project business logic operations.
type Service interface {
	Create(ctx context.Context, p *postgres.Project) error
	GetByID(ctx context.Context, id int) (*postgres.Project, error)
	GetAll(ctx context.Context) ([]postgres.Project, error)
	Update(ctx context.Context, p *postgres.Project) error
	Delete(ctx context.Context, id int) error
}

type projectService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &projectService{repo: repo}
}

func (s *projectService) Create(ctx context.Context, p *postgres.Project) error {
	return s.repo.Create(ctx, p)
}

func (s *projectService) GetByID(ctx context.Context, id int) (*postgres.Project, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *projectService) GetAll(ctx context.Context) ([]postgres.Project, error) {
	return s.repo.GetAll(ctx)
}

func (s *projectService) Update(ctx context.Context, p *postgres.Project) error {
	return s.repo.Update(ctx, p)
}

func (s *projectService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
