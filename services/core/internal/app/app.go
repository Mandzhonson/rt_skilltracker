package app

import (
	"context"
	"core_service/internal/clients/postgres"
	"core_service/internal/clients/redis"
	"core_service/internal/config"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
)

func Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	pool, err := postgres.InitPool(ctx, cfg.PostgresDSN())
	if err != nil {
		return err
	}
	defer pool.Close()
	rdb, err := redis.Init(ctx, cfg.Redis)
	if err != nil {
		return err
	}
	defer rdb.Close()
	// minio, err := minio.InitMinio(ctx, cfg.Minio)
	// if err != nil {
	// 	return err
	// }

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port),
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		}
	}()

	<-ctx.Done()
	cancel()
	shtCtx, cancel := context.WithTimeout(context.Background(), cfg.HTTP.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(shtCtx); err != nil {
		return err
	}
	return nil
}
