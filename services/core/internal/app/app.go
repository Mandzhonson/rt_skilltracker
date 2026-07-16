package app

import (
	"context"
	"core_service/internal/clients/minio"
	"core_service/internal/clients/ollama"
	"core_service/internal/clients/postgres"
	"core_service/internal/clients/redis"
	"core_service/internal/config"
	"core_service/internal/domain"
	"core_service/internal/pkg/jwt"
	minioStorage "core_service/internal/repository/minio"
	postgresRepository "core_service/internal/repository/postgres"

	redisRepository "core_service/internal/repository/redis"
	router "core_service/internal/transport/http"
	"core_service/internal/transport/http/handler"
	"core_service/internal/transport/http/middleware"
	"core_service/internal/usecase/admin"
	"core_service/internal/usecase/ai"
	"core_service/internal/usecase/auth"
	"core_service/internal/usecase/plan"
	"core_service/internal/usecase/skill"
	"core_service/internal/usecase/task"
	"core_service/internal/usecase/test"
	"core_service/internal/usecase/user"
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
		syscall.SIGKILL,
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

	minio, err := minio.InitMinio(ctx, cfg.Minio)
	if err != nil {
		return err
	}

	ollama := ollama.InitOllama(cfg.Ollama)

	minioStorage := minioStorage.NewMinioStorage(minio, cfg.Minio.Bucket)
	sessionRepository := redisRepository.NewRedisSessionRepository(rdb)
	userRepository := postgresRepository.NewUserRepository(pool)
	authRepository := postgresRepository.NewAuthRepository(pool)
	planRepository := postgresRepository.NewPlanRepository(pool)
	taskRepository := postgresRepository.NewTaskRepository(pool)
	skillRepository := postgresRepository.NewSkillRepository(pool)
	testRepository := postgresRepository.NewTestRepository(pool)

	jwtService := jwt.NewJWTService(cfg.JWT)
	aiService := ai.NewAIService(ollama)
	authService := auth.NewAuthService(authRepository, userRepository, jwtService, sessionRepository)
	userService := user.NewUserService(userRepository, minioStorage, skillRepository, planRepository)
	adminService := admin.NewAdminService(userRepository, planRepository, skillRepository, minioStorage)
	planService := plan.NewPlanService(planRepository, userRepository, taskRepository, skillRepository, testRepository, *aiService)
	taskService := task.NewTaskService(taskRepository, planRepository, planService)
	testService := test.NewTestService(testRepository, *taskService, planRepository, planService)
	skillService := skill.NewSkillService(skillRepository, userRepository)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	adminHandler := handler.NewAdminHandler(adminService)
	planHandler := handler.NewPlanHandler(planService)
	taskHandler := handler.NewTaskHandler(taskService)
	testHandler := handler.NewTestHandler(testService)
	skillHandler := handler.NewSkillHandler(skillService)

	authMiddleware := middleware.AuthMiddleware(jwtService, sessionRepository)
	adminMiddleware := middleware.AdminMiddleware()
	managerMiddleware := middleware.ManagerMiddleware()
	employeeMiddleware := middleware.EmployeeMiddleware()

	router := router.NewRouter(authHandler, userHandler, adminHandler, planHandler, taskHandler, testHandler, skillHandler, authMiddleware, adminMiddleware, managerMiddleware, employeeMiddleware)

	err = userService.CreateAdmin(
		ctx,
		user.CreateUserInput{
			Email:     cfg.Admin.Email,
			Password:  cfg.Admin.Password,
			FirstName: cfg.Admin.FirstName,
			LastName:  cfg.Admin.LastName,
			Role:      domain.RoleAdmin,
		},
	)

	if err != nil {
		return err
	}
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
