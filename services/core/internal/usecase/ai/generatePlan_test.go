package ai

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAiService_GeneratePlan(t *testing.T) {
	type mockBehavior func(
		client *mockAIClient,
	)

	testTable := []struct {
		name string

		input GeneratePlanInput

		mockBehavior mockBehavior

		expectedPlan *GeneratedPlan
		expectedErr  error
	}{
		{
			name: "Успешная генерация плана",
			input: GeneratePlanInput{
				Topic:          "Go Programming",
				Description:    "Learn Go from scratch",
				Position:       "Junior Developer",
				ExistingSkills: []string{"Python", "Docker"},
			},

			mockBehavior: func(
				client *mockAIClient,
			) {
				plan := GeneratedPlan{
					Title:       "Go Programming Plan",
					Description: "Learn Go from scratch",
					Tasks: []GeneratedTask{
						{
							Title:       "Learn Go basics",
							Description: "Study Go syntax and fundamentals",
						},
						{
							Title:       "Build a REST API",
							Description: "Create a simple REST API in Go",
						},
					},
				}

				planJSON, _ := json.Marshal(plan)

				client.generateFunc = func(ctx context.Context, prompt string) (string, error) {
					return string(planJSON), nil
				}
			},

			expectedPlan: &GeneratedPlan{
				Title:       "Go Programming Plan",
				Description: "Learn Go from scratch",
				Tasks: []GeneratedTask{
					{
						Title:       "Learn Go basics",
						Description: "Study Go syntax and fundamentals",
					},
					{
						Title:       "Build a REST API",
						Description: "Create a simple REST API in Go",
					},
				},
			},
			expectedErr: nil,
		},
		{
			name: "Успешная генерация плана (без описания)",
			input: GeneratePlanInput{
				Topic:          "React",
				Description:    "",
				Position:       "Junior Frontend Developer",
				ExistingSkills: []string{"JavaScript", "HTML", "CSS"},
			},

			mockBehavior: func(
				client *mockAIClient,
			) {
				plan := GeneratedPlan{
					Title:       "React Learning Plan",
					Description: "",
					Tasks: []GeneratedTask{
						{
							Title:       "Learn React basics",
							Description: "Components, props, state",
						},
					},
				}

				planJSON, _ := json.Marshal(plan)

				client.generateFunc = func(ctx context.Context, prompt string) (string, error) {
					return string(planJSON), nil
				}
			},

			expectedPlan: &GeneratedPlan{
				Title:       "React Learning Plan",
				Description: "",
				Tasks: []GeneratedTask{
					{
						Title:       "Learn React basics",
						Description: "Components, props, state",
					},
				},
			},
			expectedErr: nil,
		},
		{
			name: "Успешная генерация плана (с существующими навыками)",
			input: GeneratePlanInput{
				Topic:          "Docker",
				Description:    "Learn containerization",
				Position:       "DevOps Engineer",
				ExistingSkills: []string{"Linux", "Bash", "Go"},
			},

			mockBehavior: func(
				client *mockAIClient,
			) {
				plan := GeneratedPlan{
					Title:       "Docker Mastery",
					Description: "Learn containerization",
					Tasks: []GeneratedTask{
						{
							Title:       "Docker basics",
							Description: "Images, containers, Dockerfile",
						},
					},
				}

				planJSON, _ := json.Marshal(plan)

				client.generateFunc = func(ctx context.Context, prompt string) (string, error) {
					return string(planJSON), nil
				}
			},

			expectedPlan: &GeneratedPlan{
				Title:       "Docker Mastery",
				Description: "Learn containerization",
				Tasks: []GeneratedTask{
					{
						Title:       "Docker basics",
						Description: "Images, containers, Dockerfile",
					},
				},
			},
			expectedErr: nil,
		},
		{
			name: "Пустая тема - ошибка",
			input: GeneratePlanInput{
				Topic:          "   ",
				Description:    "Test",
				Position:       "Developer",
				ExistingSkills: []string{},
			},

			mockBehavior: func(
				client *mockAIClient,
			) {
				// Ничего не ожидаем
			},

			expectedPlan: nil,
			expectedErr:  ErrInvalidTopic,
		},
		{
			name: "Ошибка при генерации от AI",
			input: GeneratePlanInput{
				Topic:          "Go Programming",
				Description:    "Learn Go",
				Position:       "Developer",
				ExistingSkills: []string{},
			},

			mockBehavior: func(
				client *mockAIClient,
			) {
				client.generateFunc = func(ctx context.Context, prompt string) (string, error) {
					return "", errors.New("AI service error")
				}
			},

			expectedPlan: nil,
			expectedErr:  errors.New("AI service error"),
		},
		{
			name: "Невалидный JSON ответ от AI",
			input: GeneratePlanInput{
				Topic:          "Go Programming",
				Description:    "Learn Go",
				Position:       "Developer",
				ExistingSkills: []string{},
			},

			mockBehavior: func(
				client *mockAIClient,
			) {
				client.generateFunc = func(ctx context.Context, prompt string) (string, error) {
					return "invalid json", nil
				}
			},

			expectedPlan: nil,
			expectedErr:  ErrInvalidAIResponse,
		},
		{
			name: "Пустой ответ от AI",
			input: GeneratePlanInput{
				Topic:          "Go Programming",
				Description:    "Learn Go",
				Position:       "Developer",
				ExistingSkills: []string{},
			},

			mockBehavior: func(
				client *mockAIClient,
			) {
				client.generateFunc = func(ctx context.Context, prompt string) (string, error) {
					return "", nil
				}
			},

			expectedPlan: nil,
			expectedErr:  ErrInvalidAIResponse,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			client := &mockAIClient{}

			testCase.mockBehavior(client)

			src := NewAiService(client)

			plan, err := src.GeneratePlan(
				context.Background(),
				testCase.input,
			)

			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedErr, err)
				assert.Nil(t, plan)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, plan)
				assert.Equal(t, testCase.expectedPlan.Title, plan.Title)
				assert.Equal(t, testCase.expectedPlan.Description, plan.Description)
				assert.Equal(t, len(testCase.expectedPlan.Tasks), len(plan.Tasks))
				if len(plan.Tasks) > 0 {
					assert.Equal(t, testCase.expectedPlan.Tasks[0].Title, plan.Tasks[0].Title)
					assert.Equal(t, testCase.expectedPlan.Tasks[0].Description, plan.Tasks[0].Description)
				}
			}
		})
	}
}
