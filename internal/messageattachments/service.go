package messageattachments

import (
	"errors"
	"innotech/internal/storage/postgres"
)

// Service defines the interface for message attachment business logic operations.
type Service interface {
	Create(att *postgres.MessageAttachment) error
	GetByID(id int) (*postgres.MessageAttachment, error)
	GetByChatID(chatID int) ([]postgres.MessageAttachment, error)
	Update(att *postgres.MessageAttachment) error
	Delete(id int) error
}

type service struct {
	repo Repository
}

// NewService creates a new Service instance.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(att *postgres.MessageAttachment) error {
	if att.FilePath == "" {
		return errors.New("file_path cannot be empty")
	}
	if att.ChatID == 0 {
		return errors.New("chat_id is required")
	}
	return s.repo.Create(att)
}

func (s *service) GetByID(id int) (*postgres.MessageAttachment, error) {
	return s.repo.GetByID(id)
}

func (s *service) GetByChatID(chatID int) ([]postgres.MessageAttachment, error) {
	return s.repo.GetByChatID(chatID)
}

func (s *service) Update(att *postgres.MessageAttachment) error {
	return s.repo.Update(att)
}

func (s *service) Delete(id int) error {
	return s.repo.Delete(id)
}
