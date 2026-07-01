package app

import (
	"context"
	"core_service/internal/clients/postgres"
	"core_service/internal/clients/redis"
	"core_service/internal/config"
	postgresRepository "core_service/internal/repository/postgres"
	router "core_service/internal/transport/http"
	"core_service/internal/transport/http/handler"
	"core_service/internal/usecase/auth"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
)

func Run() error {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
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

	_ = rdb // пока не используется

	authRepository := postgresRepository.NewAuthRepository(pool)

	authService := auth.NewAuthService(authRepository)

	authHandler := handler.NewAuthHandler(authService)

	router := router.NewRouter(authHandler)

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port),
		Handler:      router,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// сюда позже добавишь slog.Error(...)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		cfg.HTTP.ShutdownTimeout,
	)
	defer cancel()

	return srv.Shutdown(shutdownCtx)
}
