package skill

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

func TestSkillService_ListByUserID(t *testing.T) {
	type mockBehavior func(
		userRepo *mock_postgres.MockUserRepository,
		skillRepo *mock_postgres.MockSkillRepository,
	)

	requesterID := uuid.New()
	userID := uuid.New()
	managerID := uuid.New()
	skill1ID := uuid.New()
	skill2ID := uuid.New()

	testTable := []struct {
		name string

		requesterID uuid.UUID
		userID      uuid.UUID

		mockBehavior mockBehavior

		expectedSkills []*domain.Skill
		expectedErr    error
	}{
		{
			name:        "Сотрудник запрашивает свои навыки",
			requesterID: userID,
			userID:      userID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				skillRepo *mock_postgres.MockSkillRepository,
			) {
				requester := &domain.User{
					ID:    userID,
					Email: "employee@mail.ru",
					Role:  domain.RoleEmployee,
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

				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(requester, nil)

				skillRepo.EXPECT().
					ListByUserID(gomock.Any(), userID).
					Return(skills, nil)
			},

			expectedSkills: []*domain.Skill{
				{
					ID:   skill1ID,
					Name: "Go",
				},
				{
					ID:   skill2ID,
					Name: "Docker",
				},
			},
			expectedErr: nil,
		},
		{
			name:        "Сотрудник пытается получить навыки другого сотрудника - запрещено",
			requesterID: requesterID,
			userID:      userID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				skillRepo *mock_postgres.MockSkillRepository,
			) {
				requester := &domain.User{
					ID:    requesterID,
					Email: "employee1@mail.ru",
					Role:  domain.RoleEmployee,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), requesterID).
					Return(requester, nil)
			},

			expectedSkills: nil,
			expectedErr:    ErrForbidden,
		},
		{
			name:        "Менеджер запрашивает навыки своего сотрудника",
			requesterID: managerID,
			userID:      userID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				skillRepo *mock_postgres.MockSkillRepository,
			) {
				requester := &domain.User{
					ID:    managerID,
					Email: "manager@mail.ru",
					Role:  domain.RoleManager,
				}

				employee := &domain.User{
					ID:        userID,
					Email:     "employee@mail.ru",
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
					Return(requester, nil)

				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(employee, nil)

				skillRepo.EXPECT().
					ListByUserID(gomock.Any(), userID).
					Return(skills, nil)
			},

			expectedSkills: []*domain.Skill{
				{
					ID:   skill1ID,
					Name: "Go",
				},
			},
			expectedErr: nil,
		},
		{
			name:        "Менеджер запрашивает навыки сотрудника, который ему не подчиняется - запрещено",
			requesterID: managerID,
			userID:      userID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				skillRepo *mock_postgres.MockSkillRepository,
			) {
				requester := &domain.User{
					ID:    managerID,
					Email: "manager@mail.ru",
					Role:  domain.RoleManager,
				}

				otherManagerID := uuid.New()
				employee := &domain.User{
					ID:        userID,
					Email:     "employee@mail.ru",
					Role:      domain.RoleEmployee,
					ManagerID: &otherManagerID,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), managerID).
					Return(requester, nil)

				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(employee, nil)
			},

			expectedSkills: nil,
			expectedErr:    ErrForbidden,
		},
		{
			name:        "Менеджер запрашивает навыки сотрудника без менеджера - запрещено",
			requesterID: managerID,
			userID:      userID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				skillRepo *mock_postgres.MockSkillRepository,
			) {
				requester := &domain.User{
					ID:    managerID,
					Email: "manager@mail.ru",
					Role:  domain.RoleManager,
				}

				employee := &domain.User{
					ID:        userID,
					Email:     "employee@mail.ru",
					Role:      domain.RoleEmployee,
					ManagerID: nil,
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), managerID).
					Return(requester, nil)

				userRepo.EXPECT().
					GetById(gomock.Any(), userID).
					Return(employee, nil)
			},

			expectedSkills: nil,
			expectedErr:    ErrForbidden,
		},
		{
			name:        "Админ запрашивает навыки сотрудника",
			requesterID: requesterID,
			userID:      userID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				skillRepo *mock_postgres.MockSkillRepository,
			) {
				requester := &domain.User{
					ID:    requesterID,
					Email: "admin@mail.ru",
					Role:  domain.RoleAdmin,
				}

				skills := []*domain.Skill{
					{
						ID:   skill1ID,
						Name: "Go",
					},
				}

				userRepo.EXPECT().
					GetById(gomock.Any(), requesterID).
					Return(requester, nil)

				skillRepo.EXPECT().
					ListByUserID(gomock.Any(), userID).
					Return(skills, nil)
			},

			expectedSkills: []*domain.Skill{
				{
					ID:   skill1ID,
					Name: "Go",
				},
			},
			expectedErr: nil,
		},
		{
			name:        "Пользователь не найден",
			requesterID: requesterID,
			userID:      userID,

			mockBehavior: func(
				userRepo *mock_postgres.MockUserRepository,
				skillRepo *mock_postgres.MockSkillRepository,
			) {
				userRepo.EXPECT().
					GetById(gomock.Any(), requesterID).
					Return(nil, postgres.ErrUserNotFound)
			},

			expectedSkills: nil,
			expectedErr:    postgres.ErrUserNotFound,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepo := mock_postgres.NewMockUserRepository(ctrl)
			skillRepo := mock_postgres.NewMockSkillRepository(ctrl)

			testCase.mockBehavior(userRepo, skillRepo)

			src := NewSkillService(
				skillRepo,
				userRepo,
			)

			skills, err := src.ListByUserID(
				context.Background(),
				testCase.requesterID,
				testCase.userID,
			)

			if testCase.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, testCase.expectedErr, err)
				assert.Nil(t, skills)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(testCase.expectedSkills), len(skills))
				if len(skills) > 0 {
					assert.Equal(t, testCase.expectedSkills[0].ID, skills[0].ID)
					assert.Equal(t, testCase.expectedSkills[0].Name, skills[0].Name)
				}
			}
		})
	}
}
