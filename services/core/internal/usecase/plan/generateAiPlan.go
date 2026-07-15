package plan

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/usecase/ai"

	"github.com/google/uuid"
)

func (s *planService) generateAIPlan(ctx context.Context, planID uuid.UUID) {

	plan, err := s.planRepo.GetByID(ctx, planID)
	if err != nil {
		return
	}

	employee, err := s.userRepo.GetById(ctx, plan.EmployeeID)
	if err != nil {
		_ = s.planRepo.UpdateGenerationStatus(ctx, planID, domain.GenerationFailed)
		return
	}

	skills, err := s.skillRepo.ListByUserID(ctx, employee.ID)
	if err != nil {
		_ = s.planRepo.UpdateGenerationStatus(ctx, planID, domain.GenerationFailed)
		return
	}

	existingSkills := make([]string, 0, len(skills))
	for _, s := range skills {
		existingSkills = append(existingSkills, s.Name)
	}

	generated, err := s.aiService.GeneratePlan(
		ctx,
		ai.GeneratePlanInput{
			Topic:          plan.Title,
			Description:    deref(plan.Description),
			ExistingSkills: existingSkills,
			Position:       employee.Position,
		},
	)
	if err != nil {
		_ = s.planRepo.UpdateGenerationStatus(ctx, planID, domain.GenerationFailed)
		return
	}

	var description *string
	if generated.Description != "" {
		description = &generated.Description
	}

	err = s.planRepo.UpdateAIContent(
		ctx,
		planID,
		generated.Title,
		description,
	)
	if err != nil {
		_ = s.planRepo.UpdateGenerationStatus(ctx, planID, domain.GenerationFailed)
		return
	}

	var tasks []*domain.Task

	for i, t := range generated.Tasks {

		var desc *string
		if t.Description != "" {
			desc = &t.Description
		}

		task := domain.NewTask(
			planID,
			t.Title,
			desc,
			i+1,
		)

		_, err = s.taskRepo.Create(ctx, task)
		if err != nil {
			_ = s.planRepo.UpdateGenerationStatus(ctx, planID, domain.GenerationFailed)
			return
		}

		tasks = append(tasks, task)
	}

	testingTask := domain.NewTask(planID, "Пройти тестирование", nil, len(tasks)+1)

	_, err = s.taskRepo.Create(ctx, testingTask)
	if err != nil {
		_ = s.planRepo.UpdateGenerationStatus(ctx, planID, domain.GenerationFailed)
		return
	}

	tasks = append(tasks, testingTask)

	plan.Title = generated.Title
	plan.Description = description

	if err := s.generateTestForPlan(ctx, plan, tasks); err != nil {
		_ = s.planRepo.UpdateGenerationStatus(ctx, planID, domain.GenerationFailed)
		return
	}

	_ = s.planRepo.UpdateGenerationStatus(
		ctx,
		planID,
		domain.GenerationReady,
	)
}
