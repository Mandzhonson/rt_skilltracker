package handler

import (
	"bytes"
	"core_service/internal/domain"
	"core_service/internal/pkg/jwt"
	"core_service/internal/transport/http/handler/mocks"
	"core_service/internal/transport/http/middleware"
	"core_service/internal/usecase/task"
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

func TestTaskHandler_Create(t *testing.T) {
	type mockBehavior func(s *mocks.MockTaskService)

	managerID := uuid.New()
	planID := uuid.New()
	taskID := uuid.New()

	testTable := []struct {
		name string

		planID string
		body   string

		mockBehavior mockBehavior

		expectedStatus int
	}{
		{
			name:   "Успешное создание задачи",
			planID: planID.String(),
			body:   `{"title":"Test Task","description":"Test Description"}`,

			mockBehavior: func(s *mocks.MockTaskService) {
				s.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(taskID, nil).
					AnyTimes()
			},

			expectedStatus: http.StatusCreated,
		},
		{
			name:   "Пустое название задачи",
			planID: planID.String(),
			body:   `{"title":"","description":"Test"}`,

			mockBehavior: func(s *mocks.MockTaskService) {
				s.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(uuid.Nil, task.ErrInvalidTitle).
					AnyTimes()
			},

			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "План не найден",
			planID: planID.String(),
			body:   `{"title":"Test Task"}`,

			mockBehavior: func(s *mocks.MockTaskService) {
				s.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(uuid.Nil, task.ErrPlanNotFound).
					AnyTimes()
			},

			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "План архивирован",
			planID: planID.String(),
			body:   `{"title":"Test Task"}`,

			mockBehavior: func(s *mocks.MockTaskService) {
				s.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(uuid.Nil, task.ErrPlanArchived).
					AnyTimes()
			},

			expectedStatus: http.StatusConflict,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockTaskService(ctrl)

			testCase.mockBehavior(service)

			handler := NewTaskHandler(service)

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
			r.POST("/manager/plans/:plan_id/tasks", handler.Create)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/manager/plans/"+testCase.planID+"/tasks", bytes.NewBufferString(testCase.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
		})
	}
}

func TestTaskHandler_GetByID(t *testing.T) {
	type mockBehavior func(s *mocks.MockTaskService)

	managerID := uuid.New()
	taskID := uuid.New()
	planID := uuid.New()
	now := time.Now()

	testTable := []struct {
		name string

		taskID string

		mockBehavior mockBehavior

		expectedStatus int
		expectedBody   func() string
	}{
		{
			name:   "Успешное получение задачи",
			taskID: taskID.String(),

			mockBehavior: func(s *mocks.MockTaskService) {
				task := &domain.Task{
					ID:          taskID,
					PlanID:      planID,
					Title:       "Test Task",
					Description: strPtr("Test Description"),
					Position:    1,
					Status:      domain.TaskTodo,
					CreatedAt:   now,
					UpdatedAt:   now,
				}

				s.EXPECT().
					GetByID(gomock.Any(), managerID, taskID).
					Return(task, nil)
			},

			expectedStatus: http.StatusOK,
			expectedBody: func() string {
				return `{"id":"` + taskID.String() + `","plan_id":"` + planID.String() + `","title":"Test Task","description":"Test Description","position":1,"status":"todo","created_at":"` + now.Format(time.RFC3339Nano) + `","updated_at":"` + now.Format(time.RFC3339Nano) + `"}`
			},
		},
		{
			name:   "Задача не найдена",
			taskID: taskID.String(),

			mockBehavior: func(s *mocks.MockTaskService) {
				s.EXPECT().
					GetByID(gomock.Any(), managerID, taskID).
					Return(nil, task.ErrTaskNotFound)
			},

			expectedStatus: http.StatusNotFound,
			expectedBody: func() string {
				return `{"error":"task not found"}`
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockTaskService(ctrl)

			testCase.mockBehavior(service)

			handler := NewTaskHandler(service)

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
			r.GET("/manager/tasks/:task_id", handler.GetByID)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/manager/tasks/"+testCase.taskID, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			assert.JSONEq(t, testCase.expectedBody(), w.Body.String())
		})
	}
}

