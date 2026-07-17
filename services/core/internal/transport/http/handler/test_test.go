package handler

import (
	"bytes"
	"core_service/internal/domain"
	"core_service/internal/pkg/jwt"
	"core_service/internal/transport/http/handler/mocks"
	"core_service/internal/transport/http/middleware"
	"core_service/internal/usecase/test"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTestHandler_GetForEmployee(t *testing.T) {
	type mockBehavior func(s *mocks.MockTestService)

	employeeID := uuid.New()
	planID := uuid.New()
	testID := uuid.New()
	question1ID := uuid.New()
	question2ID := uuid.New()

	testTable := []struct {
		name string

		planID string

		mockBehavior mockBehavior

		expectedStatus int
		expectedBody   func() string
	}{
		{
			name:   "Успешное получение теста",
			planID: planID.String(),

			mockBehavior: func(s *mocks.MockTestService) {
				testView := &domain.TestView{
					ID: testID,
					Questions: []domain.QuestionView{
						{
							ID:   question1ID,
							Text: "What is Go?",
							Options: []string{
								"Programming language",
								"Framework",
								"Library",
								"Tool",
							},
						},
						{
							ID:   question2ID,
							Text: "What is Docker?",
							Options: []string{
								"Containerization",
								"Programming language",
								"Database",
								"Framework",
							},
						},
					},
				}

				s.EXPECT().
					GetForEmployee(gomock.Any(), employeeID, planID).
					Return(testView, nil)
			},

			expectedStatus: http.StatusOK,
			expectedBody: func() string {
				return `{"test_id":"` + testID.String() + `","questions":[
					{"id":"` + question1ID.String() + `","question":"What is Go?","options":["Programming language","Framework","Library","Tool"]},
					{"id":"` + question2ID.String() + `","question":"What is Docker?","options":["Containerization","Programming language","Database","Framework"]}
				]}`
			},
		},
		{
			name:   "Тест не найден",
			planID: planID.String(),

			mockBehavior: func(s *mocks.MockTestService) {
				s.EXPECT().
					GetForEmployee(gomock.Any(), employeeID, planID).
					Return(nil, test.ErrTestNotFound)
			},

			expectedStatus: http.StatusNotFound,
			expectedBody: func() string {
				return `{"error":"test not found"}`
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockTestService(ctrl)

			testCase.mockBehavior(service)

			handler := NewTestHandler(service)

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
			r.GET("/employee/plans/:plan_id/test", handler.GetForEmployee)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/employee/plans/"+testCase.planID+"/test", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			assert.JSONEq(t, testCase.expectedBody(), w.Body.String())
		})
	}
}

func TestTestHandler_Submit(t *testing.T) {
	type mockBehavior func(s *mocks.MockTestService)

	employeeID := uuid.New()
	planID := uuid.New()
	question1ID := uuid.New()
	question2ID := uuid.New()

	testTable := []struct {
		name string

		planID string
		body   string

		mockBehavior mockBehavior

		expectedStatus int
		expectedBody   func() string
	}{
		{
			name:   "Успешная отправка теста (пройден)",
			planID: planID.String(),
			body: `{"answers":[
				{"question_id":"` + question1ID.String() + `","answer":"A"},
				{"question_id":"` + question2ID.String() + `","answer":"B"}
			]}`,

			mockBehavior: func(s *mocks.MockTestService) {
				result := &domain.TestResult{
					Score:  100,
					Total:  2,
					Passed: true,
				}

				answers := []test.AnswerInput{
					{QuestionID: question1ID, Answer: "A"},
					{QuestionID: question2ID, Answer: "B"},
				}

				s.EXPECT().
					Submit(gomock.Any(), test.SubmitTestInput{
						UserID:  employeeID,
						PlanID:  planID,
						Answers: answers,
					}).
					Return(result, nil)
			},

			expectedStatus: http.StatusOK,
			expectedBody: func() string {
				return `{"score":100,"total":2,"passed":true}`
			},
		},
		{
			name:   "Тест не найден",
			planID: planID.String(),
			body:   `{"answers":[]}`,

			mockBehavior: func(s *mocks.MockTestService) {
				s.EXPECT().
					Submit(gomock.Any(), gomock.Any()).
					Return(nil, test.ErrTestNotFound)
			},

			expectedStatus: http.StatusNotFound,
			expectedBody: func() string {
				return `{"error":"test not found"}`
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockTestService(ctrl)

			testCase.mockBehavior(service)

			handler := NewTestHandler(service)

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
			r.POST("/employee/plans/:plan_id/test", handler.Submit)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/employee/plans/"+testCase.planID+"/test", bytes.NewBufferString(testCase.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			assert.JSONEq(t, testCase.expectedBody(), w.Body.String())
		})
	}
}

func TestTestHandler_ManagerGetTest(t *testing.T) {
	type mockBehavior func(s *mocks.MockTestService)

	managerID := uuid.New()
	planID := uuid.New()
	testID := uuid.New()
	question1ID := uuid.New()
	question2ID := uuid.New()

	testTable := []struct {
		name string

		planID string

		mockBehavior mockBehavior

		expectedStatus int
		expectedBody   func() string
	}{
		{
			name:   "Успешное получение теста менеджером",
			planID: planID.String(),

			mockBehavior: func(s *mocks.MockTestService) {
				testView := &domain.TestView{
					ID: testID,
					Questions: []domain.QuestionView{
						{
							ID:   question1ID,
							Text: "What is Go?",
							Options: []string{
								"Programming language",
								"Framework",
								"Library",
								"Tool",
							},
						},
						{
							ID:   question2ID,
							Text: "What is Docker?",
							Options: []string{
								"Containerization",
								"Programming language",
								"Database",
								"Framework",
							},
						},
					},
				}

				s.EXPECT().
					GetForManager(gomock.Any(), managerID, planID).
					Return(testView, nil)
			},

			expectedStatus: http.StatusOK,
			expectedBody: func() string {
				return `{"test_id":"` + testID.String() + `","questions":[
					{"id":"` + question1ID.String() + `","question":"What is Go?","options":["Programming language","Framework","Library","Tool"]},
					{"id":"` + question2ID.String() + `","question":"What is Docker?","options":["Containerization","Programming language","Database","Framework"]}
				]}`
			},
		},
		{
			name:   "Тест не найден",
			planID: planID.String(),

			mockBehavior: func(s *mocks.MockTestService) {
				s.EXPECT().
					GetForManager(gomock.Any(), managerID, planID).
					Return(nil, test.ErrTestNotFound)
			},

			expectedStatus: http.StatusNotFound,
			expectedBody: func() string {
				return `{"error":"test not found"}`
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockTestService(ctrl)

			testCase.mockBehavior(service)

			handler := NewTestHandler(service)

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
			r.GET("/manager/plans/:plan_id/test", handler.ManagerGetTest)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/manager/plans/"+testCase.planID+"/test", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			assert.JSONEq(t, testCase.expectedBody(), w.Body.String())
		})
	}
}

func TestTestHandler_AdminGetTest(t *testing.T) {
	type mockBehavior func(s *mocks.MockTestService)

	planID := uuid.New()
	testID := uuid.New()
	question1ID := uuid.New()
	question2ID := uuid.New()

	testTable := []struct {
		name string

		planID string

		mockBehavior mockBehavior

		expectedStatus int
		expectedBody   func() string
	}{
		{
			name:   "Успешное получение теста администратором",
			planID: planID.String(),

			mockBehavior: func(s *mocks.MockTestService) {
				testView := &domain.TestView{
					ID: testID,
					Questions: []domain.QuestionView{
						{
							ID:   question1ID,
							Text: "What is Go?",
							Options: []string{
								"Programming language",
								"Framework",
								"Library",
								"Tool",
							},
						},
						{
							ID:   question2ID,
							Text: "What is Docker?",
							Options: []string{
								"Containerization",
								"Programming language",
								"Database",
								"Framework",
							},
						},
					},
				}

				s.EXPECT().
					GetForAdmin(gomock.Any(), planID).
					Return(testView, nil)
			},

			expectedStatus: http.StatusOK,
			expectedBody: func() string {
				return `{"test_id":"` + testID.String() + `","questions":[
					{"id":"` + question1ID.String() + `","question":"What is Go?","options":["Programming language","Framework","Library","Tool"]},
					{"id":"` + question2ID.String() + `","question":"What is Docker?","options":["Containerization","Programming language","Database","Framework"]}
				]}`
			},
		},
		{
			name:   "Тест не найден",
			planID: planID.String(),

			mockBehavior: func(s *mocks.MockTestService) {
				s.EXPECT().
					GetForAdmin(gomock.Any(), planID).
					Return(nil, test.ErrTestNotFound)
			},

			expectedStatus: http.StatusNotFound,
			expectedBody: func() string {
				return `{"error":"test not found"}`
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockTestService(ctrl)

			testCase.mockBehavior(service)

			handler := NewTestHandler(service)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.Use(func(c *gin.Context) {
				claims := &jwt.Claims{
					RegisteredClaims: jwtv5.RegisteredClaims{
						Subject: "admin-id",
					},
					Role: domain.RoleAdmin,
				}
				c.Set("claims", claims)
				c.Set("userID", uuid.New())
				c.Next()
			})
			r.Use(middleware.AdminMiddleware())
			r.GET("/admin/plans/:plan_id/test", handler.AdminGetTest)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/admin/plans/"+testCase.planID+"/test", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			assert.JSONEq(t, testCase.expectedBody(), w.Body.String())
		})
	}
}
