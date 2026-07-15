package ai

import (
	"context"
	"encoding/json"
	"strings"
)

func (s *AIService) GenerateTest(ctx context.Context, input GenerateTestInput) (*GeneratedTest, error) {
	prompt := buildTestPrompt(input)
	answer, err := s.client.Generate(ctx, prompt)
	if err != nil {
		return nil, err
	}

	answer = extractJSON(strings.TrimSpace(answer))

	var generated GeneratedTest

	err = json.Unmarshal([]byte(answer), &generated)
	if err != nil {
		return nil, ErrInvalidAIResponse
	}

	if len(generated.Questions) != 10 {
		return nil, ErrInvalidAIResponse
	}

	return &generated, nil
}
