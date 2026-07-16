package admin

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/postgres"
	mock_postgres "core_service/internal/repository/postgres/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAdminService_ListEmployeesByManager(t *testing.T) {
	type mockBehavior func(
		userRepo *mock_postgres.MockUserRepository,
	)

	managerID := uuid.New()
	employeeID1 := uuid.New()
	employeeID2 := uuid.New()

	testTable := []struct {
		name string

		managerID uuid.UUID

		mockBehavior mockBehavior

		expectedEmployees []*domain.User
		expectedErr       error
	}{
		{
			name:      "Успешное получение списка сотрудников",
			managerID: managerID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				manager := &domain.User{
					ID:    managerID,
					Email: "manager@mail.ru",
					Role:  domain.RoleManager,
				}

				employees := []*domain.User{
					{
						ID:        employeeID1,
						Email:     "employee1@mail.ru",
						FirstName: "Employee",
						LastName:  "One",
						Role:      domain.RoleEmployee,
						ManagerID: &managerID,
					},
					{
						ID:        employeeID2,
						Email:     "employee2@mail.ru",
						FirstName: "Employee",
						LastName:  "Two",
						Role:      domain.RoleEmployee,
						ManagerID: &managerID,
					},
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), managerID).
					Return(manager, nil)

				userRepo.EXPECT().
					ListEmployeesByManager(gomock.Any(), managerID).
					Return(employees, nil)
			},

			expectedEmployees: []*domain.User{
				{
					ID:        employeeID1,
					Email:     "employee1@mail.ru",
					FirstName: "Employee",
					LastName:  "One",
					Role:      domain.RoleEmployee,
					ManagerID: &managerID,
				},
				{
					ID:        employeeID2,
					Email:     "employee2@mail.ru",
					FirstName: "Employee",
					LastName:  "Two",
					Role:      domain.RoleEmployee,
					ManagerID: &managerID,
				},
			},
			expectedErr: nil,
		},
		{
			name:      "Менеджер не найден",
			managerID: managerID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				userRepo.EXPECT().
					GetById(gomock.Any(), managerID).
					Return(nil, postgres.ErrUserNotFound)
			},

			expectedEmployees: nil,
			expectedErr:       postgres.ErrUserNotFound,
		},
		{
			name:      "Пользователь не является менеджером",
			managerID: managerID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				employee := &domain.User{
					ID:    managerID,
					Email: "employee@mail.ru",
					Role:  domain.RoleEmployee,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), managerID).
					Return(employee, nil)
			},

			expectedEmployees: nil,
			expectedErr:       ErrInvalidManager,
		},
		{
			name:      "Пустой список сотрудников",
			managerID: managerID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
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
					ListEmployeesByManager(gomock.Any(), managerID).
					Return([]*domain.User{}, nil)
			},

			expectedEmployees: []*domain.User{},
			expectedErr:       nil,
		},
		{
			name:      "Ошибка при получении списка сотрудников",
			managerID: managerID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
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
					ListEmployeesByManager(gomock.Any(), managerID).
					Return(nil, assert.AnError)
			},

			expectedEmployees: nil,
			expectedErr:       assert.AnError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mock_postgres.NewMockUserRepository(ctrl)

			testCase.mockBehavior(userRepo)

			src := NewAdminService(
				userRepo,
				nil, // planRepo
				nil, // skillRepo
				nil, // storage
			)

			employees, err := src.ListEmployeesByManager(
				context.Background(),
				testCase.managerID,
			)

			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedErr, err)
				assert.Nil(t, employees)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(testCase.expectedEmployees), len(employees))
				if len(employees) > 0 {
					assert.Equal(t, testCase.expectedEmployees[0].ID, employees[0].ID)
					assert.Equal(t, testCase.expectedEmployees[0].Email, employees[0].Email)
				}
			}
		})
	}
}
