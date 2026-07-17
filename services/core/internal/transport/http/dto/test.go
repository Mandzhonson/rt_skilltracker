package dto

type GenerateTestRequest struct {
	PlanID string `json:"plan_id" binding:"required,uuid"`
}

type GenerateTestResponse struct {
	TestID string `json:"test_id"`
}

type StartTestResponse struct {
	TestID    string             `json:"test_id"`
	Questions []QuestionResponse `json:"questions"`
}

type QuestionResponse struct {
	ID       string   `json:"id"`
	Question string   `json:"question"`
	Options  []string `json:"options"`
}

type SubmitTestRequest struct {
	Answers []AnswerRequest `json:"answers"`
}

type AnswerRequest struct {
	QuestionID string `json:"question_id"`
	Answer     string `json:"answer"`
}

type SubmitTestResponse struct {
	Score  int  `json:"score"`
	Total  int  `json:"total"`
	Passed bool `json:"passed"`
}

