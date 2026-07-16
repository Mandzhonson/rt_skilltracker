package handler

import (
	"bytes"
	"core_service/internal/domain"
	"core_service/internal/pkg/jwt"
	"core_service/internal/transport/http/handler/mocks"
	"core_service/internal/transport/http/middleware"
	"core_service/internal/usecase/user"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_CreateUser(t *testing.T) {
	type mockBehavior func(s *mocks.MockUserService)

	userID := uuid.New()

	testTable := []struct {
		name string

		body string

		mockBehavior mockBehavior

		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Успешная регистрация",
			body: `{"email":"test@mail.ru","password":"password123","first_name":"Test","last_name":"User"}`,

			mockBehavior: func(s *mocks.MockUserService) {
				s.EXPECT().
					CreateUser(gomock.Any(), user.CreateUserInput{
						Email:     "test@mail.ru",
						Password:  "password123",
						FirstName: "Test",
						LastName:  "User",
					}).
					Return(userID, nil)
			},

			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":"` + userID.String() + `"}`,
		},
		{
			name: "Пользователь уже существует",
			body: `{"email":"test@mail.ru","password":"password123","first_name":"Test","last_name":"User"}`,

			mockBehavior: func(s *mocks.MockUserService) {
				s.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(uuid.Nil, user.ErrUserAlreadyExists)
			},

			expectedStatus: http.StatusConflict,
			expectedBody:   `{"error":"user with this email already exists"}`,
		},
		{
			name: "Неверный email",
			body: `{"email":"invalid","password":"password123","first_name":"Test","last_name":"User"}`,

			mockBehavior: func(s *mocks.MockUserService) {
				s.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Return(uuid.Nil, user.ErrInvalidEmail)
			},

			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid email format"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockUserService(ctrl)

			testCase.mockBehavior(service)

			handler := NewUserHandler(service)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.POST("/auth/register", handler.CreateUser)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString(testCase.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			assert.JSONEq(t, testCase.expectedBody, w.Body.String())
		})
	}
}

func TestUserHandler_GetProfile(t *testing.T) {
	type mockBehavior func(s *mocks.MockUserService)

	userID := uuid.New()
	now := time.Now()

	testTable := []struct {
		name string

		mockBehavior mockBehavior

		expectedStatus int
		expectedBody   func() string
	}{
		{
			name: "Успешное получение профиля",
			mockBehavior: func(s *mocks.MockUserService) {
				user := &domain.User{
					ID:        userID,
					Email:     "test@mail.ru",
					FirstName: "Test",
					LastName:  "User",
					Role:      domain.RoleEmployee,
					Position:  "Developer",
					CreatedAt: now,
					UpdatedAt: now,
				}

				s.EXPECT().
					GetProfile(gomock.Any(), userID).
					Return(user, nil)
			},

			expectedStatus: http.StatusOK,
			expectedBody: func() string {
				return `{"id":"` + userID.String() + `","email":"test@mail.ru","first_name":"Test","last_name":"User","role":"employee","position":"Developer"}`
			},
		},
		{
			name: "Пользователь не найден",
			mockBehavior: func(s *mocks.MockUserService) {
				s.EXPECT().
					GetProfile(gomock.Any(), userID).
					Return(nil, user.ErrUserNotFound)
			},

			expectedStatus: http.StatusNotFound,
			expectedBody: func() string {
				return `{"error":"user not found"}`
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockUserService(ctrl)

			testCase.mockBehavior(service)

			handler := NewUserHandler(service)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.Use(func(c *gin.Context) {
				claims := &jwt.Claims{
					RegisteredClaims: jwtv5.RegisteredClaims{
						Subject: userID.String(),
					},
					Role: domain.RoleEmployee,
				}
				c.Set("claims", claims)
				c.Set("userID", userID)
				c.Next()
			})
			r.GET("/users/me", handler.GetProfile)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/users/me", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			assert.JSONEq(t, testCase.expectedBody(), w.Body.String())
		})
	}
}

