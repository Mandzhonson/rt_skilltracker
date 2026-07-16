package handler

import (
	"bytes"
	"core_service/internal/domain"
	"core_service/internal/pkg/jwt"
	"core_service/internal/transport/http/handler/mocks"
	"core_service/internal/transport/http/middleware"
	"core_service/internal/usecase/plan"
	"core_service/internal/usecase/user"
	"errors"
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

func strPtr(s string) *string {
	return &s
}

func TestPlanHandler_Create(t *testing.T) {
	type mockBehavior func(s *mocks.MockPlanService)

	managerID := uuid.New()
	employeeID := uuid.New()
	planID := uuid.New()

	testTable := []struct {
		name string

		body string

		mockBehavior mockBehavior

		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Успешное создание плана",
			body: `{"employee_id":"` + employeeID.String() + `","title":"Test Plan","description":"Test Description"}`,

			mockBehavior: func(s *mocks.MockPlanService) {
				s.EXPECT().
					Create(gomock.Any(), plan.CreatePlanInput{
						EmployeeID:  employeeID,
						CreatedBy:   managerID,
						Title:       "Test Plan",
						Description: strPtr("Test Description"),
					}).
					Return(planID, nil)
			},

			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":"` + planID.String() + `"}`,
		},
		{
			name: "Неверный запрос - invalid employee id",
			body: `{"employee_id":"invalid","title":"Test"}`,

			mockBehavior: func(s *mocks.MockPlanService) {
			},

			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid employee id"}`,
		},
		{
			name: "Сотрудник не закреплен за менеджером",
			body: `{"employee_id":"` + employeeID.String() + `","title":"Test Plan"}`,

			mockBehavior: func(s *mocks.MockPlanService) {
				s.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(uuid.Nil, plan.ErrEmployeeNotAssigned)
			},

			expectedStatus: http.StatusForbidden,
			expectedBody:   `{"error":"employee is not assigned to manager"}`,
		},
		{
			name: "Сотрудник не найден",
			body: `{"employee_id":"` + employeeID.String() + `","title":"Test Plan"}`,

			mockBehavior: func(s *mocks.MockPlanService) {
				s.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(uuid.Nil, user.ErrUserNotFound)
			},

			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"user not found"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockPlanService(ctrl)

			testCase.mockBehavior(service)

			handler := NewPlanHandler(service)

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
			r.POST("/manager/plans", handler.Create)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/manager/plans", bytes.NewBufferString(testCase.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			assert.JSONEq(t, testCase.expectedBody, w.Body.String())
		})
	}
}

