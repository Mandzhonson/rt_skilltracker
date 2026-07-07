package admin

import "core_service/internal/repository/postgres"

type adminService struct {
	userRepo postgres.UserRepository
}

func NewAdminService(userRepo postgres.UserRepository) *adminService {
	return &adminService{
		userRepo: userRepo,
	}
}