func TestUserHandler_UpdateProfile(t *testing.T) {
	type mockBehavior func(s *mocks.MockUserService)

	userID := uuid.New()

	testTable := []struct {
		name string

		body string

		mockBehavior mockBehavior

		expectedStatus int
		expectedBody   func() string
	}{
		{
			name: "Успешное обновление профиля",
			body: `{"email":"updated@mail.ru","first_name":"Updated","last_name":"User"}`,

			mockBehavior: func(s *mocks.MockUserService) {
				updatedUser := &domain.User{
					ID:        userID,
					Email:     "updated@mail.ru",
					FirstName: "Updated",
					LastName:  "User",
					Role:      domain.RoleEmployee,
					Position:  "",
				}

				s.EXPECT().
					UpdateProfile(gomock.Any(), user.UpdateProfileInput{
						UserID:    userID,
						Email:     strPtr("updated@mail.ru"),
						FirstName: strPtr("Updated"),
						LastName:  strPtr("User"),
					}).
					Return(updatedUser, nil)
			},

			expectedStatus: http.StatusOK,
			expectedBody: func() string {
				return `{"id":"` + userID.String() + `","email":"updated@mail.ru","first_name":"Updated","last_name":"User","position":"","role":"employee"}`
			},
		},
		{
			name: "Нет данных для обновления",
			body: `{}`,

			mockBehavior: func(s *mocks.MockUserService) {
				s.EXPECT().
					UpdateProfile(gomock.Any(), user.UpdateProfileInput{
						UserID:    userID,
						Email:     nil,
						FirstName: nil,
						LastName:  nil,
					}).
					Return(nil, user.ErrNoContent)
			},

			expectedStatus: http.StatusNoContent,
			expectedBody: func() string {
				return ""
			},
		},
		{
			name: "Частичное обновление (только email)",
			body: `{"email":"new@mail.ru"}`,

			mockBehavior: func(s *mocks.MockUserService) {
				updatedUser := &domain.User{
					ID:        userID,
					Email:     "new@mail.ru",
					FirstName: "Old",
					LastName:  "User",
					Role:      domain.RoleEmployee,
					Position:  "",
				}

				s.EXPECT().
					UpdateProfile(gomock.Any(), user.UpdateProfileInput{
						UserID:    userID,
						Email:     strPtr("new@mail.ru"),
						FirstName: nil,
						LastName:  nil,
					}).
					Return(updatedUser, nil)
			},

			expectedStatus: http.StatusOK,
			expectedBody: func() string {
				return `{"id":"` + userID.String() + `","email":"new@mail.ru","first_name":"Old","last_name":"User","position":"","role":"employee"}`
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockUserService(ctrl)

			testCase.mockBehavior(service)

			handler := NewUserHandler(service)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.Use(func(c *gin.Context) {
				c.Set("userID", userID)
				c.Next()
			})
			r.PATCH("/users/me", handler.UpdateProfile)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPatch, "/users/me", bytes.NewBufferString(testCase.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			if testCase.expectedStatus == http.StatusNoContent {
				assert.Empty(t, w.Body.String())
			} else {
				assert.JSONEq(t, testCase.expectedBody(), w.Body.String())
			}
		})
	}
}

func TestUserHandler_UpdatePassword(t *testing.T) {
	type mockBehavior func(s *mocks.MockUserService)

	userID := uuid.New()

	testTable := []struct {
		name string

		body string

		mockBehavior mockBehavior

		expectedStatus int
	}{
		{
			name: "Успешная смена пароля",
			body: `{"old_password":"old123","new_password":"new123"}`,

			mockBehavior: func(s *mocks.MockUserService) {
				s.EXPECT().
					UpdatePassword(gomock.Any(), userID, "old123", "new123").
					Return(nil)
			},

			expectedStatus: http.StatusOK,
		},
		{
			name: "Неверный старый пароль",
			body: `{"old_password":"wrong","new_password":"new123"}`,

			mockBehavior: func(s *mocks.MockUserService) {
				s.EXPECT().
					UpdatePassword(gomock.Any(), userID, "wrong", "new123").
					Return(user.ErrInvalidCredentials)
			},

			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockUserService(ctrl)

			testCase.mockBehavior(service)

			handler := NewUserHandler(service)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.Use(func(c *gin.Context) {
				c.Set("userID", userID)
				c.Next()
			})
			r.PATCH("/users/me/password", handler.UpdatePassword)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPatch, "/users/me/password", bytes.NewBufferString(testCase.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
		})
	}
}

func TestUserHandler_GetEmployeesByManager(t *testing.T) {
	type mockBehavior func(s *mocks.MockUserService)

	managerID := uuid.New()
	employee1ID := uuid.New()
	employee2ID := uuid.New()

	testTable := []struct {
		name string

		mockBehavior mockBehavior

		expectedStatus int
		expectedBody   func() string
	}{
		{
			name: "Успешное получение списка сотрудников",
			mockBehavior: func(s *mocks.MockUserService) {
				employees := []*domain.User{
					{
						ID:        employee1ID,
						Email:     "employee1@mail.ru",
						FirstName: "Employee",
						LastName:  "One",
						Role:      domain.RoleEmployee,
						Position:  "Developer",
					},
					{
						ID:        employee2ID,
						Email:     "employee2@mail.ru",
						FirstName: "Employee",
						LastName:  "Two",
						Role:      domain.RoleEmployee,
						Position:  "Designer",
					},
				}

				s.EXPECT().
					GetEmployeesByManager(gomock.Any(), managerID).
					Return(employees, nil)
			},

			expectedStatus: http.StatusOK,
			expectedBody: func() string {
				return `[
					{"id":"` + employee1ID.String() + `","email":"employee1@mail.ru","first_name":"Employee","last_name":"One","position":"Developer","role":"employee"},
					{"id":"` + employee2ID.String() + `","email":"employee2@mail.ru","first_name":"Employee","last_name":"Two","position":"Designer","role":"employee"}
				]`
			},
		},
		{
			name: "Пользователь не менеджер",
			mockBehavior: func(s *mocks.MockUserService) {
				s.EXPECT().
					GetEmployeesByManager(gomock.Any(), managerID).
					Return(nil, user.ErrNotManager)
			},

			expectedStatus: http.StatusForbidden,
			expectedBody: func() string {
				return `{"error":"user is not a manager"}`
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockUserService(ctrl)

			testCase.mockBehavior(service)

			handler := NewUserHandler(service)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.Use(func(c *gin.Context) {
				claims := &jwt.Claims{
					RegisteredClaims: jwtv5.RegisteredClaims{
						Subject: managerID.String(),
					},
					Role: domain.RoleManager,
				}
				c.Set("claims", claims)
				c.Set("userID", managerID)
				c.Next()
			})
			r.Use(middleware.ManagerMiddleware())
			r.GET("/manager/employees", handler.GetEmployeesByManager)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/manager/employees", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			if testCase.expectedStatus == http.StatusOK {
				assert.JSONEq(t, testCase.expectedBody(), w.Body.String())
			} else {
				assert.JSONEq(t, testCase.expectedBody(), w.Body.String())
			}
		})
	}
}

