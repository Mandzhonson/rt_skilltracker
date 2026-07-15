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

type GeneratedQuestion struct {
	Question      string `json:"question"`
	OptionA       string `json:"option_a"`
	OptionB       string `json:"option_b"`
	OptionC       string `json:"option_c"`
	OptionD       string `json:"option_d"`
	CorrectOption string `json:"correct_option"`
}

type GenerateTestInput struct {
	PlanTitle       string
	PlanDescription string
	Tasks           []GeneratedTask
}

type GeneratedTest struct {
	Questions []GeneratedQuestion `json:"questions"`
}
