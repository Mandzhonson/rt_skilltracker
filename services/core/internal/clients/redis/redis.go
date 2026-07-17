package redis

import (
	"context"
	"core_service/internal/config"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func Init(ctx context.Context, cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           0,
		MaxIdleConns: 25,
		MinIdleConns: 5,
		MaxRetries:   5,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,
	})
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return client, nil
}