func TestUserHandler_GetEmployeeProfile(t *testing.T) {
	type mockBehavior func(s *mocks.MockUserService)

	managerID := uuid.New()
	employeeID := uuid.New()
	skill1ID := uuid.New()
	skill2ID := uuid.New()
	plan1ID := uuid.New()
	now := time.Now()

	testTable := []struct {
		name string

		employeeID string

		mockBehavior mockBehavior

		expectedStatus int
	}{
		{
			name:       "Успешное получение профиля сотрудника",
			employeeID: employeeID.String(),

			mockBehavior: func(s *mocks.MockUserService) {
				profile := &user.EmployeeProfile{
					User: &domain.User{
						ID:        employeeID,
						Email:     "employee@mail.ru",
						FirstName: "Test",
						LastName:  "Employee",
						Role:      domain.RoleEmployee,
						Position:  "Developer",
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
							Title:            "Plan 1",
							GenerationStatus: domain.GenerationReady,
							CreationType:     domain.CreationAI,
							Progress:         50,
							Status:           domain.PlanActive,
							CreatedAt:        now,
							UpdatedAt:        now,
						},
					},
				}

				s.EXPECT().
					GetEmployeeProfile(gomock.Any(), managerID, employeeID).
					Return(profile, nil)
			},

			expectedStatus: http.StatusOK,
		},
		{
			name:       "Сотрудник не найден",
			employeeID: employeeID.String(),

			mockBehavior: func(s *mocks.MockUserService) {
				s.EXPECT().
					GetEmployeeProfile(gomock.Any(), managerID, employeeID).
					Return(nil, user.ErrUserNotFound)
			},

			expectedStatus: http.StatusNotFound,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockUserService(ctrl)

			testCase.mockBehavior(service)

			handler := NewUserHandler(service)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.Use(func(c *gin.Context) {
				claims := &jwt.Claims{
					RegisteredClaims: jwtv5.RegisteredClaims{
						Subject: managerID.String(),
					},
					Role: domain.RoleManager,
				}
				c.Set("claims", claims)
				c.Set("userID", managerID)
				c.Next()
			})
			r.Use(middleware.ManagerMiddleware())
			r.GET("/manager/employees/:employee_id", handler.GetEmployeeProfile)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/manager/employees/"+testCase.employeeID, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
		})
	}
}
