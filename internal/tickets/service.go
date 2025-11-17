package tickets

import (
	"context"
	"innotech/internal/storage/postgres"
)

// Service defines the interface for ticket business logic operations.
type Service interface {
	Create(ctx context.Context, t *postgres.Ticket) error
	GetByID(ctx context.Context, id int) (*postgres.Ticket, error)
	GetAll(ctx context.Context) ([]postgres.Ticket, error)
	Update(ctx context.Context, t *postgres.Ticket) error
	Delete(ctx context.Context, id int) error
}

type ticketService struct {
	repo Repository
}

// NewService creates a new Service instance.
func NewService(repo Repository) Service {
	return &ticketService{repo: repo}
}

func (s *ticketService) Create(ctx context.Context, t *postgres.Ticket) error {
	t.Status = "open"
	return s.repo.Create(ctx, t)
}

func (s *ticketService) GetByID(ctx context.Context, id int) (*postgres.Ticket, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ticketService) GetAll(ctx context.Context) ([]postgres.Ticket, error) {
	return s.repo.GetAll(ctx)
}

func (s *ticketService) Update(ctx context.Context, t *postgres.Ticket) error {
	return s.repo.Update(ctx, t)
}

func (s *ticketService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
