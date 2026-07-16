package skill

import (
	"core_service/internal/repository/postgres"
	"errors"
)

var ErrForbidden = errors.New("forbidden")

type skillService struct {
	skillRepo postgres.SkillRepository
	userRepo  postgres.UserRepository
}

func NewSkillService(skillRepo postgres.SkillRepository, userRepo postgres.UserRepository) *skillService {
	return &skillService{
		skillRepo: skillRepo,
		userRepo:  userRepo,
	}
}
