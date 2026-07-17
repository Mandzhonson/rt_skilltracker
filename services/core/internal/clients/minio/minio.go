package minio

import (
	"context"
	"core_service/internal/config"
	"fmt"
	"sync"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	bucketInit sync.Once
	bucketErr  error
)

func InitMinio(ctx context.Context, cfg config.MinioConfig) (*minio.Client, error) {
	client, err := minio.New(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port), &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.User, cfg.Password, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}
	bucketInit.Do(func() {
		var exists bool

		exists, bucketErr = client.BucketExists(ctx, cfg.Bucket)
		if bucketErr != nil {
			return
		}

		if !exists {
			bucketErr = client.MakeBucket(
				ctx,
				cfg.Bucket,
				minio.MakeBucketOptions{},
			)
		}
	})

	if bucketErr != nil {
		return nil, bucketErr
	}

	return client, nil
}
