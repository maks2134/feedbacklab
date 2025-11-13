package documentations

import "context"

type Service interface {
	Create(ctx context.Context, d *Documentation) error
	GetByID(ctx context.Context, id int) (*Documentation, error)
	GetAll(ctx context.Context) ([]Documentation, error)
	Update(ctx context.Context, d *Documentation) error
	Delete(ctx context.Context, id int) error
}

type documentationService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &documentationService{repo: repo}
}

func (s *documentationService) Create(ctx context.Context, d *Documentation) error {
	return s.repo.Create(ctx, d)
}

func (s *documentationService) GetByID(ctx context.Context, id int) (*Documentation, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *documentationService) GetAll(ctx context.Context) ([]Documentation, error) {
	return s.repo.GetAll(ctx)
}

func (s *documentationService) Update(ctx context.Context, d *Documentation) error {
	return s.repo.Update(ctx, d)
}

func (s *documentationService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
