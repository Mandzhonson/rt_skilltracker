package ai

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAiService_GenerateTest(t *testing.T) {
	type mockBehavior func(
		client *mockAIClient,
	)

	testTable := []struct {
		name string

		input GenerateTestInput

		mockBehavior mockBehavior

		expectedTest *GeneratedTest
		expectedErr  error
	}{
		{
			name: "Успешная генерация теста",
			input: GenerateTestInput{
				PlanTitle:       "Go Programming",
				PlanDescription: "Learn Go from scratch",
				Position:        "Junior Developer",
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

			mockBehavior: func(
				client *mockAIClient,
			) {
				test := GeneratedTest{
					Questions: []GeneratedQuestion{
						{
							Question:      "What is Go?",
							OptionA:       "Programming language",
							OptionB:       "Framework",
							OptionC:       "Library",
							OptionD:       "Tool",
							CorrectOption: "A",
						},
						{
							Question:      "What is a goroutine?",
							OptionA:       "A lightweight thread",
							OptionB:       "A function",
							OptionC:       "A package",
							OptionD:       "A type",
							CorrectOption: "A",
						},
						{
							Question:      "What is a channel in Go?",
							OptionA:       "A communication mechanism",
							OptionB:       "A data type",
							OptionC:       "A function",
							OptionD:       "A package",
							CorrectOption: "A",
						},
						{
							Question:      "What is a pointer in Go?",
							OptionA:       "A variable that stores memory address",
							OptionB:       "A function",
							OptionC:       "A package",
							OptionD:       "A type",
							CorrectOption: "A",
						},
						{
							Question:      "What is a slice in Go?",
							OptionA:       "A dynamic array",
							OptionB:       "A fixed array",
							OptionC:       "A map",
							OptionD:       "A struct",
							CorrectOption: "A",
						},
						{
							Question:      "What is a map in Go?",
							OptionA:       "A key-value store",
							OptionB:       "A list",
							OptionC:       "A set",
							OptionD:       "A queue",
							CorrectOption: "A",
						},
						{
							Question:      "What is an interface in Go?",
							OptionA:       "A set of method signatures",
							OptionB:       "A struct",
							OptionC:       "A function",
							OptionD:       "A package",
							CorrectOption: "A",
						},
						{
							Question:      "What is a struct in Go?",
							OptionA:       "A collection of fields",
							OptionB:       "A function",
							OptionC:       "A package",
							OptionD:       "A type",
							CorrectOption: "A",
						},
						{
							Question:      "What is a method in Go?",
							OptionA:       "A function with receiver",
							OptionB:       "A struct",
							OptionC:       "A package",
							OptionD:       "A type",
							CorrectOption: "A",
						},
						{
							Question:      "What is a package in Go?",
							OptionA:       "A collection of Go files",
							OptionB:       "A function",
							OptionC:       "A type",
							OptionD:       "A struct",
							CorrectOption: "A",
						},
					},
				}

				testJSON, _ := json.Marshal(test)

				client.generateFunc = func(ctx context.Context, prompt string) (string, error) {
					return string(testJSON), nil
				}
			},

			expectedTest: &GeneratedTest{
				Questions: []GeneratedQuestion{
					{
						Question:      "What is Go?",
						OptionA:       "Programming language",
						OptionB:       "Framework",
						OptionC:       "Library",
						OptionD:       "Tool",
						CorrectOption: "A",
					},
					{
						Question:      "What is a goroutine?",
						OptionA:       "A lightweight thread",
						OptionB:       "A function",
						OptionC:       "A package",
						OptionD:       "A type",
						CorrectOption: "A",
					},
					{
						Question:      "What is a channel in Go?",
						OptionA:       "A communication mechanism",
						OptionB:       "A data type",
						OptionC:       "A function",
						OptionD:       "A package",
						CorrectOption: "A",
					},
					{
						Question:      "What is a pointer in Go?",
						OptionA:       "A variable that stores memory address",
						OptionB:       "A function",
						OptionC:       "A package",
						OptionD:       "A type",
						CorrectOption: "A",
					},
					{
						Question:      "What is a slice in Go?",
						OptionA:       "A dynamic array",
						OptionB:       "A fixed array",
						OptionC:       "A map",
						OptionD:       "A struct",
						CorrectOption: "A",
					},
					{
						Question:      "What is a map in Go?",
						OptionA:       "A key-value store",
						OptionB:       "A list",
						OptionC:       "A set",
						OptionD:       "A queue",
						CorrectOption: "A",
					},
					{
						Question:      "What is an interface in Go?",
						OptionA:       "A set of method signatures",
						OptionB:       "A struct",
						OptionC:       "A function",
						OptionD:       "A package",
						CorrectOption: "A",
					},
					{
						Question:      "What is a struct in Go?",
						OptionA:       "A collection of fields",
						OptionB:       "A function",
						OptionC:       "A package",
						OptionD:       "A type",
						CorrectOption: "A",
					},
					{
						Question:      "What is a method in Go?",
						OptionA:       "A function with receiver",
						OptionB:       "A struct",
						OptionC:       "A package",
						OptionD:       "A type",
						CorrectOption: "A",
					},
					{
						Question:      "What is a package in Go?",
						OptionA:       "A collection of Go files",
						OptionB:       "A function",
						OptionC:       "A type",
						OptionD:       "A struct",
						CorrectOption: "A",
					},
				},
			},
			expectedErr: nil,
		},
		{
			name: "Ошибка при генерации от AI",
			input: GenerateTestInput{
				PlanTitle:       "Go Programming",
				PlanDescription: "Learn Go",
				Position:        "Developer",
				Tasks:           []GeneratedTask{},
			},

			mockBehavior: func(
				client *mockAIClient,
			) {
				client.generateFunc = func(ctx context.Context, prompt string) (string, error) {
					return "", errors.New("AI service error")
				}
			},

			expectedTest: nil,
			expectedErr:  errors.New("AI service error"),
		},
		{
			name: "Невалидный JSON ответ от AI",
			input: GenerateTestInput{
				PlanTitle:       "Go Programming",
				PlanDescription: "Learn Go",
				Position:        "Developer",
				Tasks:           []GeneratedTask{},
			},

			mockBehavior: func(
				client *mockAIClient,
			) {
				client.generateFunc = func(ctx context.Context, prompt string) (string, error) {
					return "invalid json", nil
				}
			},

			expectedTest: nil,
			expectedErr:  ErrInvalidAIResponse,
		},
		{
			name: "Неверное количество вопросов (не 10)",
			input: GenerateTestInput{
				PlanTitle:       "Go Programming",
				PlanDescription: "Learn Go",
				Position:        "Developer",
				Tasks:           []GeneratedTask{},
			},

			mockBehavior: func(
				client *mockAIClient,
			) {
				test := GeneratedTest{
					Questions: []GeneratedQuestion{
						{
							Question:      "What is Go?",
							OptionA:       "Programming language",
							OptionB:       "Framework",
							OptionC:       "Library",
							OptionD:       "Tool",
							CorrectOption: "A",
						},
					},
				}

				testJSON, _ := json.Marshal(test)

				client.generateFunc = func(ctx context.Context, prompt string) (string, error) {
					return string(testJSON), nil
				}
			},

			expectedTest: nil,
			expectedErr:  ErrInvalidAIResponse,
		},
		{
			name: "Пустой ответ от AI",
			input: GenerateTestInput{
				PlanTitle:       "Go Programming",
				PlanDescription: "Learn Go",
				Position:        "Developer",
				Tasks:           []GeneratedTask{},
			},

			mockBehavior: func(
				client *mockAIClient,
			) {
				client.generateFunc = func(ctx context.Context, prompt string) (string, error) {
					return "", nil
				}
			},

			expectedTest: nil,
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

			test, err := src.GenerateTest(
				context.Background(),
				testCase.input,
			)

			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedErr, err)
				assert.Nil(t, test)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, test)
				assert.Equal(t, len(testCase.expectedTest.Questions), len(test.Questions))
				if len(test.Questions) > 0 {
					assert.Equal(t, testCase.expectedTest.Questions[0].Question, test.Questions[0].Question)
					assert.Equal(t, testCase.expectedTest.Questions[0].OptionA, test.Questions[0].OptionA)
					assert.Equal(t, testCase.expectedTest.Questions[0].OptionB, test.Questions[0].OptionB)
					assert.Equal(t, testCase.expectedTest.Questions[0].OptionC, test.Questions[0].OptionC)
					assert.Equal(t, testCase.expectedTest.Questions[0].OptionD, test.Questions[0].OptionD)
					assert.Equal(t, testCase.expectedTest.Questions[0].CorrectOption, test.Questions[0].CorrectOption)
				}
			}
		})
	}
}
