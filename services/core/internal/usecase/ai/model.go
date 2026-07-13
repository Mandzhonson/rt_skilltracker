package ai

type GeneratePlanInput struct {
	Topic          string
	Description    string
	EmployeeLevel  string
	ExistingSkills []string
	TargetRole     string
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
