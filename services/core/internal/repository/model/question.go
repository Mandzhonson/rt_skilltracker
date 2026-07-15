package model

import (
	"time"

	"github.com/google/uuid"
)

type QuestionModel struct {
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
