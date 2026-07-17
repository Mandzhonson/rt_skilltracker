package admin

import (
	"core_service/internal/repository/minio"
	"core_service/internal/repository/postgres"
	"errors"
)

var (
	ErrInvalidPosition = errors.New("invalid position")
	ErrInvalidEmployee = errors.New("user is not employee")
)

type adminService struct {
	userRepo  postgres.UserRepository
	planRepo  postgres.PlanRepository
	skillRepo postgres.SkillRepository
	storage   minio.Storage
}

func NewAdminService(userRepo postgres.UserRepository, planRepo postgres.PlanRepository, skillRepo postgres.SkillRepository, storage minio.Storage) *adminService {
	return &adminService{
		userRepo:  userRepo,
		planRepo:  planRepo,
		skillRepo: skillRepo,
		storage:   storage,
	}
}