func TestTaskHandler_Update(t *testing.T) {
	type mockBehavior func(s *mocks.MockTaskService)

	managerID := uuid.New()
	taskID := uuid.New()

	testTable := []struct {
		name string

		taskID string
		body   string

		mockBehavior mockBehavior

		expectedStatus int
	}{
		{
			name:   "Успешное обновление задачи",
			taskID: taskID.String(),
			body:   `{"title":"Updated Task","description":"Updated Description"}`,

			mockBehavior: func(s *mocks.MockTaskService) {
				s.EXPECT().
					Update(gomock.Any(), task.UpdateTaskInput{
						TaskID:      taskID,
						ManagerID:   managerID,
						Title:       strPtr("Updated Task"),
						Description: strPtr("Updated Description"),
					}).
					Return(nil)
			},

			expectedStatus: http.StatusOK,
		},
		{
			name:   "Задача не найдена",
			taskID: taskID.String(),
			body:   `{"title":"Updated Task"}`,

			mockBehavior: func(s *mocks.MockTaskService) {
				s.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(task.ErrTaskNotFound)
			},

			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "Нет данных для обновления",
			taskID: taskID.String(),
			body:   `{}`,

			mockBehavior: func(s *mocks.MockTaskService) {
				s.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(task.ErrInvalidUpdate)
			},

			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockTaskService(ctrl)

			testCase.mockBehavior(service)

			handler := NewTaskHandler(service)

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
			r.PATCH("/manager/tasks/:task_id", handler.Update)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPatch, "/manager/tasks/"+testCase.taskID, bytes.NewBufferString(testCase.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
		})
	}
}

func TestTaskHandler_Delete(t *testing.T) {
	type mockBehavior func(s *mocks.MockTaskService)

	managerID := uuid.New()
	taskID := uuid.New()

	testTable := []struct {
		name string

		taskID string

		mockBehavior mockBehavior

		expectedStatus int
	}{
		{
			name:   "Успешное удаление задачи",
			taskID: taskID.String(),

			mockBehavior: func(s *mocks.MockTaskService) {
				s.EXPECT().
					Delete(gomock.Any(), managerID, taskID).
					Return(nil)
			},

			expectedStatus: http.StatusNoContent,
		},
		{
			name:   "Задача не найдена",
			taskID: taskID.String(),

			mockBehavior: func(s *mocks.MockTaskService) {
				s.EXPECT().
					Delete(gomock.Any(), managerID, taskID).
					Return(task.ErrTaskNotFound)
			},

			expectedStatus: http.StatusNotFound,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockTaskService(ctrl)

			testCase.mockBehavior(service)

			handler := NewTaskHandler(service)

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
			r.DELETE("/manager/tasks/:task_id", handler.Delete)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/manager/tasks/"+testCase.taskID, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
		})
	}
}

func TestTaskHandler_UpdateStatus(t *testing.T) {
	type mockBehavior func(s *mocks.MockTaskService)

	employeeID := uuid.New()
	taskID := uuid.New()

	testTable := []struct {
		name string

		taskID string
		body   string

		mockBehavior mockBehavior

		expectedStatus int
	}{
		{
			name:   "Успешное обновление статуса задачи",
			taskID: taskID.String(),
			body:   `{"status":"in_progress"}`,

			mockBehavior: func(s *mocks.MockTaskService) {
				s.EXPECT().
					UpdateStatus(gomock.Any(), task.UpdateTaskStatusInput{
						TaskID: taskID,
						UserID: employeeID,
						Status: domain.TaskInProgress,
					}).
					Return(nil)
			},

			expectedStatus: http.StatusNoContent,
		},
		{
			name:   "Задача не найдена",
			taskID: taskID.String(),
			body:   `{"status":"done"}`,

			mockBehavior: func(s *mocks.MockTaskService) {
				s.EXPECT().
					UpdateStatus(gomock.Any(), gomock.Any()).
					Return(task.ErrTaskNotFound)
			},

			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "Неверный статус",
			taskID: taskID.String(),
			body:   `{"status":"invalid"}`,

			mockBehavior: func(s *mocks.MockTaskService) {
				s.EXPECT().
					UpdateStatus(gomock.Any(), gomock.Any()).
					Return(task.ErrInvalidStatus)
			},

			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockTaskService(ctrl)

			testCase.mockBehavior(service)

			handler := NewTaskHandler(service)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.Use(func(c *gin.Context) {
				claims := &jwt.Claims{
					RegisteredClaims: jwtv5.RegisteredClaims{
						Subject: employeeID.String(),
					},
					Role: domain.RoleEmployee,
				}
				c.Set("claims", claims)
				c.Set("userID", employeeID)
				c.Next()
			})
			r.Use(middleware.EmployeeMiddleware())
			r.PATCH("/employee/tasks/:task_id/status", handler.UpdateStatus)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPatch, "/employee/tasks/"+testCase.taskID+"/status", bytes.NewBufferString(testCase.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
		})
	}
}
