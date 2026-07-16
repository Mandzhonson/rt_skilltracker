package minio

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/minio/minio-go/v7"
)
//go:generate mockgen -source=storage.go -destination=mocks/mock_storage.go -package=mocks
type Storage interface {
	UploadAvatar(ctx context.Context, objectName string, file io.Reader, size int64, contentType string) (string, error)
	DeleteAvatar(ctx context.Context, objectName string) error
	GetAvatar(ctx context.Context, objectName string) (io.ReadCloser, string, error)
}

type minioStorage struct {
	client     *minio.Client
	bucketName string
	log        *slog.Logger
}

func NewMinioStorage(client *minio.Client, bucketName string, log *slog.Logger) *minioStorage {
	return &minioStorage{
		client:     client,
		bucketName: bucketName,
		log:        log,
	}
}

func (s *minioStorage) UploadAvatar(ctx context.Context, objectName string, file io.Reader, size int64, contentType string) (string, error) {
	_, err := s.client.PutObject(ctx, s.bucketName, objectName, file, size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		s.log.Error("Failed to upload avatar to MinIO",
			slog.String("error", err.Error()),
			slog.String("object_name", objectName),
			slog.String("bucket", s.bucketName),
			slog.Int64("size", size),
			slog.String("content_type", contentType),
		)
		return "", fmt.Errorf("minio.UploadAvatar: %w", err)
	}

	s.log.Info("Avatar uploaded successfully to MinIO",
		slog.String("object_name", objectName),
		slog.String("bucket", s.bucketName),
		slog.Int64("size", size),
	)
	return objectName, nil
}

func (s *minioStorage) DeleteAvatar(ctx context.Context, objectName string) error {
	err := s.client.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		s.log.Error("Failed to delete avatar from MinIO",
			slog.String("error", err.Error()),
			slog.String("object_name", objectName),
			slog.String("bucket", s.bucketName),
		)
		return fmt.Errorf("minio.DeleteAvatar: %w", err)
	}

	s.log.Info("Avatar deleted successfully from MinIO",
		slog.String("object_name", objectName),
		slog.String("bucket", s.bucketName),
	)
	return nil
}

func (s *minioStorage) GetAvatar(ctx context.Context, objectName string) (io.ReadCloser, string, error) {
	object, err := s.client.GetObject(ctx, s.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		s.log.Error("Failed to get avatar from MinIO",
			slog.String("error", err.Error()),
			slog.String("object_name", objectName),
			slog.String("bucket", s.bucketName),
		)
		return nil, "", fmt.Errorf("minio.GetAvatar: %w", err)
	}

	info, err := object.Stat()
	if err != nil {
		s.log.Error("Failed to stat avatar object from MinIO",
			slog.String("error", err.Error()),
			slog.String("object_name", objectName),
			slog.String("bucket", s.bucketName),
		)
		object.Close()
		return nil, "", fmt.Errorf("minio.GetAvatar: %w", err)
	}

	s.log.Info("Avatar retrieved successfully from MinIO",
		slog.String("object_name", objectName),
		slog.String("bucket", s.bucketName),
		slog.String("content_type", info.ContentType),
	)
	return object, info.ContentType, nil
}
