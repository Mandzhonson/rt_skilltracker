package plan

import (
	"context"

	"core_service/internal/domain"
)

func (s *planService) attachTesting(ctx context.Context, plan *domain.Plan) error {

	tasks, err := s.taskRepo.ListByPlanID(ctx, plan.ID)
	if err != nil {
		return err
	}

	testingTask := domain.NewTask(plan.ID, "Пройти тестирование", new("Необходимо успешно пройти итоговый тест (минимум 70%)."), len(tasks)+1)
	_, err = s.taskRepo.Create(ctx, testingTask)
	if err != nil {
		return err
	}
	tasks = append(tasks, testingTask)
	return s.generateTestForPlan(ctx, plan, tasks)
}
