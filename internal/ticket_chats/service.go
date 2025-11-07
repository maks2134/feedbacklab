package ticket_chats

import "errors"

type Service interface {
	Create(chat *TicketChat) error
	GetByID(id int) (*TicketChat, error)
	GetByTicketID(ticketID int) ([]TicketChat, error)
	Update(chat *TicketChat) error
	Delete(id int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(chat *TicketChat) error {
	if chat.Message == "" {
		return errors.New("message cannot be empty")
	}
	return s.repo.Create(chat)
}

func (s *service) GetByID(id int) (*TicketChat, error) {
	return s.repo.GetByID(id)
}

func (s *service) GetByTicketID(ticketID int) ([]TicketChat, error) {
	return s.repo.GetByTicketID(ticketID)
}

func (s *service) Update(chat *TicketChat) error {
	return s.repo.Update(chat)
}

func (s *service) Delete(id int) error {
	return s.repo.Delete(id)
}
