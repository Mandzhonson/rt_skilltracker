package test

import "github.com/google/uuid"

type SubmitTestInput struct {
	UserID uuid.UUID

	PlanID uuid.UUID

	Answers []AnswerInput
}

type AnswerInput struct {
	QuestionID uuid.UUID

	Answer string
}
