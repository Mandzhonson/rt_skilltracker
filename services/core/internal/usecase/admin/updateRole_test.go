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

func TestAdminService_UpdateRole(t *testing.T) {
	type mockBehavior func(
		userRepo *mock_postgres.MockUserRepository,
	)

	userID := uuid.New()
	actorID := uuid.New()

	testTable := []struct {
		name string

		input UpdateRoleInput

		mockBehavior mockBehavior

		expectedErr error
	}{
		{
			name: "Успешное обновление роли сотрудника на менеджера",
			input: UpdateRoleInput{
				ActorID: actorID,
				UserID:  userID,
				Role:    domain.RoleManager,
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
					UpdateRole(gomock.Any(), userID, domain.RoleManager).
					Return(nil)
			},

			expectedErr: nil,
		},
		{
			name: "Успешное обновление роли менеджера на сотрудника (с очисткой назначений)",
			input: UpdateRoleInput{
				ActorID: actorID,
				UserID:  userID,
				Role:    domain.RoleEmployee,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				user := &domain.User{
					ID:    userID,
					Email: "manager@mail.ru",
					Role:  domain.RoleManager,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(user, nil)

				userRepo.EXPECT().
					ClearManagerAssignments(gomock.Any(), userID).
					Return(nil)

				userRepo.EXPECT().
					UpdateRole(gomock.Any(), userID, domain.RoleEmployee).
					Return(nil)
			},

			expectedErr: nil,
		},
		{
			name: "Успешное обновление роли сотрудника на администратора",
			input: UpdateRoleInput{
				ActorID: actorID,
				UserID:  userID,
				Role:    domain.RoleAdmin,
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
					UpdateRole(gomock.Any(), userID, domain.RoleAdmin).
					Return(nil)
			},

			expectedErr: nil,
		},
		{
			name: "Пользователь не найден",
			input: UpdateRoleInput{
				ActorID: actorID,
				UserID:  userID,
				Role:    domain.RoleManager,
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
			name: "Ошибка при очистке назначений менеджера",
			input: UpdateRoleInput{
				ActorID: actorID,
				UserID:  userID,
				Role:    domain.RoleEmployee,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				user := &domain.User{
					ID:    userID,
					Email: "manager@mail.ru",
					Role:  domain.RoleManager,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(user, nil)

				userRepo.EXPECT().
					ClearManagerAssignments(gomock.Any(), userID).
					Return(assert.AnError)
			},

			expectedErr: assert.AnError,
		},
		{
			name: "Ошибка при обновлении роли",
			input: UpdateRoleInput{
				ActorID: actorID,
				UserID:  userID,
				Role:    domain.RoleManager,
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
					UpdateRole(gomock.Any(), userID, domain.RoleManager).
					Return(assert.AnError)
			},

			expectedErr: assert.AnError,
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

			err := src.UpdateRole(
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
