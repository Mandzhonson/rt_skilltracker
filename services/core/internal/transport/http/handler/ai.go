package handler

import (
	"context"
	"core_service/internal/transport/http/dto"
	"core_service/internal/usecase/ai"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

func (h *AIHandler) GeneratePlan(c *gin.Context) {
	var req dto.GeneratePlanRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
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
