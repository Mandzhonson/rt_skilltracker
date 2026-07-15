package plan

import (
	"context"
	"core_service/internal/domain"
)

func (s *planService) generateManualPlan(ctx context.Context, plan *domain.Plan) {
	ctx = context.Background()
	err := s.generateTestForPlan(ctx, plan, nil)
	if err != nil {
		_ = s.planRepo.UpdateGenerationStatus(ctx, plan.ID, domain.GenerationFailed)
		return
	}

	position, err := s.taskRepo.GetNextPosition(ctx, plan.ID)
	if err == nil {

		testingTask := domain.NewTask(plan.ID, "Пройти тестирование", nil, position)
		_, _ = s.taskRepo.Create(ctx, testingTask)
	}

	_ = s.planRepo.UpdateGenerationStatus(ctx, plan.ID, domain.GenerationReady)
}
