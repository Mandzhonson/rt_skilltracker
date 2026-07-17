package test

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/postgres"
	mock_postgres "core_service/internal/repository/postgres/mocks"
	"core_service/internal/usecase/task"
	"core_service/internal/usecase/user"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTestService_GetForManager(t *testing.T) {
	type mockBehavior func(
		testRepo *mock_postgres.MockTestRepository,
		planRepo *mock_postgres.MockPlanRepository,
	)

	managerID := uuid.New()
	otherManagerID := uuid.New()
	planID := uuid.New()
	testID := uuid.New()
	question1ID := uuid.New()
	question2ID := uuid.New()

	testTable := []struct {
		name string

		managerID uuid.UUID
		planID    uuid.UUID

		mockBehavior mockBehavior

		expectedView *domain.TestView
		expectedErr  error
	}{
		{
			name:      "Успешное получение теста для менеджера",
			managerID: managerID,
			planID:    planID,

			mockBehavior: func(
				testRepo *mock_postgres.MockTestRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
				planRepo.EXPECT().
					ManagerOwnsPlan(gomock.Any(), managerID, planID).
					Return(true, nil)

				test := &domain.Test{
					ID:     testID,
					PlanID: planID,
				}

				questions := []*domain.Question{
					{
						ID:           question1ID,
						PlanID:       planID,
						QuestionText: "What is Go?",
						OptionA:      "Programming language",
						OptionB:      "Framework",
						OptionC:      "Library",
						OptionD:      "Tool",
					},
					{
						ID:           question2ID,
						PlanID:       planID,
						QuestionText: "What is Docker?",
						OptionA:      "Containerization",
						OptionB:      "Programming language",
						OptionC:      "Database",
						OptionD:      "Framework",
					},
				}

				testRepo.EXPECT().
					GetByPlanID(gomock.Any(), planID).
					Return(test, nil)

				testRepo.EXPECT().
					GetQuestions(gomock.Any(), testID).
					Return(questions, nil)
			},

			expectedView: &domain.TestView{
				ID: testID,
				Questions: []domain.QuestionView{
					{
						ID:   question1ID,
						Text: "What is Go?",
						Options: []string{
							"Programming language",
							"Framework",
							"Library",
							"Tool",
						},
					},
					{
						ID:   question2ID,
						Text: "What is Docker?",
						Options: []string{
							"Containerization",
							"Programming language",
							"Database",
							"Framework",
						},
					},
				},
			},
			expectedErr: nil,
		},
		{
			name:      "Неверный ID менеджера (пустой UUID)",
			managerID: uuid.Nil,
			planID:    planID,

			mockBehavior: func(
				testRepo *mock_postgres.MockTestRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
			},

			expectedView: nil,
			expectedErr:  ErrInvalidUserID,
		},
		{
			name:      "Неверный ID плана (пустой UUID)",
			managerID: managerID,
			planID:    uuid.Nil,

			mockBehavior: func(
				testRepo *mock_postgres.MockTestRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
			},

			expectedView: nil,
			expectedErr:  ErrInvalidPlanID,
		},
		{
			name:      "Менеджер не владеет планом",
			managerID: otherManagerID,
			planID:    planID,

			mockBehavior: func(
				testRepo *mock_postgres.MockTestRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
				planRepo.EXPECT().
					ManagerOwnsPlan(gomock.Any(), otherManagerID, planID).
					Return(false, nil)
			},

			expectedView: nil,
			expectedErr:  user.ErrForbidden,
		},
		{
			name:      "Ошибка при проверке прав менеджера",
			managerID: managerID,
			planID:    planID,

			mockBehavior: func(
				testRepo *mock_postgres.MockTestRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
				planRepo.EXPECT().
					ManagerOwnsPlan(gomock.Any(), managerID, planID).
					Return(false, assert.AnError)
			},

			expectedView: nil,
			expectedErr:  assert.AnError,
		},
		{
			name:      "Тест не найден",
			managerID: managerID,
			planID:    planID,

			mockBehavior: func(
				testRepo *mock_postgres.MockTestRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
				planRepo.EXPECT().
					ManagerOwnsPlan(gomock.Any(), managerID, planID).
					Return(true, nil)

				testRepo.EXPECT().
					GetByPlanID(gomock.Any(), planID).
					Return(nil, postgres.ErrTestNotFound)
			},

			expectedView: nil,
			expectedErr:  ErrTestNotFound,
		},
		{
			name:      "Ошибка при получении вопросов",
			managerID: managerID,
			planID:    planID,

			mockBehavior: func(
				testRepo *mock_postgres.MockTestRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
				planRepo.EXPECT().
					ManagerOwnsPlan(gomock.Any(), managerID, planID).
					Return(true, nil)

				test := &domain.Test{
					ID:     testID,
					PlanID: planID,
				}

				testRepo.EXPECT().
					GetByPlanID(gomock.Any(), planID).
					Return(test, nil)

				testRepo.EXPECT().
					GetQuestions(gomock.Any(), testID).
					Return(nil, assert.AnError)
			},

			expectedView: nil,
			expectedErr:  assert.AnError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			testRepo := mock_postgres.NewMockTestRepository(ctrl)
			planRepo := mock_postgres.NewMockPlanRepository(ctrl)
			taskRepo := mock_postgres.NewMockTaskRepository(ctrl)

			testCase.mockBehavior(testRepo, planRepo)

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

			view, err := src.GetForManager(
				context.Background(),
				testCase.managerID,
				testCase.planID,
			)

			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedErr, err)
				assert.Nil(t, view)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, view)
				assert.Equal(t, testCase.expectedView.ID, view.ID)
				assert.Equal(t, len(testCase.expectedView.Questions), len(view.Questions))
				if len(view.Questions) > 0 {
					assert.Equal(t, testCase.expectedView.Questions[0].ID, view.Questions[0].ID)
					assert.Equal(t, testCase.expectedView.Questions[0].Text, view.Questions[0].Text)
					assert.Equal(t, len(testCase.expectedView.Questions[0].Options), len(view.Questions[0].Options))
				}
			}
		})
	}
}
