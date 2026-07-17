package handler

import (
	"bytes"
	"core_service/internal/transport/http/handler/mocks"
	"core_service/internal/usecase/ai"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAIHandler_GeneratePlan(t *testing.T) {
	type mockBehavior func(s *mocks.MockAIService) *ai.GeneratedPlan

	testTable := []struct {
		name string

		body string

		mockBehavior mockBehavior

		expectedStatus int
		expectedBody   func(plan *ai.GeneratedPlan) string
	}{
		{
			name: "Успешная генерация плана",
			body: `{"topic":"Go Programming","description":"Learn Go from scratch","position":"Junior Developer"}`,

			mockBehavior: func(s *mocks.MockAIService) *ai.GeneratedPlan {
				plan := &ai.GeneratedPlan{
					Title:       "Go Programming Plan",
					Description: "Learn Go from scratch",
					Tasks: []ai.GeneratedTask{
						{
							Title:       "Learn Go basics",
							Description: "Study Go syntax and fundamentals",
						},
						{
							Title:       "Build a REST API",
							Description: "Create a simple REST API in Go",
						},
					},
				}

				s.EXPECT().
					GeneratePlan(gomock.Any(), ai.GeneratePlanInput{
						Topic:       "Go Programming",
						Description: "Learn Go from scratch",
						Position:    "Junior Developer",
					}).
					Return(plan, nil)

				return plan
			},

			expectedStatus: http.StatusOK,
			expectedBody: func(plan *ai.GeneratedPlan) string {
				body, _ := json.Marshal(plan)
				return string(body)
			},
		},
		{
			name: "Неверный запрос (невалидный JSON)",
			body: `{"topic":"Go Programming","description":"Learn Go from scratch"`,

			mockBehavior: func(s *mocks.MockAIService) *ai.GeneratedPlan {
				return nil
			},

			expectedStatus: http.StatusBadRequest,
			expectedBody: func(plan *ai.GeneratedPlan) string {
				return `{"error":"invalid request body"}`
			},
		},
		{
			name: "Ошибка при генерации",
			body: `{"topic":"Go Programming","description":"Learn Go"}`,

			mockBehavior: func(s *mocks.MockAIService) *ai.GeneratedPlan {
				s.EXPECT().
					GeneratePlan(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("AI service error"))

				return nil
			},

			expectedStatus: http.StatusInternalServerError,
			expectedBody: func(plan *ai.GeneratedPlan) string {
				return `{"error":"AI service error"}`
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockAIService(ctrl)

			plan := testCase.mockBehavior(service)

			handler := NewAIHandler(service)

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.POST("/ai/generate-plan", handler.GeneratePlan)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/ai/generate-plan", bytes.NewBufferString(testCase.body))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatus, w.Code)
			if testCase.expectedStatus == http.StatusOK {
				assert.JSONEq(t, testCase.expectedBody(plan), w.Body.String())
			} else {
				assert.JSONEq(t, testCase.expectedBody(plan), w.Body.String())
			}
		})
	}
}
