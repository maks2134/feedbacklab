package files

import (
	"context"
	"innotech/pkg/minio"
	"log/slog"
)

type Service struct {
	minio  *minio.MinioClient
	logger *slog.Logger
}

func NewService(minioClient *minio.MinioClient, logger *slog.Logger) *Service {
	return &Service{minio: minioClient, logger: logger}
}

func (s *Service) UploadFile(ctx context.Context, filename, path string) (string, error) {
	s.logger.Info("uploading file", "filename", filename)

	err := s.minio.Upload(ctx, filename, path)
	if err != nil {
		s.logger.Error("upload failed", "err", err)
		return "", err
	}

	url, err := s.minio.GetFileURL(filename)
	if err != nil {
		s.logger.Error("cannot get file url", "err", err)
		return "", err
	}

	s.logger.Info("file uploaded", "url", url)
	return url, nil
}
