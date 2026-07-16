package admin

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/postgres"
	mock_postgres "core_service/internal/repository/postgres/mocks"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
)

func TestAdminService_AssignManager(t *testing.T) {
	type mockBehavior func(
		userRepo *mock_postgres.MockUserRepository,
	)

	userID := uuid.New()
	managerID := uuid.New()
	otherManagerID := uuid.New()
	topManagerID := uuid.New()

	testTable := []struct {
		name string

		input AssignManagerInput

		mockBehavior mockBehavior

		expectedErr error
	}{
		{
			name: "Успешное назначение менеджера (нет цепочки)",
			input: AssignManagerInput{
				UserID:    userID,
				ManagerID: managerID,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				user := &domain.User{
					ID:    userID,
					Email: "employee@mail.ru",
					Role:  domain.RoleEmployee,
				}

				manager := &domain.User{
					ID:        managerID,
					Email:     "manager@mail.ru",
					Role:      domain.RoleManager,
					ManagerID: nil,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(user, nil)

				userRepo.EXPECT().
					GetById(gomock.Any(), managerID).
					Return(manager, nil)

				userRepo.EXPECT().
					GetById(gomock.Any(), gomock.Any()).
					Return(&domain.User{
						ID:        uuid.New(),
						Email:     "any@mail.ru",
						Role:      domain.RoleManager,
						ManagerID: nil,
					}, nil).
					AnyTimes()

				userRepo.EXPECT().
					AssignManager(gomock.Any(), userID, managerID).
					Return(nil)
			},

			expectedErr: nil,
		},
		{
			name: "Успешное назначение менеджера (с цепочкой)",
			input: AssignManagerInput{
				UserID:    userID,
				ManagerID: managerID,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				user := &domain.User{
					ID:    userID,
					Email: "employee@mail.ru",
					Role:  domain.RoleEmployee,
				}

				manager := &domain.User{
					ID:        managerID,
					Email:     "manager@mail.ru",
					Role:      domain.RoleManager,
					ManagerID: &topManagerID,
				}

				topManager := &domain.User{
					ID:        topManagerID,
					Email:     "topmanager@mail.ru",
					Role:      domain.RoleManager,
					ManagerID: nil,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(user, nil)

				userRepo.EXPECT().
					GetById(gomock.Any(), managerID).
					Return(manager, nil)

				userRepo.EXPECT().
					GetById(gomock.Any(), managerID).
					Return(manager, nil).
					AnyTimes()

				userRepo.EXPECT().
					GetById(gomock.Any(), topManagerID).
					Return(topManager, nil).
					AnyTimes()

				userRepo.EXPECT().
					AssignManager(gomock.Any(), userID, managerID).
					Return(nil)
			},

			expectedErr: nil,
		},
		{
			name: "Попытка назначить самого себя",
			input: AssignManagerInput{
				UserID:    userID,
				ManagerID: userID,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
			},

			expectedErr: ErrAssignYourself,
		},
		{
			name: "Сотрудник не найден",
			input: AssignManagerInput{
				UserID:    userID,
				ManagerID: managerID,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(nil, postgres.ErrUserNotFound)
			},

			expectedErr: postgres.ErrUserNotFound,
		},
		{
			name: "Попытка назначить менеджера админу",
			input: AssignManagerInput{
				UserID:    userID,
				ManagerID: managerID,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				admin := &domain.User{
					ID:    userID,
					Email: "admin@mail.ru",
					Role:  domain.RoleAdmin,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(admin, nil)
			},

			expectedErr: ErrInvalidManager,
		},
		{
			name: "Менеджер не найден",
			input: AssignManagerInput{
				UserID:    userID,
				ManagerID: managerID,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				user := &domain.User{
					ID:    userID,
					Email: "employee@mail.ru",
					Role:  domain.RoleEmployee,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(user, nil)

				userRepo.EXPECT().
					GetById(gomock.Any(), managerID).
					Return(nil, postgres.ErrUserNotFound)
			},

			expectedErr: postgres.ErrUserNotFound,
		},
		{
			name: "Попытка назначить менеджером не-менеджера",
			input: AssignManagerInput{
				UserID:    userID,
				ManagerID: otherManagerID,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				user := &domain.User{
					ID:    userID,
					Email: "employee@mail.ru",
					Role:  domain.RoleEmployee,
				}

				notManager := &domain.User{
					ID:    otherManagerID,
					Email: "notmanager@mail.ru",
					Role:  domain.RoleEmployee,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(user, nil)

				userRepo.EXPECT().
					GetById(gomock.Any(), otherManagerID).
					Return(notManager, nil)
			},

			expectedErr: ErrInvalidManager,
		},
		{
			name: "Ошибка валидации иерархии (циклическая зависимость)",
			input: AssignManagerInput{
				UserID:    userID,
				ManagerID: managerID,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				user := &domain.User{
					ID:    userID,
					Email: "employee@mail.ru",
					Role:  domain.RoleEmployee,
				}

				manager := &domain.User{
					ID:        managerID,
					Email:     "manager@mail.ru",
					Role:      domain.RoleManager,
					ManagerID: &userID,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(user, nil)

				userRepo.EXPECT().
					GetById(gomock.Any(), managerID).
					Return(manager, nil)

				userRepo.EXPECT().
					GetById(gomock.Any(), managerID).
					Return(manager, nil).
					AnyTimes()

				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(user, nil).
					AnyTimes()

			},

			expectedErr: ErrManagerCycle,
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

			err := src.AssignManager(
				context.Background(),
				testCase.input,
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
