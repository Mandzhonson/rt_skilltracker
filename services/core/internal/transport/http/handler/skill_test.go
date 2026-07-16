package handler

import (
	"core_service/internal/domain"
	"core_service/internal/pkg/jwt"
	"core_service/internal/transport/http/handler/mocks"
	"core_service/internal/transport/http/middleware"
	"core_service/internal/usecase/skill"
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

func TestSkillHandler_EmployeeList(t *testing.T) {
	type mockBehavior func(s *mocks.MockSkillService)

	employeeID := uuid.New()
	skill1ID := uuid.New()
	skill2ID := uuid.New()
	now := time.Now()

	testTable := []struct {
		name string

		mockBehavior mockBehavior

		expectedStatus int
		expectedBody   func() string
	}{
		{
			name: "Успешное получение навыков сотрудника",
			mockBehavior: func(s *mocks.MockSkillService) {
				skills := []*domain.Skill{
					{
						ID:          skill1ID,
						Name:        "Go",
						Category:    "Programming Language",
						Description: strPtr("Proficient in Go"),
						CreatedAt:   now,
					},
					{
						ID:          skill2ID,
						Name:        "Docker",
						Category:    "DevOps",
						Description: strPtr("Containerization"),
						CreatedAt:   now,
					},
				}

				s.EXPECT().
					ListByUserID(gomock.Any(), employeeID, employeeID).
					Return(skills, nil)
			},

			expectedStatus: http.StatusOK,
			expectedBody: func() string {
				return `{"skills":[
					{"id":"` + skill1ID.String() + `","name":"Go","category":"Programming Language","description":"Proficient in Go","created_at":"` + now.Format(time.RFC3339Nano) + `"},
					{"id":"` + skill2ID.String() + `","name":"Docker","category":"DevOps","description":"Containerization","created_at":"` + now.Format(time.RFC3339Nano) + `"}
				]}`
			},
		},
		{
			name: "Пустой список навыков",
			mockBehavior: func(s *mocks.MockSkillService) {
				s.EXPECT().
					ListByUserID(gomock.Any(), employeeID, employeeID).
					Return([]*domain.Skill{}, nil)
			},

			expectedStatus: http.StatusOK,
			expectedBody: func() string {
				return `{"skills":[]}`
			},
		},
		{
			name: "Ошибка при получении навыков",
			mockBehavior: func(s *mocks.MockSkillService) {
				s.EXPECT().
					ListByUserID(gomock.Any(), employeeID, employeeID).
					Return(nil, errors.New("database error"))
			},

			expectedStatus: http.StatusInternalServerError,
			expectedBody: func() string {
				return `{"error":"internal server error"}`
			},
		},
		{
			name: "Доступ запрещен",
			mockBehavior: func(s *mocks.MockSkillService) {
				s.EXPECT().
					ListByUserID(gomock.Any(), employeeID, employeeID).
					Return(nil, skill.ErrForbidden)
			},

			expectedStatus: http.StatusForbidden,
			expectedBody: func() string {
				return `{"error":"forbidden"}`
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			service := mocks.NewMockSkillService(ctrl)

			testCase.mockBehavior(service)

			handler := NewSkillHandler(service)

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
			r.GET("/employee/skills", handler.EmployeeList)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/employee/skills", nil)

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
