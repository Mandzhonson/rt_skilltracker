package ai

import (
	"context"
	"errors"
)

var (
	ErrInvalidTopic      = errors.New("invalid topic")
	ErrInvalidPlan       = errors.New("invalid plan")
	ErrInvalidAIResponse = errors.New("invalid ai response")
)

type AIClient interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

type AiService struct {
	client AIClient
}

func NewAiService(client AIClient) *AiService {
	return &AiService{
		client: client,
	}
}
