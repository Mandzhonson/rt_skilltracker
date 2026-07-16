package handler

import (
	"bytes"
	"core_service/internal/domain"
	"core_service/internal/transport/http/handler/mocks"
	"core_service/internal/usecase/admin"
	"core_service/internal/usecase/user"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAdminHandler_ListUsers(t *testing.T) {
	type mockBehavior func(s *mocks.MockAdminService) []*domain.User

	testTable := []struct {
		name string

		url string

		mockBehavior mockBehavior

		expectedStatus int
		expectedBody   func(users []*domain.User) string
	}{
		{
			name: "Успешное получение списка пользователей",
			url:  "/admin/users?page=1&limit=10",

			mockBehavior: func(s *mocks.MockAdminService) []*domain.User {
				user1ID := uuid.New()
				user2ID := uuid.New()

				users := []*domain.User{
					{
						ID:        user1ID,
						Email:     "user1@mail.ru",
						FirstName: "User",
						LastName:  "One",
						Role:      domain.RoleEmployee,
					},
					{
						ID:        user2ID,
						Email:     "user2@mail.ru",
						FirstName: "User",
						LastName:  "Two",
						Role:      domain.RoleManager,
					},
				}

				s.EXPECT().
					ListUsers(gomock.Any(), admin.ListUsersInput{
						Page:   1,
						Limit:  10,
						Role:   nil,
						Search: nil,
					}).
					Return(users, nil)

				return users
			},

			expectedStatus: http.StatusOK,
			expectedBody: func(users []*domain.User) string {
				return `{"users":[
					{"id":"` + users[0].ID.String() + `","email":"user1@mail.ru","first_name":"User","last_name":"One","role":"employee","position":""},
					{"id":"` + users[1].ID.String() + `","email":"user2@mail.ru","first_name":"User","last_name":"Two","role":"manager","position":""}
				]}`
			},
		},
		{
			name: "Ошибка при получении списка пользователей",
			url:  "/admin/users",

			mockBehavior: func(s *mocks.MockAdminService) []*domain.User {
				s.EXPECT().
					ListUsers(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("database error"))

				return nil
			},

			expectedStatus: http.StatusInternalServerError,
			expectedBody: func(users []*domain.User) string {
				return `{"error":"internal server error"}`
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockAdminService(ctrl)

			users := testCase.mockBehavior(service)

			handler := NewAdminHandler(service)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.GET("/admin/users", handler.ListUsers)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, testCase.url, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			if testCase.expectedStatus == http.StatusOK {
				assert.JSONEq(t, testCase.expectedBody(users), w.Body.String())
			}
		})
	}
}

func TestAdminHandler_GetUser(t *testing.T) {
	type mockBehavior func(s *mocks.MockAdminService) *domain.User

	userID := uuid.New()
	managerID := uuid.New()
	now := time.Now()

	testTable := []struct {
		name string

		userID string

		mockBehavior mockBehavior

		expectedStatus int
		expectedBody   func(user *domain.User) string
	}{
		{
			name:   "Успешное получение пользователя",
			userID: userID.String(),

			mockBehavior: func(s *mocks.MockAdminService) *domain.User {
				user := &domain.User{
					ID:        userID,
					Email:     "test@mail.ru",
					FirstName: "Test",
					LastName:  "User",
					Role:      domain.RoleEmployee,
					Position:  "Developer",
					ManagerID: &managerID,
					CreatedAt: now,
					UpdatedAt: now,
				}

				s.EXPECT().
					GetUser(gomock.Any(), userID).
					Return(user, nil)

				return user
			},

			expectedStatus: http.StatusOK,
			expectedBody: func(user *domain.User) string {
				managerIDStr := ""
				if user.ManagerID != nil {
					managerIDStr = user.ManagerID.String()
				}
				return `{"id":"` + user.ID.String() + `","email":"test@mail.ru","first_name":"Test","last_name":"User","role":"employee","position":"Developer","manager_id":"` + managerIDStr + `","created_at":"` + user.CreatedAt.Format(time.RFC3339Nano) + `","updated_at":"` + user.UpdatedAt.Format(time.RFC3339Nano) + `"}`
			},
		},
		{
			name:   "Пользователь не найден",
			userID: userID.String(),

			mockBehavior: func(s *mocks.MockAdminService) *domain.User {
				s.EXPECT().
					GetUser(gomock.Any(), userID).
					Return(nil, user.ErrUserNotFound)

				return nil
			},

			expectedStatus: http.StatusNotFound,
			expectedBody: func(user *domain.User) string {
				return `{"error":"user not found"}`
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockAdminService(ctrl)

			user := testCase.mockBehavior(service)

			handler := NewAdminHandler(service)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.GET("/admin/users/:id", handler.GetUser)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/admin/users/"+testCase.userID, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			assert.JSONEq(t, testCase.expectedBody(user), w.Body.String())
		})
	}
}

func TestAdminHandler_UpdateRole(t *testing.T) {
	type mockBehavior func(s *mocks.MockAdminService)

	userID := uuid.New()
	actorID := uuid.New()

	testTable := []struct {
		name string

		userID string
		body   string

		mockBehavior mockBehavior

		expectedStatus int
	}{
		{
			name:   "Успешное обновление роли",
			userID: userID.String(),
			body:   `{"role":"manager"}`,

			mockBehavior: func(s *mocks.MockAdminService) {
				s.EXPECT().
					UpdateRole(gomock.Any(), admin.UpdateRoleInput{
						ActorID: actorID,
						UserID:  userID,
						Role:    domain.RoleManager,
					}).
					Return(nil)
			},

			expectedStatus: http.StatusOK,
		},
		{
			name:   "Неверный запрос",
			userID: userID.String(),
			body:   `{"role":"invalid"}`,

			mockBehavior: func(s *mocks.MockAdminService) {
				s.EXPECT().
					UpdateRole(gomock.Any(), gomock.Any()).
					Return(admin.ErrInvalidRole)
			},

			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockAdminService(ctrl)

			testCase.mockBehavior(service)

			handler := NewAdminHandler(service)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.Use(func(c *gin.Context) {
				c.Set("user_id", actorID)
				c.Next()
			})
			r.PATCH("/admin/users/:id/role", handler.UpdateRole)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPatch, "/admin/users/"+testCase.userID+"/role", bytes.NewBufferString(testCase.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
		})
	}
}

func TestAdminHandler_AssignManager(t *testing.T) {
	type mockBehavior func(s *mocks.MockAdminService)

	userID := uuid.New()
	managerID := uuid.New()

	testTable := []struct {
		name string

		userID string
		body   string

		mockBehavior mockBehavior

		expectedStatus int
	}{
		{
			name:   "Успешное назначение менеджера",
			userID: userID.String(),
			body:   `{"manager_id":"` + managerID.String() + `"}`,

			mockBehavior: func(s *mocks.MockAdminService) {
				s.EXPECT().
					AssignManager(gomock.Any(), admin.AssignManagerInput{
						UserID:    userID,
						ManagerID: managerID,
					}).
					Return(nil)
			},

			expectedStatus: http.StatusOK,
		},
		{
			name:   "Назначение самого себя",
			userID: userID.String(),
			body:   `{"manager_id":"` + userID.String() + `"}`,

			mockBehavior: func(s *mocks.MockAdminService) {
				s.EXPECT().
					AssignManager(gomock.Any(), gomock.Any()).
					Return(admin.ErrAssignYourself)
			},

			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockAdminService(ctrl)

			testCase.mockBehavior(service)

			handler := NewAdminHandler(service)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.PATCH("/admin/users/:id/manager", handler.AssignManager)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPatch, "/admin/users/"+testCase.userID+"/manager", bytes.NewBufferString(testCase.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
		})
	}
}

func TestAdminHandler_RemoveManager(t *testing.T) {
	type mockBehavior func(s *mocks.MockAdminService)

	userID := uuid.New()

	testTable := []struct {
		name string

		userID string

		mockBehavior mockBehavior

		expectedStatus int
	}{
		{
			name:   "Успешное удаление менеджера",
			userID: userID.String(),

			mockBehavior: func(s *mocks.MockAdminService) {
				s.EXPECT().
					RemoveManager(gomock.Any(), userID).
					Return(nil)
			},

			expectedStatus: http.StatusOK,
		},
		{
			name:   "Менеджер не назначен",
			userID: userID.String(),

			mockBehavior: func(s *mocks.MockAdminService) {
				s.EXPECT().
					RemoveManager(gomock.Any(), userID).
					Return(admin.ErrManagerNotAssigned)
			},

			expectedStatus: http.StatusNoContent,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockAdminService(ctrl)

			testCase.mockBehavior(service)

			handler := NewAdminHandler(service)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.DELETE("/admin/users/:id/manager", handler.RemoveManager)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/admin/users/"+testCase.userID+"/manager", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
		})
	}
}

func TestAdminHandler_UpdatePosition(t *testing.T) {
	type mockBehavior func(s *mocks.MockAdminService)

	userID := uuid.New()

	testTable := []struct {
		name string

		userID string
		body   string

		mockBehavior mockBehavior

		expectedStatus int
	}{
		{
			name:   "Успешное обновление должности",
			userID: userID.String(),
			body:   `{"position":"Senior Developer"}`,

			mockBehavior: func(s *mocks.MockAdminService) {
				s.EXPECT().
					UpdatePosition(gomock.Any(), userID, "Senior Developer").
					Return(nil)
			},

			expectedStatus: http.StatusOK,
		},
		{
			name:   "Неверный запрос",
			userID: userID.String(),
			body:   `{"position":""}`,

			mockBehavior: func(s *mocks.MockAdminService) {
				s.EXPECT().
					UpdatePosition(gomock.Any(), userID, "").
					Return(admin.ErrInvalidPosition)
			},

			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockAdminService(ctrl)

			testCase.mockBehavior(service)

			handler := NewAdminHandler(service)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.PATCH("/admin/users/:id/position", handler.UpdatePosition)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPatch, "/admin/users/"+testCase.userID+"/position", bytes.NewBufferString(testCase.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
		})
	}
}
