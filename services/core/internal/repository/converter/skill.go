package converter

import (
	"core_service/internal/domain"
	"core_service/internal/repository/model"
)

func ToUserSkillEntity(skill *model.UserSkillModel) *domain.UserSkill {
	return &domain.UserSkill{
		ID:          skill.ID,
		UserID:      skill.UserID,
		PlanID:      skill.PlanID,
		Name:        skill.Name,
		ConfirmedAt: skill.ConfirmedAt,
	}
}

func ToUserSkillModel(skill *domain.UserSkill) *model.UserSkillModel {
	return &model.UserSkillModel{
		ID:          skill.ID,
		UserID:      skill.UserID,
		PlanID:      skill.PlanID,
		Name:        skill.Name,
		ConfirmedAt: skill.ConfirmedAt,
	}
}
