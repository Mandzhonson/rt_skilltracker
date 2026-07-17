package admin

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/repository/model"
	mock_postgres "core_service/internal/repository/postgres/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAdminService_ListUsers(t *testing.T) {
	type mockBehavior func(
		userRepo *mock_postgres.MockUserRepository,
	)

	userID1 := uuid.New()
	userID2 := uuid.New()
	roleAdmin := domain.RoleAdmin
	roleManager := domain.RoleManager
	searchTerm := "test"

	testTable := []struct {
		name string

		input ListUsersInput

		mockBehavior mockBehavior

		expectedUsers []*domain.User
		expectedErr   error
	}{
		{
			name: "Успешное получение списка пользователей (с пагинацией)",
			input: ListUsersInput{
				Page:   1,
				Limit:  20,
				Role:   nil,
				Search: nil,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				users := []*domain.User{
					{
						ID:        userID1,
						Email:     "user1@mail.ru",
						FirstName: "User",
						LastName:  "One",
						Role:      domain.RoleEmployee,
					},
					{
						ID:        userID2,
						Email:     "user2@mail.ru",
						FirstName: "User",
						LastName:  "Two",
						Role:      domain.RoleManager,
					},
				}

				expectedParams := model.ListUsersParams{
					Offset: 0,
					Limit:  20,
					Role:   nil,
					Search: nil,
				}

				userRepo.EXPECT().
					ListUsers(gomock.Any(), expectedParams).
					Return(users, nil)
			},

			expectedUsers: []*domain.User{
				{
					ID:        userID1,
					Email:     "user1@mail.ru",
					FirstName: "User",
					LastName:  "One",
					Role:      domain.RoleEmployee,
				},
				{
					ID:        userID2,
					Email:     "user2@mail.ru",
					FirstName: "User",
					LastName:  "Two",
					Role:      domain.RoleManager,
				},
			},
			expectedErr: nil,
		},
		{
			name: "Успешное получение списка пользователей (с фильтром по роли)",
			input: ListUsersInput{
				Page:   1,
				Limit:  20,
				Role:   &roleAdmin,
				Search: nil,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				users := []*domain.User{
					{
						ID:        userID1,
						Email:     "admin@mail.ru",
						FirstName: "Admin",
						LastName:  "User",
						Role:      domain.RoleAdmin,
					},
				}

				expectedParams := model.ListUsersParams{
					Offset: 0,
					Limit:  20,
					Role:   &roleAdmin,
					Search: nil,
				}

				userRepo.EXPECT().
					ListUsers(gomock.Any(), expectedParams).
					Return(users, nil)
			},

			expectedUsers: []*domain.User{
				{
					ID:        userID1,
					Email:     "admin@mail.ru",
					FirstName: "Admin",
					LastName:  "User",
					Role:      domain.RoleAdmin,
				},
			},
			expectedErr: nil,
		},
		{
			name: "Успешное получение списка пользователей (с поиском)",
			input: ListUsersInput{
				Page:   1,
				Limit:  20,
				Role:   nil,
				Search: &searchTerm,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				users := []*domain.User{
					{
						ID:        userID1,
						Email:     "test@mail.ru",
						FirstName: "Test",
						LastName:  "User",
						Role:      domain.RoleEmployee,
					},
				}

				expectedParams := model.ListUsersParams{
					Offset: 0,
					Limit:  20,
					Role:   nil,
					Search: &searchTerm,
				}

				userRepo.EXPECT().
					ListUsers(gomock.Any(), expectedParams).
					Return(users, nil)
			},

			expectedUsers: []*domain.User{
				{
					ID:        userID1,
					Email:     "test@mail.ru",
					FirstName: "Test",
					LastName:  "User",
					Role:      domain.RoleEmployee,
				},
			},
			expectedErr: nil,
		},
		{
			name: "Успешное получение списка пользователей (с фильтром по роли и поиском)",
			input: ListUsersInput{
				Page:   1,
				Limit:  20,
				Role:   &roleManager,
				Search: &searchTerm,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				users := []*domain.User{
					{
						ID:        userID1,
						Email:     "test@mail.ru",
						FirstName: "Test",
						LastName:  "Manager",
						Role:      domain.RoleManager,
					},
				}

				expectedParams := model.ListUsersParams{
					Offset: 0,
					Limit:  20,
					Role:   &roleManager,
					Search: &searchTerm,
				}

				userRepo.EXPECT().
					ListUsers(gomock.Any(), expectedParams).
					Return(users, nil)
			},

			expectedUsers: []*domain.User{
				{
					ID:        userID1,
					Email:     "test@mail.ru",
					FirstName: "Test",
					LastName:  "Manager",
					Role:      domain.RoleManager,
				},
			},
			expectedErr: nil,
		},
		{
			name: "Пустой список пользователей",
			input: ListUsersInput{
				Page:   1,
				Limit:  20,
				Role:   nil,
				Search: nil,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				expectedParams := model.ListUsersParams{
					Offset: 0,
					Limit:  20,
					Role:   nil,
					Search: nil,
				}

				userRepo.EXPECT().
					ListUsers(gomock.Any(), expectedParams).
					Return([]*domain.User{}, nil)
			},

			expectedUsers: []*domain.User{},
			expectedErr:   nil,
		},
		{
			name: "Ошибка при получении списка пользователей",
			input: ListUsersInput{
				Page:   1,
				Limit:  20,
				Role:   nil,
				Search: nil,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				expectedParams := model.ListUsersParams{
					Offset: 0,
					Limit:  20,
					Role:   nil,
					Search: nil,
				}

				userRepo.EXPECT().
					ListUsers(gomock.Any(), expectedParams).
					Return(nil, assert.AnError)
			},

			expectedUsers: nil,
			expectedErr:   assert.AnError,
		},
		{
			name: "Page меньше 1 - устанавливается в 1",
			input: ListUsersInput{
				Page:   0,
				Limit:  20,
				Role:   nil,
				Search: nil,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				expectedParams := model.ListUsersParams{
					Offset: 0,
					Limit:  20,
					Role:   nil,
					Search: nil,
				}

				userRepo.EXPECT().
					ListUsers(gomock.Any(), expectedParams).
					Return([]*domain.User{}, nil)
			},

			expectedUsers: []*domain.User{},
			expectedErr:   nil,
		},
		{
			name: "Limit меньше 1 - устанавливается в 20",
			input: ListUsersInput{
				Page:   1,
				Limit:  0,
				Role:   nil,
				Search: nil,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				expectedParams := model.ListUsersParams{
					Offset: 0,
					Limit:  20,
					Role:   nil,
					Search: nil,
				}

				userRepo.EXPECT().
					ListUsers(gomock.Any(), expectedParams).
					Return([]*domain.User{}, nil)
			},

			expectedUsers: []*domain.User{},
			expectedErr:   nil,
		},
		{
			name: "Пагинация на второй странице",
			input: ListUsersInput{
				Page:   2,
				Limit:  10,
				Role:   nil,
				Search: nil,
			},

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
			) {
				users := []*domain.User{
					{
						ID:        userID1,
						Email:     "user1@mail.ru",
						FirstName: "User",
						LastName:  "One",
						Role:      domain.RoleEmployee,
					},
				}

				expectedParams := model.ListUsersParams{
					Offset: 10,
					Limit:  10,
					Role:   nil,
					Search: nil,
				}

				userRepo.EXPECT().
					ListUsers(gomock.Any(), expectedParams).
					Return(users, nil)
			},

			expectedUsers: []*domain.User{
				{
					ID:        userID1,
					Email:     "user1@mail.ru",
					FirstName: "User",
					LastName:  "One",
					Role:      domain.RoleEmployee,
				},
			},
			expectedErr: nil,
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

			users, err := src.ListUsers(
				context.Background(),
				testCase.input,
			)

			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedErr, err)
				assert.Nil(t, users)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(testCase.expectedUsers), len(users))
				if len(users) > 0 {
					assert.Equal(t, testCase.expectedUsers[0].ID, users[0].ID)
					assert.Equal(t, testCase.expectedUsers[0].Email, users[0].Email)
				}
			}
		})
	}
}