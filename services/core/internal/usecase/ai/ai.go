package ai

import (
	"core_service/internal/clients/ollama"
	"errors"
)

var (
	ErrInvalidTopic      = errors.New("invalid topic")
	ErrInvalidPlan       = errors.New("invalid plan")
	ErrInvalidAIResponse = errors.New("invalid ai response")
)

type AIService struct {
	client *ollama.Client
}

func NewAIService(client *ollama.Client) *AIService {
	return &AIService{
		client: client,
	}
}
