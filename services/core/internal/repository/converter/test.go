package converter

import (
	"core_service/internal/domain"
	"core_service/internal/repository/model"
)

func ToTestEntity(test *model.TestModel) *domain.Test {
	return &domain.Test{
		ID:        test.ID,
		PlanID:    test.PlanID,
		CreatedAt: test.CreatedAt,
	}
}

func ToTestModel(test *domain.Test) *model.TestModel {
	return &model.TestModel{
		ID:        test.ID,
		PlanID:    test.PlanID,
		CreatedAt: test.CreatedAt,
	}
}
