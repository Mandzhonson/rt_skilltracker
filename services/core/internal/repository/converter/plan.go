package converter

import (
	"core_service/internal/domain"
	"core_service/internal/repository/model"
)

func ToPlanModel(plan *domain.Plan) *model.PlanModel {
	return &model.PlanModel{
		ID:           plan.ID,
		EmployeeID:   plan.EmployeeID,
		CreatedBy:    plan.CreatedBy,
		Title:        plan.Title,
		Description:  plan.Description,
		CreationType: string(plan.CreationType),
		Progress:     int16(plan.Progress),
		Status:       string(plan.Status),
		CreatedAt:    plan.CreatedAt,
		UpdatedAt:    plan.UpdatedAt,
	}
}

func ToPlanEntity(model *model.PlanModel) *domain.Plan {
	return &domain.Plan{
		ID:           model.ID,
		EmployeeID:   model.EmployeeID,
		CreatedBy:    model.CreatedBy,
		Title:        model.Title,
		Description:  model.Description,
		CreationType: domain.CreationType(model.CreationType),
		Progress:     int(model.Progress),
		Status:       domain.PlanStatus(model.Status),
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}
}
