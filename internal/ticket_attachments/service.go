package ticket_attachments

import "errors"

type Service interface {
	Create(att *TicketAttachment) error
	GetByID(id int) (*TicketAttachment, error)
	GetByTicketID(ticketID int) ([]TicketAttachment, error)
	Update(att *TicketAttachment) error
	Delete(id int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(att *TicketAttachment) error {
	if att.FilePath == "" {
		return errors.New("file_path cannot be empty")
	}
	return s.repo.Create(att)
}

func (s *service) GetByID(id int) (*TicketAttachment, error) {
	return s.repo.GetByID(id)
}

func (s *service) GetByTicketID(ticketID int) ([]TicketAttachment, error) {
	return s.repo.GetByTicketID(ticketID)
}

func (s *service) Update(att *TicketAttachment) error {
	return s.repo.Update(att)
}

func (s *service) Delete(id int) error {
	return s.repo.Delete(id)
}
