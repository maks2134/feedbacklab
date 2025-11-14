package ticket_attachments

import (
	"errors"
	"innotech/internal/storage/postgres"
)

type Service interface {
	Create(att *postgres.TicketAttachment) error
	GetByID(id int) (*postgres.TicketAttachment, error)
	GetByTicketID(ticketID int) ([]postgres.TicketAttachment, error)
	Update(att *postgres.TicketAttachment) error
	Delete(id int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(att *postgres.TicketAttachment) error {
	if att.FilePath == "" {
		return errors.New("file_path cannot be empty")
	}
	return s.repo.Create(att)
}

func (s *service) GetByID(id int) (*postgres.TicketAttachment, error) {
	return s.repo.GetByID(id)
}

func (s *service) GetByTicketID(ticketID int) ([]postgres.TicketAttachment, error) {
	return s.repo.GetByTicketID(ticketID)
}

func (s *service) Update(att *postgres.TicketAttachment) error {
	return s.repo.Update(att)
}

func (s *service) Delete(id int) error {
	return s.repo.Delete(id)
}