func TestPlanHandler_GetByID(t *testing.T) {
	type mockBehavior func(s *mocks.MockPlanService)

	managerID := uuid.New()
	planID := uuid.New()
	employeeID := uuid.New()
	now := time.Now()

	testTable := []struct {
		name string

		planID string

		mockBehavior mockBehavior

		expectedStatus int
		expectedBody   func() string
	}{
		{
			name:   "Успешное получение плана",
			planID: planID.String(),

			mockBehavior: func(s *mocks.MockPlanService) {
				planWithTasks := &domain.PlanWithTasks{
					Plan: &domain.Plan{
						ID:               planID,
						EmployeeID:       employeeID,
						CreatedBy:        managerID,
						Title:            "Test Plan",
						Description:      strPtr("Test Description"),
						GenerationStatus: domain.GenerationReady,
						CreationType:     domain.CreationAI,
						Progress:         50,
						Status:           domain.PlanActive,
						CreatedAt:        now,
						UpdatedAt:        now,
					},
					Tasks: []*domain.Task{},
				}

				s.EXPECT().
					GetByIDWithTasks(gomock.Any(), managerID, planID).
					Return(planWithTasks, nil)
			},

			expectedStatus: http.StatusOK,
			expectedBody: func() string {
				return `{"plan":{"id":"` + planID.String() + `","employee_id":"` + employeeID.String() + `","created_by":"` + managerID.String() + `","title":"Test Plan","description":"Test Description","generation_status":"ready","creation_type":"ai","progress":50,"status":"active","created_at":"` + now.Format(time.RFC3339Nano) + `","updated_at":"` + now.Format(time.RFC3339Nano) + `"},"tasks":[]}`
			},
		},
		{
			name:   "План не найден",
			planID: planID.String(),

			mockBehavior: func(s *mocks.MockPlanService) {
				s.EXPECT().
					GetByIDWithTasks(gomock.Any(), managerID, planID).
					Return(nil, plan.ErrPlanNotFound)
			},

			expectedStatus: http.StatusNotFound,
			expectedBody: func() string {
				return `{"error":"plan not found"}`
			},
		},
		{
			name:   "Нет прав (другой менеджер)",
			planID: planID.String(),

			mockBehavior: func(s *mocks.MockPlanService) {
				s.EXPECT().
					GetByIDWithTasks(gomock.Any(), managerID, planID).
					Return(nil, plan.ErrManagerForbidden)
			},

			expectedStatus: http.StatusForbidden,
			expectedBody: func() string {
				return `{"error":"manager has no access"}`
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockPlanService(ctrl)

			testCase.mockBehavior(service)

			handler := NewPlanHandler(service)

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
			r.GET("/manager/plans/:plan_id", handler.GetByID)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/manager/plans/"+testCase.planID, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			assert.JSONEq(t, testCase.expectedBody(), w.Body.String())
		})
	}
}

func TestPlanHandler_ListByManager(t *testing.T) {
	type mockBehavior func(s *mocks.MockPlanService)

	managerID := uuid.New()
	plan1ID := uuid.New()
	plan2ID := uuid.New()

	testTable := []struct {
		name string

		mockBehavior mockBehavior

		expectedStatus int
	}{
		{
			name: "Успешное получение списка планов",
			mockBehavior: func(s *mocks.MockPlanService) {
				plans := []*domain.Plan{
					{
						ID:       plan1ID,
						Title:    "Plan 1",
						Status:   domain.PlanActive,
						Progress: 50,
					},
					{
						ID:       plan2ID,
						Title:    "Plan 2",
						Status:   domain.PlanActive,
						Progress: 80,
					},
				}

				s.EXPECT().
					ListByManager(gomock.Any(), managerID).
					Return(plans, nil)
			},

			expectedStatus: http.StatusOK,
		},
		{
			name: "Ошибка при получении списка",
			mockBehavior: func(s *mocks.MockPlanService) {
				s.EXPECT().
					ListByManager(gomock.Any(), managerID).
					Return(nil, errors.New("database error"))
			},

			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockPlanService(ctrl)

			testCase.mockBehavior(service)

			handler := NewPlanHandler(service)

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
			r.GET("/manager/plans", handler.ListByManager)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/manager/plans", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
		})
	}
}

func TestPlanHandler_Update(t *testing.T) {
	type mockBehavior func(s *mocks.MockPlanService)

	managerID := uuid.New()
	planID := uuid.New()

	testTable := []struct {
		name string

		planID string
		body   string

		mockBehavior mockBehavior

		expectedStatus int
	}{
		{
			name:   "Успешное обновление плана",
			planID: planID.String(),
			body:   `{"title":"Updated Plan","description":"Updated Description"}`,

			mockBehavior: func(s *mocks.MockPlanService) {
				s.EXPECT().
					Update(gomock.Any(), plan.UpdatePlanInput{
						PlanID:      planID,
						ManagerID:   managerID,
						Title:       "Updated Plan",
						Description: strPtr("Updated Description"),
					}).
					Return(nil)
			},

			expectedStatus: http.StatusOK,
		},
		{
			name:   "План не найден",
			planID: planID.String(),
			body:   `{"title":"Updated Plan"}`,

			mockBehavior: func(s *mocks.MockPlanService) {
				s.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(plan.ErrPlanNotFound)
			},

			expectedStatus: http.StatusNotFound,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockPlanService(ctrl)

			testCase.mockBehavior(service)

			handler := NewPlanHandler(service)

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
			r.PATCH("/manager/plans/:plan_id", handler.Update)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPatch, "/manager/plans/"+testCase.planID, bytes.NewBufferString(testCase.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
		})
	}
}

