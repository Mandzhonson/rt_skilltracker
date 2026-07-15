package test

import (
	"context"

	"core_service/internal/domain"

	"github.com/google/uuid"
)

func (s *testService) Submit(ctx context.Context, input SubmitTestInput) (*domain.TestResult, error) {
	if input.PlanID == uuid.Nil {
		return nil, ErrInvalidPlanID
	}

	if input.UserID == uuid.Nil {
		return nil, ErrInvalidUserID
	}

	if len(input.Answers) == 0 {
		return nil, ErrInvalidAnswers
	}

	testEntity, err := s.testRepo.GetByPlanID(ctx, input.PlanID)
	if err != nil {
		return nil, ErrTestNotFound
	}

	questions, err := s.testRepo.GetQuestions(ctx, testEntity.ID)
	if err != nil {
		return nil, err
	}

	answerMap := make(map[uuid.UUID]string)

	for _, answer := range input.Answers {

		answerMap[answer.QuestionID] = answer.Answer
	}

	correct := 0

	attempt := &domain.TestAttempt{
		TestID: testEntity.ID,
		UserID: input.UserID,
		Total:  len(questions),
	}

	answers := make([]*domain.TestAnswer, 0, len(questions))

	for _, q := range questions {

		selected := answerMap[q.ID]

		isCorrect := selected == q.CorrectOption

		if isCorrect {
			correct++
		}

		answers = append(answers, &domain.TestAnswer{
			QuestionID:     q.ID,
			SelectedOption: selected,
			IsCorrect:      isCorrect,
		},
		)
	}

	score := correct * 100 / len(questions)

	attempt.Score = score
	attempt.Passed = score >= 70

	attemptID, err := s.testRepo.CreateAttempt(ctx, attempt)
	if err != nil {
		return nil, err
	}

	for _, answer := range answers {
		answer.AttemptID = attemptID
	}

	err = s.testRepo.CreateAnswers(ctx, answers)
	if err != nil {
		return nil, err
	}

	if attempt.Passed {

		err = s.taskService.CompleteTestingTask(
			ctx,
			input.PlanID,
			input.UserID,
		)

		if err != nil {
			return nil, err
		}

		progress, err := s.planRepo.RecalculateProgress(
			ctx,
			input.PlanID,
		)

		if err != nil {
			return nil, err
		}

		if progress == 100 {

			go func(planID uuid.UUID) {

				err := s.planCompletionService.GenerateSkillsIfCompleted(
					context.Background(),
					planID,
				)

				if err != nil {
					// TODO slog.Error()
				}

			}(input.PlanID)
		}
	}

	return &domain.TestResult{
		Score:  score,
		Total:  len(questions),
		Passed: attempt.Passed,
	}, nil
}
