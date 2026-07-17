package test

import (
	"context"
	"core_service/internal/domain"

	"github.com/google/uuid"
)

func (s *testService) getTestView(ctx context.Context, planID uuid.UUID) (*domain.TestView, error) {

	if planID == uuid.Nil {
		return nil, ErrInvalidPlanID
	}

	testEntity, err := s.testRepo.GetByPlanID(ctx, planID)
	if err != nil {
		return nil, ErrTestNotFound
	}

	questions, err := s.testRepo.GetQuestions(ctx, testEntity.ID)
	if err != nil {
		return nil, err
	}

	result := make([]domain.QuestionView, 0, len(questions))
	for _, q := range questions {
		result = append(result, domain.QuestionView{
			ID:   q.ID,
			Text: q.QuestionText,
			Options: []string{
				q.OptionA,
				q.OptionB,
				q.OptionC,
				q.OptionD,
			},
		})
	}
	return &domain.TestView{
		ID:        testEntity.ID,
		Questions: result,
	}, nil
}
