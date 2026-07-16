package handler

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/transport/http/dto"
	"core_service/internal/transport/http/middleware"
	"core_service/internal/usecase/test"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TestService interface {
	GetForEmployee(ctx context.Context, userID uuid.UUID, planID uuid.UUID) (*domain.TestView, error)
	Submit(ctx context.Context, input test.SubmitTestInput) (*domain.TestResult, error)
	GetForManager(ctx context.Context, managerID uuid.UUID, planID uuid.UUID) (*domain.TestView, error)
	GetForAdmin(ctx context.Context, planID uuid.UUID) (*domain.TestView, error)
}

type TestHandler struct {
	service TestService
}

func NewTestHandler(service TestService) *TestHandler {
	return &TestHandler{
		service: service,
	}
}

func (h *TestHandler) GetForEmployee(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	planID, err := uuid.Parse(c.Param("plan_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan id"})
		return
	}

	result, err := h.service.GetForEmployee(c.Request.Context(), userID, planID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.StartTestResponse{TestID: result.ID.String()}
	response.Questions = make([]dto.QuestionResponse, 0, len(result.Questions))
	for _, q := range result.Questions {

		response.Questions = append(response.Questions, dto.QuestionResponse{
			ID:       q.ID.String(),
			Question: q.Text,
			Options:  q.Options,
		},
		)
	}
	c.JSON(http.StatusOK, response)
}

func (h *TestHandler) Submit(c *gin.Context) {

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	planID, err := uuid.Parse(c.Param("plan_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan id"})
		return
	}

	var req dto.SubmitTestRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	answers := make([]test.AnswerInput, 0, len(req.Answers))
	for _, a := range req.Answers {
		questionID, err := uuid.Parse(a.QuestionID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid question id"})
			return
		}
		answers = append(answers, test.AnswerInput{QuestionID: questionID, Answer: a.Answer})
	}

	result, err := h.service.Submit(c.Request.Context(), test.SubmitTestInput{
		UserID:  userID,
		PlanID:  planID,
		Answers: answers,
	},
	)
	if err != nil {
		switch {
		case errors.Is(err, test.ErrTestNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, dto.SubmitTestResponse{
		Score:  result.Score,
		Total:  result.Total,
		Passed: result.Passed,
	},
	)
}

func (h *TestHandler) ManagerGetTest(c *gin.Context) {
	managerID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	planID, err := uuid.Parse(c.Param("plan_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan id"})
		return
	}

	testEntity, err := h.service.GetForManager(c.Request.Context(), managerID, planID)
	if err != nil {

		switch {
		case errors.Is(err, test.ErrTestNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	response := dto.StartTestResponse{
		TestID:    testEntity.ID.String(),
		Questions: make([]dto.QuestionResponse, 0, len(testEntity.Questions)),
	}

	for _, q := range testEntity.Questions {

		response.Questions = append(
			response.Questions,
			dto.QuestionResponse{
				ID:       q.ID.String(),
				Question: q.Text,
				Options:  q.Options,
			},
		)
	}

	c.JSON(http.StatusOK, response)
}

func (h *TestHandler) AdminGetTest(c *gin.Context) {

	planID, err := uuid.Parse(c.Param("plan_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid plan id",
		})
		return
	}

	testEntity, err := h.service.GetForAdmin(
		c.Request.Context(),
		planID,
	)

	if err != nil {

		switch {
		case errors.Is(err, test.ErrTestNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}

		return
	}

	response := dto.StartTestResponse{
		TestID:    testEntity.ID.String(),
		Questions: make([]dto.QuestionResponse, 0, len(testEntity.Questions)),
	}

	for _, q := range testEntity.Questions {

		response.Questions = append(
			response.Questions,
			dto.QuestionResponse{
				ID:       q.ID.String(),
				Question: q.Text,
				Options:  q.Options,
			},
		)
	}

	c.JSON(http.StatusOK, response)
}
