// Package minio - package to create basic operation in MinIO storage
package minio

import (
	"context"
	"log"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Client - struct to create basic operation in MinIO storage
type Client struct {
	Client     *minio.Client
	BucketName string
}

// New - functional to create new connection in MinIO
func New(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*Client, error) {
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

	return &Client{
		Client:     client,
		BucketName: bucket,
	}, nil
}

// Upload - functional to upload new file in MinIO
func (m *Client) Upload(ctx context.Context, objectName, filePath string) error {
	_, err := m.Client.FPutObject(ctx, m.BucketName, objectName, filePath, minio.PutObjectOptions{})
	return err
}

// GetFileURL - functional to getting file URL
func (m *Client) GetFileURL(objectName string) (string, error) {
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
