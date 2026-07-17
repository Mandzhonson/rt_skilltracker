package converter

import (
	"core_service/internal/domain"
	"core_service/internal/repository/model"
)

func ToEmployeeProfileEntity(
	user *model.EmployeeProfileModel,
	skills []*domain.Skill,
	plans []*domain.Plan,
) *domain.EmployeeProfile {

	return &domain.EmployeeProfile{
		User: &domain.User{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Position:  user.Position,
			CreatedAt: user.CreatedAt,
		},
		Skills: skills,
		Plans:  plans,
	}
}
