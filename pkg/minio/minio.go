package minio

import (
	"context"
	"log"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	Client     *minio.Client
	BucketName string
}

func New(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*MinioClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	exists, _ := client.BucketExists(ctx, bucket)
	if !exists {
		err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
		log.Println("Bucket created:", bucket)
	}

	return &MinioClient{
		Client:     client,
		BucketName: bucket,
	}, nil
}

func (m *MinioClient) Upload(ctx context.Context, objectName, filePath string) error {
	_, err := m.Client.FPutObject(ctx, m.BucketName, objectName, filePath, minio.PutObjectOptions{})
	return err
}

func (m *MinioClient) GetFileURL(objectName string) (string, error) {
	reqParams := make(url.Values)

	presignedURL, err := m.Client.PresignedGetObject(
		context.Background(),
		m.BucketName,
		objectName,
		time.Hour*24,
		reqParams,
	)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}
