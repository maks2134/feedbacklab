package ticketattachments

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
	Create(ctx context.Context, att *postgres.TicketAttachment, file *multipart.FileHeader) error
	GetByID(ctx context.Context, id int) (*postgres.TicketAttachment, error)
	GetByTicketID(ctx context.Context, ticketID int) ([]postgres.TicketAttachment, error)
	Update(ctx context.Context, att *postgres.TicketAttachment) error
	Delete(ctx context.Context, id int) error
}

type service struct {
	repo        Repository
	minioClient *minio.MinioClient
}

func NewService(repo Repository, minioClient *minio.MinioClient) Service {
	return &service{
		repo:        repo,
		minioClient: minioClient,
	}
}

func (s *service) Create(ctx context.Context, att *postgres.TicketAttachment, file *multipart.FileHeader) error {
	ext := filepath.Ext(file.Filename)
	objectName := fmt.Sprintf("tickets/%d/%s%s", att.TicketID, uuid.New().String(), ext)

	if err := s.minioClient.UploadFile(ctx, objectName, file); err != nil {
		return err
	}

	att.FilePath = objectName

	contentType := file.Header.Get("Content-Type")
	att.FileType = &contentType

	return s.repo.Create(ctx, att)
}

func (s *service) GetByID(ctx context.Context, id int) (*postgres.TicketAttachment, error) {
	att, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	signedURL, err := s.minioClient.GetFileURL(att.FilePath)
	if err == nil {
		att.FilePath = signedURL
	}

	return att, nil
}

func (s *service) GetByTicketID(ctx context.Context, ticketID int) ([]postgres.TicketAttachment, error) {
	list, err := s.repo.GetByTicketID(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	for i := range list {
		signedURL, err := s.minioClient.GetFileURL(list[i].FilePath)
		if err == nil {
			list[i].FilePath = signedURL
		}
	}

	return list, nil
}

func (s *service) Update(ctx context.Context, att *postgres.TicketAttachment) error {
	return s.repo.Update(ctx, att)
}

func (s *service) Delete(ctx context.Context, id int) error {
	att, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	_ = s.minioClient.DeleteFile(ctx, att.FilePath)

	return s.repo.Delete(ctx, id)
}
