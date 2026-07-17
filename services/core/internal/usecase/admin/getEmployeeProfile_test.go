package admin

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/postgres"
	mock_postgres "core_service/internal/repository/postgres/mocks"
	"core_service/internal/usecase/user"
	"errors"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestAdminService_GetEmployeeProfileForAdmin(t *testing.T) {
	type mockBehavior func(
		userRepo *mock_postgres.MockUserRepository,
		planRepo *mock_postgres.MockPlanRepository,
		skillRepo *mock_postgres.MockSkillRepository,
	)

	employeeID := uuid.New()
	planID1 := uuid.New()
	planID2 := uuid.New()
	skillID1 := uuid.New()
	skillID2 := uuid.New()

	testTable := []struct {
		name string

		employeeID uuid.UUID

		mockBehavior mockBehavior

		expectedProfile *domain.EmployeeProfile
		expectedErr     error
	}{
		{
			name:       "Успешное получение профиля сотрудника",
			employeeID: employeeID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				planRepo *mock_postgres.MockPlanRepository,
				skillRepo *mock_postgres.MockSkillRepository,
			) {
				employee := &domain.User{
					ID:        employeeID,
					Email:     "employee@mail.ru",
					FirstName: "Test",
					LastName:  "Employee",
					Role:      domain.RoleEmployee,
					Position:  "Developer",
				}

				plans := []*domain.Plan{
					{
						ID:     planID1,
						Title:  "Plan 1",
						Status: domain.PlanActive,
					},
					{
						ID:     planID2,
						Title:  "Plan 2",
						Status: domain.PlanCompleted,
					},
				}

				skills := []*domain.Skill{
					{
						ID:   skillID1,
						Name: "Go",
					},
					{
						ID:   skillID2,
						Name: "React",
					},
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(employee, nil)

				planRepo.EXPECT().
					ListAllByEmployeeID(gomock.Any(), employeeID).
					Return(plans, nil)

				skillRepo.EXPECT().
					ListByUserID(gomock.Any(), employeeID).
					Return(skills, nil)
			},

			expectedProfile: &domain.EmployeeProfile{
				User: &domain.User{
					ID:        employeeID,
					Email:     "employee@mail.ru",
					FirstName: "Test",
					LastName:  "Employee",
					Role:      domain.RoleEmployee,
					Position:  "Developer",
				},
				Plans: []*domain.Plan{
					{
						ID:     planID1,
						Title:  "Plan 1",
						Status: domain.PlanActive,
					},
					{
						ID:     planID2,
						Title:  "Plan 2",
						Status: domain.PlanCompleted,
					},
				},
				Skills: []*domain.Skill{
					{
						ID:   skillID1,
						Name: "Go",
					},
					{
						ID:   skillID2,
						Name: "React",
					},
				},
			},
			expectedErr: nil,
		},
		{
			name:       "Сотрудник не найден",
			employeeID: employeeID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				planRepo *mock_postgres.MockPlanRepository,
				skillRepo *mock_postgres.MockSkillRepository,
			) {
				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(nil, postgres.ErrUserNotFound)
			},

			expectedProfile: nil,
			expectedErr:     user.ErrUserNotFound,
		},
		{
			name:       "Пользователь не является сотрудником",
			employeeID: employeeID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				planRepo *mock_postgres.MockPlanRepository,
				skillRepo *mock_postgres.MockSkillRepository,
			) {
				manager := &domain.User{
					ID:        employeeID,
					Email:     "manager@mail.ru",
					FirstName: "Test",
					LastName:  "Manager",
					Role:      domain.RoleManager,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(manager, nil)
			},

			expectedProfile: nil,
			expectedErr:     ErrInvalidEmployee,
		},
		{
			name:       "Ошибка при получении планов",
			employeeID: employeeID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				planRepo *mock_postgres.MockPlanRepository,
				skillRepo *mock_postgres.MockSkillRepository,
			) {
				employee := &domain.User{
					ID:    employeeID,
					Email: "employee@mail.ru",
					Role:  domain.RoleEmployee,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(employee, nil)

				planRepo.EXPECT().
					ListAllByEmployeeID(gomock.Any(), employeeID).
					Return(nil, errors.New("database error"))
			},

			expectedProfile: nil,
			expectedErr:     errors.New("list plans: database error"),
		},
		{
			name:       "Ошибка при получении навыков",
			employeeID: employeeID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				planRepo *mock_postgres.MockPlanRepository,
				skillRepo *mock_postgres.MockSkillRepository,
			) {
				employee := &domain.User{
					ID:    employeeID,
					Email: "employee@mail.ru",
					Role:  domain.RoleEmployee,
				}

				plans := []*domain.Plan{
					{
						ID:     planID1,
						Title:  "Plan 1",
						Status: domain.PlanActive,
					},
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), employeeID).
					Return(employee, nil)

				planRepo.EXPECT().
					ListAllByEmployeeID(gomock.Any(), employeeID).
					Return(plans, nil)

				skillRepo.EXPECT().
					ListByUserID(gomock.Any(), employeeID).
					Return(nil, errors.New("database error"))
			},

			expectedProfile: nil,
			expectedErr:     errors.New("list skills: database error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mock_postgres.NewMockUserRepository(ctrl)
			planRepo := mock_postgres.NewMockPlanRepository(ctrl)
			skillRepo := mock_postgres.NewMockSkillRepository(ctrl)

			testCase.mockBehavior(userRepo, planRepo, skillRepo)

			src := NewAdminService(
				userRepo,
				planRepo,
				skillRepo,
				nil, // storage
			)

			profile, err := src.GetEmployeeProfileForAdmin(
				context.Background(),
				testCase.employeeID,
			)

			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedErr.Error(), err.Error())
				assert.Nil(t, profile)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, profile)
				assert.Equal(t, testCase.expectedProfile.User.ID, profile.User.ID)
				assert.Equal(t, testCase.expectedProfile.User.Email, profile.User.Email)
				assert.Equal(t, len(testCase.expectedProfile.Plans), len(profile.Plans))
				assert.Equal(t, len(testCase.expectedProfile.Skills), len(profile.Skills))
				if len(profile.Plans) > 0 {
					assert.Equal(t, testCase.expectedProfile.Plans[0].Title, profile.Plans[0].Title)
				}
				if len(profile.Skills) > 0 {
					assert.Equal(t, testCase.expectedProfile.Skills[0].Name, profile.Skills[0].Name)
				}
			}
		})
	}
}
