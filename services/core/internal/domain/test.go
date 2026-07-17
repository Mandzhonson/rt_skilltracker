package domain

import (
	"time"

	"github.com/google/uuid"
)

type Test struct {
	ID        uuid.UUID
	PlanID    uuid.UUID
	CreatedAt time.Time
}

func NewTest(planID uuid.UUID) *Test {
	return &Test{
		PlanID: planID,
	}
}

type QuestionView struct {
	ID      uuid.UUID
	Text    string
	Options []string
}

type TestView struct {
	ID        uuid.UUID
	Questions []QuestionView
}

type TestQuestion struct {
	ID         uuid.UUID
	TestID     uuid.UUID
	QuestionID uuid.UUID
	Position   int
}

func NewTestQuestion(testID uuid.UUID, questionID uuid.UUID, position int) *TestQuestion {
	return &TestQuestion{
		TestID:     testID,
		QuestionID: questionID,
		Position:   position,
	}
}

type TestAttempt struct {
	ID         uuid.UUID
	TestID     uuid.UUID
	UserID     uuid.UUID
	Score      int
	Total      int
	Passed     bool
	AIFeedback *string
	StartedAt  time.Time
	FinishedAt *time.Time
}

type TestAnswer struct {
	ID             uuid.UUID
	AttemptID      uuid.UUID
	QuestionID     uuid.UUID
	SelectedOption string
	IsCorrect      bool
}

type TestResult struct {
	Score  int
	Total  int
	Passed bool
}
