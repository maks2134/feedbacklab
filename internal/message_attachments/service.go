package message_attachments

import "errors"

type Service interface {
	Create(att *MessageAttachment) error
	GetByID(id int) (*MessageAttachment, error)
	GetByChatID(chatID int) ([]MessageAttachment, error)
	Update(att *MessageAttachment) error
	Delete(id int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(att *MessageAttachment) error {
	if att.FilePath == "" {
		return errors.New("file_path cannot be empty")
	}
	if att.ChatID == 0 {
		return errors.New("chat_id is required")
	}
	return s.repo.Create(att)
}

func (s *service) GetByID(id int) (*MessageAttachment, error) {
	return s.repo.GetByID(id)
}

func (s *service) GetByChatID(chatID int) ([]MessageAttachment, error) {
	return s.repo.GetByChatID(chatID)
}

func (s *service) Update(att *MessageAttachment) error {
	return s.repo.Update(att)
}

func (s *service) Delete(id int) error {
	return s.repo.Delete(id)
}
