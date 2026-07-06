package minio

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
)

type Storage interface {
	UploadAvatar(ctx context.Context, objectName string, file io.Reader, size int64, contentType string) (string, error)
	DeleteAvatar(ctx context.Context, objectName string) error
	GetAvatar(ctx context.Context, objectName string) (io.ReadCloser, string, error)
}

type minioStorage struct {
	client     *minio.Client
	bucketName string
}

func NewMinioStorage(client *minio.Client, bucketName string) *minioStorage {
	return &minioStorage{
		client:     client,
		bucketName: bucketName,
	}
}

func (s *minioStorage) UploadAvatar(ctx context.Context, objectName string, file io.Reader, size int64, contentType string) (string, error) {
	_, err := s.client.PutObject(ctx, s.bucketName, objectName, file, size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", fmt.Errorf("minio.UploadAvatar: %w", err)
	}
	return objectName, nil
}

func (s *minioStorage) DeleteAvatar(ctx context.Context, objectName string) error {
	err := s.client.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("minio.DeleteAvatar: %w", err)
	}
	return nil
}

func (s *minioStorage) GetAvatar(ctx context.Context, objectName string) (io.ReadCloser, string, error) {
	object, err := s.client.GetObject(ctx, s.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, "", fmt.Errorf("minio.GetAvatar: %w", err)
	}

	info, err := object.Stat()
	if err != nil {
		return nil, "", fmt.Errorf("minio.GetAvatar: %w", err)
	}

	return object, info.ContentType, nil
}
