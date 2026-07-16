package ai

import (
	"context"
	"encoding/json"
	"strings"
)

func (s *AiService) GeneratePlan(ctx context.Context, input GeneratePlanInput) (*GeneratedPlan, error) {
	if strings.TrimSpace(input.Topic) == "" {
		return nil, ErrInvalidTopic
	}
	prompt := buildPlanPrompt(input)
	answer, err := s.client.Generate(ctx, prompt)
	if err != nil {
		return nil, err
	}
	answer = strings.TrimSpace(answer)

	var plan GeneratedPlan

	if err = json.Unmarshal([]byte(answer), &plan); err != nil {
		return nil, ErrInvalidAIResponse
	}
	return &plan, nil
}
