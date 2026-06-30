package minio

import (
	"context"
	"core_service/internal/config"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func InitMinio(ctx context.Context, cfg config.MinioConfig) (*minio.Client, error) {
	client, err := minio.New(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port), &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.User, cfg.Password, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}
	// TODO: дописать в .env имя бакета и создать однажды его через sync.Once.Do
	return client, nil
}
