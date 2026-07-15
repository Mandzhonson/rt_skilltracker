package plan

import (
	"context"

	"core_service/internal/domain"
	"core_service/internal/usecase/ai"
)

func (s *planService) generateTestForPlan(ctx context.Context, plan *domain.Plan, tasks []*domain.Task) error {
	generatedTasks := make([]ai.GeneratedTask, 0, len(tasks))

	for _, task := range tasks {

		generatedTasks = append(generatedTasks, ai.GeneratedTask{
			Title:       task.Title,
			Description: deref(task.Description),
		},
		)
	}
	employee, err := s.userRepo.GetById(ctx, plan.EmployeeID)
	if err != nil {
		return err
	}
	generatedTest, err := s.aiService.GenerateTest(ctx, ai.GenerateTestInput{
		PlanTitle:       plan.Title,
		PlanDescription: deref(plan.Description),
		Tasks:           generatedTasks,
		Position:        employee.Position,
	},
	)

	if err != nil {
		return err
	}

	if len(generatedTest.Questions) == 0 {
		return ErrTestGenerationFailed
	}

	test := domain.NewTest(plan.ID)

	questions := make(
		[]*domain.Question,
		0,
		len(generatedTest.Questions),
	)

	for _, q := range generatedTest.Questions {

		correct := normalizeCorrectOption(q.CorrectOption)

		if correct == "" {
			continue
		}

		questions = append(
			questions,
			domain.NewQuestion(
				plan.ID,
				q.Question,
				q.OptionA,
				q.OptionB,
				q.OptionC,
				q.OptionD,
				correct,
			),
		)
	}

	if len(questions) == 0 {
		return ErrTestGenerationFailed
	}

	_, err = s.testRepo.CreateWithQuestions(
		ctx,
		test,
		questions,
	)

	return err
}
