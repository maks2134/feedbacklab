package ticketchats

import (
	"errors"
	"innotech/internal/storage/postgres"
)

// Service defines the interface for ticket chat business logic operations.
type Service interface {
	Create(chat *postgres.TicketChat) error
	GetByID(id int) (*postgres.TicketChat, error)
	GetByTicketID(ticketID int) ([]postgres.TicketChat, error)
	Update(chat *postgres.TicketChat) error
	Delete(id int) error
}

type service struct {
	repo Repository
}

// NewService creates a new Service instance.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(chat *postgres.TicketChat) error {
	if chat.Message == "" {
		return errors.New("message cannot be empty")
	}
	return s.repo.Create(chat)
}

func (s *service) GetByID(id int) (*postgres.TicketChat, error) {
	return s.repo.GetByID(id)
}

func (s *service) GetByTicketID(ticketID int) ([]postgres.TicketChat, error) {
	return s.repo.GetByTicketID(ticketID)
}

func (s *service) Update(chat *postgres.TicketChat) error {
	return s.repo.Update(chat)
}

func (s *service) Delete(id int) error {
	return s.repo.Delete(id)
}