func TestPlanHandler_Delete(t *testing.T) {
	type mockBehavior func(s *mocks.MockPlanService)

	managerID := uuid.New()
	planID := uuid.New()

	testTable := []struct {
		name string

		planID string

		mockBehavior mockBehavior

		expectedStatus int
	}{
		{
			name:   "Успешное удаление плана",
			planID: planID.String(),

			mockBehavior: func(s *mocks.MockPlanService) {
				s.EXPECT().
					Delete(gomock.Any(), managerID, planID).
					Return(nil)
			},

			expectedStatus: http.StatusNoContent,
		},
		{
			name:   "План не найден",
			planID: planID.String(),

			mockBehavior: func(s *mocks.MockPlanService) {
				s.EXPECT().
					Delete(gomock.Any(), managerID, planID).
					Return(plan.ErrPlanNotFound)
			},

			expectedStatus: http.StatusNotFound,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockPlanService(ctrl)

			testCase.mockBehavior(service)

			handler := NewPlanHandler(service)

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
			r.DELETE("/manager/plans/:plan_id", handler.Delete)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/manager/plans/"+testCase.planID, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
		})
	}
}

func TestPlanHandler_Archive(t *testing.T) {
	type mockBehavior func(s *mocks.MockPlanService)

	managerID := uuid.New()
	planID := uuid.New()

	testTable := []struct {
		name string

		planID string

		mockBehavior mockBehavior

		expectedStatus int
	}{
		{
			name:   "Успешное архивирование плана",
			planID: planID.String(),

			mockBehavior: func(s *mocks.MockPlanService) {
				s.EXPECT().
					Archive(gomock.Any(), managerID, planID).
					Return(nil)
			},

			expectedStatus: http.StatusOK,
		},
		{
			name:   "План не найден",
			planID: planID.String(),

			mockBehavior: func(s *mocks.MockPlanService) {
				s.EXPECT().
					Archive(gomock.Any(), managerID, planID).
					Return(plan.ErrPlanNotFound)
			},

			expectedStatus: http.StatusNotFound,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockPlanService(ctrl)

			testCase.mockBehavior(service)

			handler := NewPlanHandler(service)

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
			r.PATCH("/manager/plans/:plan_id/archive", handler.Archive)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPatch, "/manager/plans/"+testCase.planID+"/archive", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
		})
	}
}

