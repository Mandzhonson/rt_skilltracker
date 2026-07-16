package handler

import (
	"context"
	"core_service/internal/transport/http/dto"
	"core_service/internal/usecase/ai"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:generate mockgen -source=ai.go -destination=mocks/mock_ai_handler.go -package=mocks
type AIService interface {
	GeneratePlan(ctx context.Context, input ai.GeneratePlanInput) (*ai.GeneratedPlan, error)
	ExtractSkills(ctx context.Context, input ai.ExtractSkillsInput) ([]ai.SkillCandidate, error)
}

type AIHandler struct {
	service AIService
}

func NewAIHandler(service AIService) *AIHandler {
	return &AIHandler{
		service: service,
	}
}

// GeneratePlan godoc
// @Summary Сгенерировать план с помощью ИИ
// @Description Генерирует план обучения на основе темы и описания
// @Tags AI
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.GeneratePlanRequest true "Данные для генерации плана"
// @Success 200 {object} ai.GeneratedPlan
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /ai/generate-plan [post]
func (h *AIHandler) GeneratePlan(c *gin.Context) {
	var req dto.GeneratePlanRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if strings.TrimSpace(req.Topic) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "topic is required"})
		return
	}
	plan, err := h.service.GeneratePlan(c.Request.Context(), ai.GeneratePlanInput{
		Topic:       req.Topic,
		Description: req.Description,
		Position:    req.Position,
	},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, plan)
}
