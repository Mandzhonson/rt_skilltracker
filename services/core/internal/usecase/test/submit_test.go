package test

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/postgres"
	mock_postgres "core_service/internal/repository/postgres/mocks"
	"core_service/internal/usecase/task"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTestService_Submit(t *testing.T) {
	type mockBehavior func(
		testRepo *mock_postgres.MockTestRepository,
		planRepo *mock_postgres.MockPlanRepository,
		taskRepo *mock_postgres.MockTaskRepository,
	)

	planID := uuid.New()
	userID := uuid.New()
	testID := uuid.New()
	question1ID := uuid.New()
	question2ID := uuid.New()
	attemptID := uuid.New()

	testTable := []struct {
		name string

		input SubmitTestInput

		mockBehavior mockBehavior

		expectedResult *domain.TestResult
		expectedErr    error
	}{
		{
			name: "Успешная сдача теста (пройден)",
			input: SubmitTestInput{
				PlanID: planID,
				UserID: userID,
				Answers: []AnswerInput{
					{
						QuestionID: question1ID,
						Answer:     "A",
					},
					{
						QuestionID: question2ID,
						Answer:     "B",
					},
				},
			},

			mockBehavior: func(
				testRepo *mock_postgres.MockTestRepository,
				planRepo *mock_postgres.MockPlanRepository,
				taskRepo *mock_postgres.MockTaskRepository,
			) {
				test := &domain.Test{
					ID:     testID,
					PlanID: planID,
				}

				questions := []*domain.Question{
					{
						ID:            question1ID,
						PlanID:        planID,
						QuestionText:  "Q1",
						OptionA:       "A1",
						OptionB:       "B1",
						OptionC:       "C1",
						OptionD:       "D1",
						CorrectOption: "A",
					},
					{
						ID:            question2ID,
						PlanID:        planID,
						QuestionText:  "Q2",
						OptionA:       "A2",
						OptionB:       "B2",
						OptionC:       "C2",
						OptionD:       "D2",
						CorrectOption: "B",
					},
				}

				plan := &domain.Plan{
					ID:         planID,
					EmployeeID: userID,
					Title:      "Test Plan",
					Status:     domain.PlanActive,
				}

				testRepo.EXPECT().
					GetByPlanID(gomock.Any(), planID).
					Return(test, nil)

				testRepo.EXPECT().
					GetQuestions(gomock.Any(), testID).
					Return(questions, nil)

				testRepo.EXPECT().
					CreateAttempt(gomock.Any(), gomock.Any()).
					Return(attemptID, nil)

				testRepo.EXPECT().
					CreateAnswers(gomock.Any(), gomock.Any()).
					Return(nil)

				taskRepo.EXPECT().
					CompleteTestingTask(gomock.Any(), planID).
					Return(nil)

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)

				planRepo.EXPECT().
					RecalculateProgress(gomock.Any(), planID).
					Return(100, nil)

				planRepo.EXPECT().
					RecalculateProgress(gomock.Any(), planID).
					Return(100, nil)
			},

			expectedResult: &domain.TestResult{
				Score:  100,
				Total:  2,
				Passed: true,
			},
			expectedErr: nil,
		},
		{
			name: "Успешная сдача теста (не пройден)",
			input: SubmitTestInput{
				PlanID: planID,
				UserID: userID,
				Answers: []AnswerInput{
					{
						QuestionID: question1ID,
						Answer:     "B",
					},
					{
						QuestionID: question2ID,
						Answer:     "C",
					},
				},
			},

			mockBehavior: func(
				testRepo *mock_postgres.MockTestRepository,
				planRepo *mock_postgres.MockPlanRepository,
				taskRepo *mock_postgres.MockTaskRepository,
			) {
				test := &domain.Test{
					ID:     testID,
					PlanID: planID,
				}

				questions := []*domain.Question{
					{
						ID:            question1ID,
						PlanID:        planID,
						QuestionText:  "Q1",
						OptionA:       "A1",
						OptionB:       "B1",
						OptionC:       "C1",
						OptionD:       "D1",
						CorrectOption: "A",
					},
					{
						ID:            question2ID,
						PlanID:        planID,
						QuestionText:  "Q2",
						OptionA:       "A2",
						OptionB:       "B2",
						OptionC:       "C2",
						OptionD:       "D2",
						CorrectOption: "B",
					},
				}

				testRepo.EXPECT().
					GetByPlanID(gomock.Any(), planID).
					Return(test, nil)

				testRepo.EXPECT().
					GetQuestions(gomock.Any(), testID).
					Return(questions, nil)

				testRepo.EXPECT().
					CreateAttempt(gomock.Any(), gomock.Any()).
					Return(attemptID, nil)

				testRepo.EXPECT().
					CreateAnswers(gomock.Any(), gomock.Any()).
					Return(nil)

			},

			expectedResult: &domain.TestResult{
				Score:  0,
				Total:  2,
				Passed: false,
			},
			expectedErr: nil,
		},
		{
			name: "Неверный ID плана (пустой UUID)",
			input: SubmitTestInput{
				PlanID:  uuid.Nil,
				UserID:  userID,
				Answers: []AnswerInput{},
			},

			mockBehavior: func(
				testRepo *mock_postgres.MockTestRepository,
				planRepo *mock_postgres.MockPlanRepository,
				taskRepo *mock_postgres.MockTaskRepository,
			) {
			},

			expectedResult: nil,
			expectedErr:    ErrInvalidPlanID,
		},
		{
			name: "Неверный ID пользователя (пустой UUID)",
			input: SubmitTestInput{
				PlanID:  planID,
				UserID:  uuid.Nil,
				Answers: []AnswerInput{},
			},

			mockBehavior: func(
				testRepo *mock_postgres.MockTestRepository,
				planRepo *mock_postgres.MockPlanRepository,
				taskRepo *mock_postgres.MockTaskRepository,
			) {
			},

			expectedResult: nil,
			expectedErr:    ErrInvalidUserID,
		},
		{
			name: "Пустые ответы",
			input: SubmitTestInput{
				PlanID:  planID,
				UserID:  userID,
				Answers: []AnswerInput{},
			},

			mockBehavior: func(
				testRepo *mock_postgres.MockTestRepository,
				planRepo *mock_postgres.MockPlanRepository,
				taskRepo *mock_postgres.MockTaskRepository,
			) {
			},

			expectedResult: nil,
			expectedErr:    ErrInvalidAnswers,
		},
		{
			name: "Тест не найден",
			input: SubmitTestInput{
				PlanID: planID,
				UserID: userID,
				Answers: []AnswerInput{
					{
						QuestionID: question1ID,
						Answer:     "A",
					},
				},
			},

			mockBehavior: func(
				testRepo *mock_postgres.MockTestRepository,
				planRepo *mock_postgres.MockPlanRepository,
				taskRepo *mock_postgres.MockTaskRepository,
			) {
				testRepo.EXPECT().
					GetByPlanID(gomock.Any(), planID).
					Return(nil, postgres.ErrTestNotFound)
			},

			expectedResult: nil,
			expectedErr:    ErrTestNotFound,
		},
		{
			name: "Ошибка при получении вопросов",
			input: SubmitTestInput{
				PlanID: planID,
				UserID: userID,
				Answers: []AnswerInput{
					{
						QuestionID: question1ID,
						Answer:     "A",
					},
				},
			},

			mockBehavior: func(
				testRepo *mock_postgres.MockTestRepository,
				planRepo *mock_postgres.MockPlanRepository,
				taskRepo *mock_postgres.MockTaskRepository,
			) {
				test := &domain.Test{
					ID:     testID,
					PlanID: planID,
				}

				testRepo.EXPECT().
					GetByPlanID(gomock.Any(), planID).
					Return(test, nil)

				testRepo.EXPECT().
					GetQuestions(gomock.Any(), testID).
					Return(nil, errors.New("get questions error"))
			},

			expectedResult: nil,
			expectedErr:    errors.New("get questions error"),
		},
		{
			name: "Ошибка при создании попытки",
			input: SubmitTestInput{
				PlanID: planID,
				UserID: userID,
				Answers: []AnswerInput{
					{
						QuestionID: question1ID,
						Answer:     "A",
					},
				},
			},

			mockBehavior: func(
				testRepo *mock_postgres.MockTestRepository,
				planRepo *mock_postgres.MockPlanRepository,
				taskRepo *mock_postgres.MockTaskRepository,
			) {
				test := &domain.Test{
					ID:     testID,
					PlanID: planID,
				}

				questions := []*domain.Question{
					{
						ID:            question1ID,
						PlanID:        planID,
						QuestionText:  "Q1",
						OptionA:       "A1",
						OptionB:       "B1",
						OptionC:       "C1",
						OptionD:       "D1",
						CorrectOption: "A",
					},
				}

				testRepo.EXPECT().
					GetByPlanID(gomock.Any(), planID).
					Return(test, nil)

				testRepo.EXPECT().
					GetQuestions(gomock.Any(), testID).
					Return(questions, nil)

				testRepo.EXPECT().
					CreateAttempt(gomock.Any(), gomock.Any()).
					Return(uuid.Nil, errors.New("create attempt error"))
			},

			expectedResult: nil,
			expectedErr:    errors.New("create attempt error"),
		},
		{
			name: "Ошибка при создании ответов",
			input: SubmitTestInput{
				PlanID: planID,
				UserID: userID,
				Answers: []AnswerInput{
					{
						QuestionID: question1ID,
						Answer:     "A",
					},
				},
			},

			mockBehavior: func(
				testRepo *mock_postgres.MockTestRepository,
				planRepo *mock_postgres.MockPlanRepository,
				taskRepo *mock_postgres.MockTaskRepository,
			) {
				test := &domain.Test{
					ID:     testID,
					PlanID: planID,
				}

				questions := []*domain.Question{
					{
						ID:            question1ID,
						PlanID:        planID,
						QuestionText:  "Q1",
						OptionA:       "A1",
						OptionB:       "B1",
						OptionC:       "C1",
						OptionD:       "D1",
						CorrectOption: "A",
					},
				}

				testRepo.EXPECT().
					GetByPlanID(gomock.Any(), planID).
					Return(test, nil)

				testRepo.EXPECT().
					GetQuestions(gomock.Any(), testID).
					Return(questions, nil)

				testRepo.EXPECT().
					CreateAttempt(gomock.Any(), gomock.Any()).
					Return(attemptID, nil)

				testRepo.EXPECT().
					CreateAnswers(gomock.Any(), gomock.Any()).
					Return(errors.New("create answers error"))
			},

			expectedResult: nil,
			expectedErr:    errors.New("create answers error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			testRepo := mock_postgres.NewMockTestRepository(ctrl)
			planRepo := mock_postgres.NewMockPlanRepository(ctrl)
			taskRepo := mock_postgres.NewMockTaskRepository(ctrl)

			testCase.mockBehavior(testRepo, planRepo, taskRepo)

			mockPlanCompletion := &mockPlanCompletionService{}

			taskService := task.NewTaskService(
				taskRepo,
				planRepo,
				mockPlanCompletion,
			)

			src := NewTestService(
				testRepo,
				*taskService,
				planRepo,
				mockPlanCompletion,
			)

			result, err := src.Submit(
				context.Background(),
				testCase.input,
			)

			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedErr, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, testCase.expectedResult.Score, result.Score)
				assert.Equal(t, testCase.expectedResult.Total, result.Total)
				assert.Equal(t, testCase.expectedResult.Passed, result.Passed)
			}
		})
	}
}
