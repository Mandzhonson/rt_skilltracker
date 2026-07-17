package converter

import (
	"core_service/internal/domain"
	"core_service/internal/repository/model"
)

func ToTaskModel(entity *domain.Task) *model.TaskModel {
	return &model.TaskModel{
		ID:          entity.ID,
		PlanID:      entity.PlanID,
		Title:       entity.Title,
		Description: entity.Description,
		Position:    int16(entity.Position),
		Status:      string(entity.Status),
	}
}

func ToTaskEntity(model *model.TaskModel) *domain.Task {
	return &domain.Task{
		ID:          model.ID,
		PlanID:      model.PlanID,
		Title:       model.Title,
		Description: model.Description,
		Position:    int(model.Position),
		Status:      domain.TaskStatus(model.Status),
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}
