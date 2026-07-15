package converter

import (
	"core_service/internal/domain"
	"core_service/internal/repository/model"
)

func ToSkillEntity(skill *model.SkillModel) *domain.Skill {
	return &domain.Skill{
		ID:          skill.ID,
		Name:        skill.Name,
		Category:    skill.Category,
		Description: skill.Description,
		CreatedAt:   skill.CreatedAt,
	}
}

func ToSkillModel(skill *domain.Skill) *model.SkillModel {
	return &model.SkillModel{
		ID:          skill.ID,
		Name:        skill.Name,
		Category:    skill.Category,
		Description: skill.Description,
		CreatedAt:   skill.CreatedAt,
	}
}
