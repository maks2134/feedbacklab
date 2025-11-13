package modules

import "context"

type Service interface {
	Create(ctx context.Context, m *Module) error
	GetByID(ctx context.Context, id int) (*Module, error)
	GetAll(ctx context.Context) ([]Module, error)
	Update(ctx context.Context, m *Module) error
	Delete(ctx context.Context, id int) error
}

type moduleService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &moduleService{repo: repo}
}

func (s *moduleService) Create(ctx context.Context, m *Module) error {
	return s.repo.Create(ctx, m)
}

func (s *moduleService) GetByID(ctx context.Context, id int) (*Module, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *moduleService) GetAll(ctx context.Context) ([]Module, error) {
	return s.repo.GetAll(ctx)
}

func (s *moduleService) Update(ctx context.Context, m *Module) error {
	return s.repo.Update(ctx, m)
}

func (s *moduleService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
