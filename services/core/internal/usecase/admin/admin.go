package admin

import (
	"core_service/internal/repository/minio"
	"core_service/internal/repository/postgres"
)

type adminService struct {
	userRepo postgres.UserRepository
	storage  minio.Storage
}

func NewAdminService(userRepo postgres.UserRepository, storage minio.Storage) *adminService {
	return &adminService{
		userRepo: userRepo,
		storage:  storage,
	}
}
