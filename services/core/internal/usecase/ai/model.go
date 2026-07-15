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

type ExtractSkillsInput struct {
	PlanTitle       string
	PlanDescription string
	Tasks           []string
	ExistingSkills  []string
}

type ExtractSkillsResult struct {
	Skills []SkillCandidate
}

type SkillCandidate struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
}
