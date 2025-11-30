// Package minio - package to create basic operation in MinIO storage
package minio

import (
	"context"
	"innotech/pkg/logger"
	"log"
	"mime/multipart"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioClient provides helper methods for interacting with MinIO storage.
type MinioClient struct {
	Client     *minio.Client
	BucketName string
}

func New(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*MinioClient, error) {
	logger.Info("initializing MinIO client",
		"endpoint", endpoint,
		"bucket", bucket,
		"useSSL", useSSL,
	)

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		logger.Error("failed to create MinIO client",
			"endpoint", endpoint,
			"error", err,
		)
		return nil, err
	}

	ctx := context.Background()
	exists, _ := client.BucketExists(ctx, bucket)
	if !exists {
		logger.Info("bucket does not exist, creating new bucket",
			"bucket", bucket,
		)

		err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			logger.Error("failed to create bucket",
				"bucket", bucket,
				"error", err,
			)
			return nil, err
		}
		log.Println("Bucket created:", bucket)
		logger.Info("bucket created successfully",
			"bucket", bucket,
		)
	} else {
		logger.Debug("bucket already exists",
			"bucket", bucket,
		)
	}

	logger.Info("MinIO client initialized successfully",
		"endpoint", endpoint,
		"bucket", bucket,
	)
	return &MinioClient{
		Client:     client,
		BucketName: bucket,
	}, nil
}

func (m *MinioClient) UploadFile(ctx context.Context, objectName string, file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			return
		}
	}(src)

	logger.Info("uploading file stream to MinIO",
		"bucket", m.BucketName,
		"object_name", objectName,
		"size", file.Size,
	)

	_, err = m.Client.PutObject(ctx, m.BucketName, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		logger.Error("failed to upload file", "error", err)
		return err
	}

	return nil
}

func (m *MinioClient) GetFileURL(objectName string) (string, error) {
	expiry := time.Hour * 24

	pressigndURL, err := m.Client.PresignedGetObject(context.Background(), m.BucketName, objectName, expiry, make(url.Values))
	if err != nil {
		logger.Error("failed to get presigned url", "error", err)
	}

	return pressigndURL.String(), nil
}

func (m *MinioClient) DeleteFile(ctx context.Context, objectName string) error {
	return m.Client.RemoveObject(ctx, m.BucketName, objectName, minio.RemoveObjectOptions{})
}
