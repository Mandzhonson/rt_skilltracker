package converter

import (
	"core_service/internal/domain"
	"core_service/internal/repository/model"
)

func ToQuestionEntity(q *model.QuestionModel) *domain.Question {
	return &domain.Question{
		ID:            q.ID,
		PlanID:        q.PlanID,
		QuestionText:  q.QuestionText,
		OptionA:       q.OptionA,
		OptionB:       q.OptionB,
		OptionC:       q.OptionC,
		OptionD:       q.OptionD,
		CorrectOption: q.CorrectOption,
		AIGenerated:   q.AIGenerated,
		CreatedAt:     q.CreatedAt,
	}
}

func ToQuestionModel(q *domain.Question) *model.QuestionModel {
	return &model.QuestionModel{
		ID:            q.ID,
		PlanID:        q.PlanID,
		QuestionText:  q.QuestionText,
		OptionA:       q.OptionA,
		OptionB:       q.OptionB,
		OptionC:       q.OptionC,
		OptionD:       q.OptionD,
		CorrectOption: q.CorrectOption,
		AIGenerated:   q.AIGenerated,
		CreatedAt:     q.CreatedAt,
	}
}
