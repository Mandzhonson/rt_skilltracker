package user

import (
	"context"
	"core_service/internal/domain"
	mock_minio "core_service/internal/repository/minio/mocks"
	"core_service/internal/repository/postgres"
	mock_postgres "core_service/internal/repository/postgres/mocks"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserService_GetEmployeeProfile(t *testing.T) {
	type mockBehavior func(
		userRepo *mock_postgres.MockUserRepository,
		skillRepo *mock_postgres.MockSkillRepository,
		planRepo *mock_postgres.MockPlanRepository,
	)

	managerID := uuid.New()
	employeeID := uuid.New()
	otherManagerID := uuid.New()
	skill1ID := uuid.New()
	skill2ID := uuid.New()
	plan1ID := uuid.New()
	plan2ID := uuid.New()
	now := time.Now()

	testTable := []struct {
		name string

		managerID  uuid.UUID
		employeeID uuid.UUID

		mockBehavior mockBehavior

		expectedProfile *EmployeeProfile
		expectedErr     error
	}{
		{
			name:       "Успешное получение профиля сотрудника",
			managerID:  managerID,
			employeeID: employeeID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				skillRepo *mock_postgres.MockSkillRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
				manager := &domain.User{
					ID:    managerID,
					Email: "manager@mail.ru",
					Role:  domain.RoleManager,
				}

				employee := &domain.User{
					ID:        employeeID,
					Email:     "employee@mail.ru",
					FirstName: "Test",
					LastName:  "Employee",
					Role:      domain.RoleEmployee,
					Position:  "Developer",
					ManagerID: &managerID,
				}

				skills := []*domain.Skill{
					{
						ID:   skill1ID,
						Name: "Go",
					},
					{
						ID:   skill2ID,
						Name: "Docker",
					},
				}

				plans := []*domain.Plan{
					{
						ID:               plan1ID,
						EmployeeID:       employeeID,
						CreatedBy:        managerID,
						Title:            "Plan 1",
						GenerationStatus: domain.GenerationReady,
						CreationType:     domain.CreationAI,
						Progress:         50,
						Status:           domain.PlanActive,
						CreatedAt:        now,
						UpdatedAt:        now,
					},
					{
						ID:               plan2ID,
						EmployeeID:       employeeID,
						CreatedBy:        managerID,
						Title:            "Plan 2",
						GenerationStatus: domain.GenerationReady,
						CreationType:     domain.CreationAI,
						Progress:         80,
						Status:           domain.PlanActive,
						CreatedAt:        now,
						UpdatedAt:        now,
					},
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), managerID).
					Return(manager, nil)

				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(employee, nil)

				skillRepo.EXPECT().
					ListByUserID(gomock.Any(), employeeID).
					Return(skills, nil)

				planRepo.EXPECT().
					ListAllByEmployeeID(gomock.Any(), employeeID).
					Return(plans, nil)
			},

			expectedProfile: &EmployeeProfile{
				User: &domain.User{
					ID:        employeeID,
					Email:     "employee@mail.ru",
					FirstName: "Test",
					LastName:  "Employee",
					Role:      domain.RoleEmployee,
					Position:  "Developer",
					ManagerID: &managerID,
				},
				Skills: []*domain.Skill{
					{
						ID:   skill1ID,
						Name: "Go",
					},
					{
						ID:   skill2ID,
						Name: "Docker",
					},
				},
				Plans: []*domain.Plan{
					{
						ID:               plan1ID,
						EmployeeID:       employeeID,
						CreatedBy:        managerID,
						Title:            "Plan 1",
						GenerationStatus: domain.GenerationReady,
						CreationType:     domain.CreationAI,
						Progress:         50,
						Status:           domain.PlanActive,
						CreatedAt:        now,
						UpdatedAt:        now,
					},
					{
						ID:               plan2ID,
						EmployeeID:       employeeID,
						CreatedBy:        managerID,
						Title:            "Plan 2",
						GenerationStatus: domain.GenerationReady,
						CreationType:     domain.CreationAI,
						Progress:         80,
						Status:           domain.PlanActive,
						CreatedAt:        now,
						UpdatedAt:        now,
					},
				},
			},
			expectedErr: nil,
		},
		{
			name:       "Менеджер не найден",
			managerID:  managerID,
			employeeID: employeeID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				skillRepo *mock_postgres.MockSkillRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
				userRepo.EXPECT().
					GetById(gomock.Any(), managerID).
					Return(nil, postgres.ErrUserNotFound)
			},

			expectedProfile: nil,
			expectedErr:     ErrUserNotFound,
		},
		{
			name:       "Пользователь не является менеджером",
			managerID:  managerID,
			employeeID: employeeID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				skillRepo *mock_postgres.MockSkillRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
				manager := &domain.User{
					ID:    managerID,
					Email: "employee@mail.ru",
					Role:  domain.RoleEmployee,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), managerID).
					Return(manager, nil)
			},

			expectedProfile: nil,
			expectedErr:     ErrNotManager,
		},
		{
			name:       "Сотрудник не найден",
			managerID:  managerID,
			employeeID: employeeID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				skillRepo *mock_postgres.MockSkillRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
				manager := &domain.User{
					ID:    managerID,
					Email: "manager@mail.ru",
					Role:  domain.RoleManager,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), managerID).
					Return(manager, nil)

				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(nil, postgres.ErrUserNotFound)
			},

			expectedProfile: nil,
			expectedErr:     ErrUserNotFound,
		},
		{
			name:       "Сотрудник не закреплен за менеджером",
			managerID:  managerID,
			employeeID: employeeID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				skillRepo *mock_postgres.MockSkillRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
				manager := &domain.User{
					ID:    managerID,
					Email: "manager@mail.ru",
					Role:  domain.RoleManager,
				}

				employee := &domain.User{
					ID:        employeeID,
					Email:     "employee@mail.ru",
					FirstName: "Test",
					LastName:  "Employee",
					Role:      domain.RoleEmployee,
					ManagerID: nil,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), managerID).
					Return(manager, nil)

				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(employee, nil)
			},

			expectedProfile: nil,
			expectedErr:     ErrForbidden,
		},
		{
			name:       "Сотрудник закреплен за другим менеджером",
			managerID:  otherManagerID,
			employeeID: employeeID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				skillRepo *mock_postgres.MockSkillRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
				manager := &domain.User{
					ID:    otherManagerID,
					Email: "manager@mail.ru",
					Role:  domain.RoleManager,
				}

				employee := &domain.User{
					ID:        employeeID,
					Email:     "employee@mail.ru",
					FirstName: "Test",
					LastName:  "Employee",
					Role:      domain.RoleEmployee,
					ManagerID: &managerID,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), otherManagerID).
					Return(manager, nil)

				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(employee, nil)
			},

			expectedProfile: nil,
			expectedErr:     ErrForbidden,
		},
		{
			name:       "Ошибка при получении навыков",
			managerID:  managerID,
			employeeID: employeeID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				skillRepo *mock_postgres.MockSkillRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
				manager := &domain.User{
					ID:    managerID,
					Email: "manager@mail.ru",
					Role:  domain.RoleManager,
				}

				employee := &domain.User{
					ID:        employeeID,
					Email:     "employee@mail.ru",
					FirstName: "Test",
					LastName:  "Employee",
					Role:      domain.RoleEmployee,
					ManagerID: &managerID,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), managerID).
					Return(manager, nil)

				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(employee, nil)

				skillRepo.EXPECT().
					ListByUserID(gomock.Any(), employeeID).
					Return(nil, errors.New("database error"))
			},

			expectedProfile: nil,
			expectedErr:     errors.New("database error"),
		},
		{
			name:       "Ошибка при получении планов",
			managerID:  managerID,
			employeeID: employeeID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				skillRepo *mock_postgres.MockSkillRepository,
				planRepo *mock_postgres.MockPlanRepository,
			) {
				manager := &domain.User{
					ID:    managerID,
					Email: "manager@mail.ru",
					Role:  domain.RoleManager,
				}

				employee := &domain.User{
					ID:        employeeID,
					Email:     "employee@mail.ru",
					FirstName: "Test",
					LastName:  "Employee",
					Role:      domain.RoleEmployee,
					ManagerID: &managerID,
				}

				skills := []*domain.Skill{
					{
						ID:   skill1ID,
						Name: "Go",
					},
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), managerID).
					Return(manager, nil)

				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(employee, nil)

				skillRepo.EXPECT().
					ListByUserID(gomock.Any(), employeeID).
					Return(skills, nil)

				planRepo.EXPECT().
					ListAllByEmployeeID(gomock.Any(), employeeID).
					Return(nil, errors.New("database error"))
			},

			expectedProfile: nil,
			expectedErr:     errors.New("database error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mock_postgres.NewMockUserRepository(ctrl)
			skillRepo := mock_postgres.NewMockSkillRepository(ctrl)
			planRepo := mock_postgres.NewMockPlanRepository(ctrl)
			storage := mock_minio.NewMockStorage(ctrl)

			testCase.mockBehavior(userRepo, skillRepo, planRepo)

			src := NewUserService(
				userRepo,
				storage,
				skillRepo,
				planRepo,
			)

			profile, err := src.GetEmployeeProfile(
				context.Background(),
				testCase.managerID,
				testCase.employeeID,
			)

			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedErr, err)
				assert.Nil(t, profile)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, profile)
				assert.Equal(t, testCase.expectedProfile.User.ID, profile.User.ID)
				assert.Equal(t, testCase.expectedProfile.User.Email, profile.User.Email)
				assert.Equal(t, len(testCase.expectedProfile.Skills), len(profile.Skills))
				assert.Equal(t, len(testCase.expectedProfile.Plans), len(profile.Plans))
			}
		})
	}
}
