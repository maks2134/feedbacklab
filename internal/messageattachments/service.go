package messageattachments

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"innotech/internal/storage/postgres"
	"innotech/pkg/minio"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, att *postgres.MessageAttachment, file *multipart.FileHeader) error
	GetByID(ctx context.Context, id int) (*postgres.MessageAttachment, error)
	GetByChatID(ctx context.Context, chatID int) ([]postgres.MessageAttachment, error)
	Update(ctx context.Context, att *postgres.MessageAttachment) error
	Delete(ctx context.Context, id int) error
}

type service struct {
	repo        Repository
	minioClient *minio.MinioClient
}

func NewService(repo Repository, minioClient *minio.MinioClient) Service {
	return &service{repo: repo, minioClient: minioClient}
}

func (s *service) Create(ctx context.Context, att *postgres.MessageAttachment, file *multipart.FileHeader) error {
	ext := filepath.Ext(file.Filename)
	objectName := fmt.Sprintf("chat/%d/%s%s", att.ChatID, uuid.New().String(), ext)

	if err := s.minioClient.UploadFile(ctx, objectName, file); err != nil {
		return err
	}

	contentType := file.Header.Get("Content-Type")
	att.FileType = &contentType

	return s.repo.Create(ctx, att)
}

func (s *service) GetByID(ctx context.Context, id int) (*postgres.MessageAttachment, error) {
	att, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	url, err := s.minioClient.GetFileURL(att.FilePath)
	if err == nil {
		att.FilePath = url
	}
	return att, nil
}

func (s *service) GetByChatID(ctx context.Context, chatID int) ([]postgres.MessageAttachment, error) {
	list, err := s.repo.GetByChatID(ctx, chatID)
	if err != nil {
		return nil, err
	}
	for i := range list {
		url, err := s.minioClient.GetFileURL(list[i].FilePath)
		if err == nil {
			list[i].FilePath = url
		}
	}
	return list, nil
}

func (s *service) Delete(ctx context.Context, id int) error {
	att, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	_ = s.minioClient.DeleteFile(ctx, att.FilePath)
	return s.repo.Delete(ctx, id)
}
func (s *service) Update(ctx context.Context, att *postgres.MessageAttachment) error {
	return s.repo.Update(ctx, att)
}
