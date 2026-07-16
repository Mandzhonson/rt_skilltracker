package plan

import "context"

type MockAIClient struct{}

func (m *MockAIClient) Generate(ctx context.Context, prompt string) (string, error) {
	return "", nil
}
