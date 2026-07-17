package domain

import (
	"time"

	"github.com/google/uuid"
)

type Question struct {
	ID            uuid.UUID
	PlanID        uuid.UUID
	QuestionText  string
	OptionA       string
	OptionB       string
	OptionC       string
	OptionD       string
	CorrectOption string
	AIGenerated   bool
	CreatedAt     time.Time
}

func NewQuestion(planID uuid.UUID, question, a, b, c, d, correct string) *Question {
	return &Question{
		PlanID:        planID,
		QuestionText:  question,
		OptionA:       a,
		OptionB:       b,
		OptionC:       c,
		OptionD:       d,
		CorrectOption: correct,
		AIGenerated:   true,
	}
}
