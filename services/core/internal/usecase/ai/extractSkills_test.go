package ai

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mockAIClient struct {
	generateFunc func(ctx context.Context, prompt string) (string, error)
}

func (m *mockAIClient) Generate(ctx context.Context, prompt string) (string, error) {
	if m.generateFunc != nil {
		return m.generateFunc(ctx, prompt)
	}
	return `[
		{"name":"go","category":"Programming Language","description":"Proficient in Go"},
		{"name":"docker","category":"DevOps","description":"Containerization"}
	]`, nil
}

func TestAiService_ExtractSkills(t *testing.T) {
	type mockBehavior func(
		client *mockAIClient,
	)

	testTable := []struct {
		name string

		input ExtractSkillsInput

		mockBehavior mockBehavior

		expectedSkills []SkillCandidate
		expectedErr    error
	}{
		{
			name: "Успешное извлечение навыков",
			input: ExtractSkillsInput{
				PlanTitle:       "Go Developer Plan",
				PlanDescription: "Learn Go programming language",
				Tasks:           []string{"Write Go code", "Build microservices"},
				ExistingSkills:  []string{"Python", "Docker"},
			},

			mockBehavior: func(
				client *mockAIClient,
			) {
				skills := []SkillCandidate{
					{
						Name:        "go",
						Category:    "Programming Language",
						Description: "Proficient in Go",
					},
					{
						Name:        "docker",
						Category:    "DevOps",
						Description: "Containerization",
					},
				}

				skillsJSON, _ := json.Marshal(skills)

				client.generateFunc = func(ctx context.Context, prompt string) (string, error) {
					return string(skillsJSON), nil
				}
			},

			expectedSkills: []SkillCandidate{
				{
					Name:        "Go",
					Category:    "Programming Language",
					Description: "Proficient in Go",
				},
				{
					Name:        "docker",
					Category:    "DevOps",
					Description: "Containerization",
				},
			},
			expectedErr: nil,
		},
		{
			name: "Успешное извлечение навыков (без нормализации)",
			input: ExtractSkillsInput{
				PlanTitle:       "Test Plan",
				PlanDescription: "Test Description",
				Tasks:           []string{"Task 1", "Task 2"},
				ExistingSkills:  []string{},
			},

			mockBehavior: func(
				client *mockAIClient,
			) {
				skills := []SkillCandidate{
					{
						Name:        "python",
						Category:    "Programming Language",
						Description: "Python programming",
					},
				}

				skillsJSON, _ := json.Marshal(skills)

				client.generateFunc = func(ctx context.Context, prompt string) (string, error) {
					return string(skillsJSON), nil
				}
			},

			expectedSkills: []SkillCandidate{
				{
					Name:        "python",
					Category:    "Programming Language",
					Description: "Python programming",
				},
			},
			expectedErr: nil,
		},
		{
			name: "Ошибка при генерации от AI",
			input: ExtractSkillsInput{
				PlanTitle:       "Test Plan",
				PlanDescription: "Test description",
				Tasks:           []string{"Task 1"},
				ExistingSkills:  []string{},
			},

			mockBehavior: func(
				client *mockAIClient,
			) {
				client.generateFunc = func(ctx context.Context, prompt string) (string, error) {
					return "", errors.New("AI service error")
				}
			},

			expectedSkills: nil,
			expectedErr:    errors.New("AI service error"),
		},
		{
			name: "Невалидный JSON ответ от AI",
			input: ExtractSkillsInput{
				PlanTitle:       "Test Plan",
				PlanDescription: "Test description",
				Tasks:           []string{"Task 1"},
				ExistingSkills:  []string{},
			},

			mockBehavior: func(
				client *mockAIClient,
			) {
				client.generateFunc = func(ctx context.Context, prompt string) (string, error) {
					return "invalid json", nil
				}
			},

			expectedSkills: nil,
			expectedErr:    ErrInvalidAIResponse,
		},
		{
			name: "Пустой ответ от AI",
			input: ExtractSkillsInput{
				PlanTitle:       "Test Plan",
				PlanDescription: "Test description",
				Tasks:           []string{"Task 1"},
				ExistingSkills:  []string{},
			},

			mockBehavior: func(
				client *mockAIClient,
			) {
				client.generateFunc = func(ctx context.Context, prompt string) (string, error) {
					return "", nil
				}
			},

			expectedSkills: nil,
			expectedErr:    ErrInvalidAIResponse,
		},
		{
			name: "Ответ с пустым массивом навыков",
			input: ExtractSkillsInput{
				PlanTitle:       "Test Plan",
				PlanDescription: "Test description",
				Tasks:           []string{"Task 1"},
				ExistingSkills:  []string{},
			},

			mockBehavior: func(
				client *mockAIClient,
			) {
				client.generateFunc = func(ctx context.Context, prompt string) (string, error) {
					return "[]", nil
				}
			},

			expectedSkills: []SkillCandidate{},
			expectedErr:    nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			client := &mockAIClient{}

			testCase.mockBehavior(client)

			src := NewAiService(client)

			skills, err := src.ExtractSkills(
				context.Background(),
				testCase.input,
			)

			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedErr, err)
				assert.Nil(t, skills)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(testCase.expectedSkills), len(skills))
				for i, expected := range testCase.expectedSkills {
					assert.Equal(t, expected.Name, skills[i].Name)
					assert.Equal(t, expected.Category, skills[i].Category)
					assert.Equal(t, expected.Description, skills[i].Description)
				}
			}
		})
	}
}
