package ai

import (
	"context"
	"encoding/json"
)

func (s *AIService) ExtractSkills(ctx context.Context, input ExtractSkillsInput) ([]SkillCandidate, error) {
	prompt := buildSkillPrompt(input)
	answer, err := s.client.Generate(ctx, prompt)
	if err != nil {
		return nil, err
	}

	var skills []SkillCandidate

	if err := json.Unmarshal([]byte(answer), &skills); err != nil {
		return nil, ErrInvalidAIResponse
	}
	
	return normalizeSkills(skills), nil
}
