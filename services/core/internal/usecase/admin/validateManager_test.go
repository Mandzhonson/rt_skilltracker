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

func TestAdminService_validateManagerHierarchy(t *testing.T) {
	type mockBehavior func(
		userRepo *mock_postgres.MockUserRepository,
	)

	userID := uuid.New()
	managerID := uuid.New()
	topManagerID := uuid.New()

	testTable := []struct {
		name string

		userID    uuid.UUID
		managerID uuid.UUID

		mockBehavior mockBehavior

		expectedErr error
	}{
		{
			name:      "Успешная валидация (менеджер без начальника)",
			userID:    userID,
			managerID: managerID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				manager := &domain.User{
					ID:        managerID,
					Email:     "manager@mail.ru",
					Role:      domain.RoleManager,
					ManagerID: nil,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), managerID).
					Return(manager, nil)
			},

			expectedErr: nil,
		},
		{
			name:      "Успешная валидация (цепочка менеджеров без цикла)",
			userID:    userID,
			managerID: managerID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
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
					GetById(gomock.Any(), managerID).
					Return(manager, nil)

				userRepo.EXPECT().
					GetById(gomock.Any(), topManagerID).
					Return(topManager, nil)
			},

			expectedErr: nil,
		},
		{
			name:      "Циклическая зависимость (менеджер указывает на пользователя)",
			userID:    userID,
			managerID: managerID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				manager := &domain.User{
					ID:        managerID,
					Email:     "manager@mail.ru",
					Role:      domain.RoleManager,
					ManagerID: &userID,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), managerID).
					Return(manager, nil)
			},

			expectedErr: ErrManagerCycle,
		},
		{
			name:      "Циклическая зависимость (глубокая цепочка)",
			userID:    userID,
			managerID: managerID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
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
					ManagerID: &userID,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), managerID).
					Return(manager, nil)

				userRepo.EXPECT().
					GetById(gomock.Any(), topManagerID).
					Return(topManager, nil)

			},

			expectedErr: ErrManagerCycle,
		},
		{
			name:      "Менеджер не найден",
			userID:    userID,
			managerID: managerID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				userRepo.EXPECT().
					GetById(gomock.Any(), managerID).
					Return(nil, postgres.ErrUserNotFound)
			},

			expectedErr: postgres.ErrUserNotFound,
		},
		{
			name:      "Промежуточный менеджер не найден",
			userID:    userID,
			managerID: managerID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				manager := &domain.User{
					ID:        managerID,
					Email:     "manager@mail.ru",
					Role:      domain.RoleManager,
					ManagerID: &topManagerID,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), managerID).
					Return(manager, nil)

				userRepo.EXPECT().
					GetById(gomock.Any(), topManagerID).
					Return(nil, postgres.ErrUserNotFound)
			},

			expectedErr: postgres.ErrUserNotFound,
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

			err := src.validateManagerHierarchy(
				context.Background(),
				testCase.userID,
				testCase.managerID,
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
