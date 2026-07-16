package plan

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/postgres"
	mock_postgres "core_service/internal/repository/postgres/mocks"
	"core_service/internal/usecase/ai"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockAIClient struct{}

func (m *mockAIClient) Generate(ctx context.Context, prompt string) (string, error) {
	return `[
		{"name":"go","category":"Programming Language","description":"Proficient in Go"}
	]`, nil
}

func TestPlanService_GenerateSkillsIfCompleted(t *testing.T) {
	type mockBehavior func(
		planRepo *mock_postgres.MockPlanRepository,
		userRepo *mock_postgres.MockUserRepository,
		taskRepo *mock_postgres.MockTaskRepository,
		skillRepo *mock_postgres.MockSkillRepository,
		testRepo *mock_postgres.MockTestRepository,
	)

	planID := uuid.New()
	employeeID := uuid.New()
	managerID := uuid.New()
	skillID := uuid.New()
	description := "Test Description"

	testTable := []struct {
		name string

		planID uuid.UUID

		mockBehavior mockBehavior

		expectedErr error
	}{
		{
			name:   "Успешная генерация навыков (план завершен)",
			planID: planID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				userRepo *mock_postgres.MockUserRepository,
				taskRepo *mock_postgres.MockTaskRepository,
				skillRepo *mock_postgres.MockSkillRepository,
				testRepo *mock_postgres.MockTestRepository,
			) {
				plan := &domain.Plan{
					ID:               planID,
					EmployeeID:       employeeID,
					CreatedBy:        managerID,
					Title:            "Test Plan",
					Description:      &description,
					GenerationStatus: domain.GenerationReady,
					CreationType:     domain.CreationAI,
					Progress:         100,
					Status:           domain.PlanCompleted,
				}

				tasks := []*domain.Task{
					{
						ID:     uuid.New(),
						PlanID: planID,
						Title:  "Task 1",
						Status: domain.TaskDone,
					},
					{
						ID:     uuid.New(),
						PlanID: planID,
						Title:  "Task 2",
						Status: domain.TaskDone,
					},
				}

				existingSkills := []*domain.Skill{
					{
						ID:   uuid.New(),
						Name: "Python",
					},
				}

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)

				taskRepo.EXPECT().
					ListByPlanID(gomock.Any(), planID).
					Return(tasks, nil)

				skillRepo.EXPECT().
					ListByUserID(gomock.Any(), employeeID).
					Return(existingSkills, nil)

				skillRepo.EXPECT().
					GetByName(gomock.Any(), "Go").
					Return(nil, errors.New("skill not found"))

				skillRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(skillID, nil)

				skillRepo.EXPECT().
					AttachToUser(gomock.Any(), employeeID, skillID, planID).
					Return(nil)
			},

			expectedErr: nil,
		},
		{
			name:   "План не завершен (прогресс < 100)",
			planID: planID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				userRepo *mock_postgres.MockUserRepository,
				taskRepo *mock_postgres.MockTaskRepository,
				skillRepo *mock_postgres.MockSkillRepository,
				testRepo *mock_postgres.MockTestRepository,
			) {
				plan := &domain.Plan{
					ID:               planID,
					EmployeeID:       employeeID,
					CreatedBy:        managerID,
					Title:            "Test Plan",
					Description:      &description,
					GenerationStatus: domain.GenerationReady,
					Progress:         50,
					Status:           domain.PlanActive,
				}

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)
			},

			expectedErr: nil,
		},
		{
			name:   "План не найден",
			planID: planID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				userRepo *mock_postgres.MockUserRepository,
				taskRepo *mock_postgres.MockTaskRepository,
				skillRepo *mock_postgres.MockSkillRepository,
				testRepo *mock_postgres.MockTestRepository,
			) {
				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(nil, postgres.ErrPlanNotFound)
			},

			expectedErr: postgres.ErrPlanNotFound,
		},
		{
			name:   "Ошибка при получении задач",
			planID: planID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				userRepo *mock_postgres.MockUserRepository,
				taskRepo *mock_postgres.MockTaskRepository,
				skillRepo *mock_postgres.MockSkillRepository,
				testRepo *mock_postgres.MockTestRepository,
			) {
				plan := &domain.Plan{
					ID:               planID,
					EmployeeID:       employeeID,
					CreatedBy:        managerID,
					Title:            "Test Plan",
					Description:      &description,
					GenerationStatus: domain.GenerationReady,
					Progress:         100,
					Status:           domain.PlanCompleted,
				}

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)

				taskRepo.EXPECT().
					ListByPlanID(gomock.Any(), planID).
					Return(nil, errors.New("database error"))
			},

			expectedErr: errors.New("database error"),
		},
		{
			name:   "Ошибка при получении существующих навыков",
			planID: planID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				userRepo *mock_postgres.MockUserRepository,
				taskRepo *mock_postgres.MockTaskRepository,
				skillRepo *mock_postgres.MockSkillRepository,
				testRepo *mock_postgres.MockTestRepository,
			) {
				plan := &domain.Plan{
					ID:               planID,
					EmployeeID:       employeeID,
					CreatedBy:        managerID,
					Title:            "Test Plan",
					Description:      &description,
					GenerationStatus: domain.GenerationReady,
					Progress:         100,
					Status:           domain.PlanCompleted,
				}

				tasks := []*domain.Task{
					{
						ID:     uuid.New(),
						PlanID: planID,
						Title:  "Task 1",
						Status: domain.TaskDone,
					},
				}

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)

				taskRepo.EXPECT().
					ListByPlanID(gomock.Any(), planID).
					Return(tasks, nil)

				skillRepo.EXPECT().
					ListByUserID(gomock.Any(), employeeID).
					Return(nil, errors.New("database error"))
			},

			expectedErr: errors.New("database error"),
		},
		{
			name:   "Ошибка при создании навыка",
			planID: planID,

			mockBehavior: func(
				planRepo *mock_postgres.MockPlanRepository,
				userRepo *mock_postgres.MockUserRepository,
				taskRepo *mock_postgres.MockTaskRepository,
				skillRepo *mock_postgres.MockSkillRepository,
				testRepo *mock_postgres.MockTestRepository,
			) {
				plan := &domain.Plan{
					ID:               planID,
					EmployeeID:       employeeID,
					CreatedBy:        managerID,
					Title:            "Test Plan",
					Description:      &description,
					GenerationStatus: domain.GenerationReady,
					Progress:         100,
					Status:           domain.PlanCompleted,
				}

				tasks := []*domain.Task{
					{
						ID:     uuid.New(),
						PlanID: planID,
						Title:  "Task 1",
						Status: domain.TaskDone,
					},
				}

				existingSkills := []*domain.Skill{}

				planRepo.EXPECT().
					GetByID(gomock.Any(), planID).
					Return(plan, nil)

				taskRepo.EXPECT().
					ListByPlanID(gomock.Any(), planID).
					Return(tasks, nil)

				skillRepo.EXPECT().
					ListByUserID(gomock.Any(), employeeID).
					Return(existingSkills, nil)

				skillRepo.EXPECT().
					GetByName(gomock.Any(), "Go").
					Return(nil, errors.New("skill not found"))

				skillRepo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(uuid.Nil, errors.New("create skill error"))
			},

			expectedErr: errors.New("create skill error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			planRepo := mock_postgres.NewMockPlanRepository(ctrl)
			userRepo := mock_postgres.NewMockUserRepository(ctrl)
			taskRepo := mock_postgres.NewMockTaskRepository(ctrl)
			skillRepo := mock_postgres.NewMockSkillRepository(ctrl)
			testRepo := mock_postgres.NewMockTestRepository(ctrl)

			testCase.mockBehavior(planRepo, userRepo, taskRepo, skillRepo, testRepo)

			mockClient := &mockAIClient{}
			aiService := ai.NewAiService(mockClient)

			src := NewPlanService(
				planRepo,
				userRepo,
				taskRepo,
				skillRepo,
				testRepo,
				aiService,
			)

			err := src.GenerateSkillsIfCompleted(
				context.Background(),
				testCase.planID,
			)

			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