func TestPlanHandler_CreateAI(t *testing.T) {
	type mockBehavior func(s *mocks.MockPlanService)

	managerID := uuid.New()
	employeeID := uuid.New()
	planID := uuid.New()

	testTable := []struct {
		name string

		body string

		mockBehavior mockBehavior

		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Успешное создание AI плана",
			body: `{"employee_id":"` + employeeID.String() + `","topic":"Go Programming","description":"Learn Go"}`,

			mockBehavior: func(s *mocks.MockPlanService) {
				s.EXPECT().
					CreateAI(gomock.Any(), plan.CreateAIInput{
						EmployeeID:  employeeID,
						CreatedBy:   managerID,
						Topic:       "Go Programming",
						Description: "Learn Go",
					}).
					Return(planID, nil)
			},

			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":"` + planID.String() + `"}`,
		},
		{
			name: "Сотрудник не закреплен за менеджером",
			body: `{"employee_id":"` + employeeID.String() + `","topic":"Go Programming"}`,

			mockBehavior: func(s *mocks.MockPlanService) {
				s.EXPECT().
					CreateAI(gomock.Any(), gomock.Any()).
					Return(uuid.Nil, plan.ErrEmployeeNotAssigned)
			},

			expectedStatus: http.StatusForbidden,
			expectedBody:   `{"error":"employee is not assigned to manager"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockPlanService(ctrl)

			testCase.mockBehavior(service)

			handler := NewPlanHandler(service)

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
			r.POST("/manager/plans/ai", handler.CreateAI)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/manager/plans/ai", bytes.NewBufferString(testCase.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			assert.JSONEq(t, testCase.expectedBody, w.Body.String())
		})
	}
}

func TestPlanHandler_EmployeeGetPlans(t *testing.T) {
	type mockBehavior func(s *mocks.MockPlanService)

	employeeID := uuid.New()
	plan1ID := uuid.New()
	now := time.Now()

	testTable := []struct {
		name string

		mockBehavior mockBehavior

		expectedStatus int
	}{
		{
			name: "Успешное получение планов сотрудника",
			mockBehavior: func(s *mocks.MockPlanService) {
				plans := []*domain.PlanWithTasks{
					{
						Plan: &domain.Plan{
							ID:               plan1ID,
							EmployeeID:       employeeID,
							Title:            "Plan 1",
							GenerationStatus: domain.GenerationReady,
							CreationType:     domain.CreationAI,
							Progress:         50,
							Status:           domain.PlanActive,
							CreatedAt:        now,
							UpdatedAt:        now,
						},
						Tasks: []*domain.Task{},
					},
				}

				s.EXPECT().
					ListEmployeePlans(gomock.Any(), employeeID).
					Return(plans, nil)
			},

			expectedStatus: http.StatusOK,
		},
		{
			name: "Ошибка при получении планов",
			mockBehavior: func(s *mocks.MockPlanService) {
				s.EXPECT().
					ListEmployeePlans(gomock.Any(), employeeID).
					Return(nil, errors.New("database error"))
			},

			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockPlanService(ctrl)

			testCase.mockBehavior(service)

			handler := NewPlanHandler(service)

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
			r.GET("/employee/plans", handler.EmployeeGetPlans)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/employee/plans", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
		})
	}
}

func TestPlanHandler_EmployeeGetPlan(t *testing.T) {
	type mockBehavior func(s *mocks.MockPlanService)

	employeeID := uuid.New()
	planID := uuid.New()
	now := time.Now()

	testTable := []struct {
		name string

		planID string

		mockBehavior mockBehavior

		expectedStatus int
	}{
		{
			name:   "Успешное получение плана сотрудника",
			planID: planID.String(),

			mockBehavior: func(s *mocks.MockPlanService) {
				planWithTasks := &domain.PlanWithTasks{
					Plan: &domain.Plan{
						ID:               planID,
						EmployeeID:       employeeID,
						Title:            "Test Plan",
						GenerationStatus: domain.GenerationReady,
						CreationType:     domain.CreationAI,
						Progress:         50,
						Status:           domain.PlanActive,
						CreatedAt:        now,
						UpdatedAt:        now,
					},
					Tasks: []*domain.Task{},
				}

				s.EXPECT().
					GetEmployeePlan(gomock.Any(), employeeID, planID).
					Return(planWithTasks, nil)
			},

			expectedStatus: http.StatusOK,
		},
		{
			name:   "План не найден",
			planID: planID.String(),

			mockBehavior: func(s *mocks.MockPlanService) {
				s.EXPECT().
					GetEmployeePlan(gomock.Any(), employeeID, planID).
					Return(nil, plan.ErrPlanNotFound)
			},

			expectedStatus: http.StatusNotFound,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockPlanService(ctrl)

			testCase.mockBehavior(service)

			handler := NewPlanHandler(service)

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
			r.GET("/employee/plans/:plan_id", handler.EmployeeGetPlan)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/employee/plans/"+testCase.planID, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
		})
	}
}
