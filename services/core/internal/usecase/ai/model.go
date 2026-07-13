package ai

type GeneratePlanInput struct {
	Topic          string
	Description    string
	ExistingSkills []string
}

type GeneratedTask struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type GeneratedPlan struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Tasks       []GeneratedTask `json:"tasks"`
}
