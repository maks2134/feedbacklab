package messageattachments

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"innotech/internal/storage/postgres"
	"innotech/pkg/errors"
	"innotech/pkg/logger"
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
		return errors.NewInternalError("service", "upload file to minio", err)
	}

	contentType := file.Header.Get("Content-Type")
	att.FileType = &contentType
	att.FilePath = objectName

	if err := s.repo.Create(ctx, att); err != nil {
		return errors.NewInternalError("service", "create message attachment", err)
	}

	return nil
}

func (s *service) GetByID(ctx context.Context, id int) (*postgres.MessageAttachment, error) {
	att, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.NewNotFoundError("message_attachment.error.not_found", "message attachment", err)
	}

	url, err := s.minioClient.GetFileURL(att.FilePath)
	if err != nil {
		logger.Warn("failed to get presigned URL for message attachment",
			"id", id,
			"file_path", att.FilePath,
			"error", err,
		)
		// Continue with original file path if URL generation fails
	} else {
		att.FilePath = url
	}
	return att, nil
}

func (s *service) GetByChatID(ctx context.Context, chatID int) ([]postgres.MessageAttachment, error) {
	list, err := s.repo.GetByChatID(ctx, chatID)
	if err != nil {
		return nil, errors.NewInternalError("service", "get message attachments by chat id", err)
	}

	for i := range list {
		url, err := s.minioClient.GetFileURL(list[i].FilePath)
		if err != nil {
			logger.Warn("failed to get presigned URL for message attachment",
				"chat_id", chatID,
				"file_path", list[i].FilePath,
				"error", err,
			)
			// Continue with original file path if URL generation fails
		} else {
			list[i].FilePath = url
		}
	}
	return list, nil
}

func (s *service) Delete(ctx context.Context, id int) error {
	att, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.NewNotFoundError("message_attachment.error.not_found", "message attachment", err)
	}

	if err := s.minioClient.DeleteFile(ctx, att.FilePath); err != nil {
		logger.Warn("failed to delete file from minio",
			"id", id,
			"file_path", att.FilePath,
			"error", err,
		)
		// Continue with database deletion even if minio deletion fails
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return errors.NewInternalError("service", "delete message attachment", err)
	}

	return nil
}

func (s *service) Update(ctx context.Context, att *postgres.MessageAttachment) error {
	if err := s.repo.Update(ctx, att); err != nil {
		return errors.NewInternalError("service", "update message attachment", err)
	}
	return nil
}
