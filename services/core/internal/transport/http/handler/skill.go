package handler

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/transport/http/dto"
	"core_service/internal/transport/http/middleware"
	"core_service/internal/usecase/skill"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//go:generate mockgen -source=skill.go -destination=mocks/mock_skill_handler.go -package=mocks
type SkillService interface {
	ListByUserID(ctx context.Context, requesterID, userID uuid.UUID) ([]*domain.Skill, error)
}

type SkillHandler struct {
	service SkillService
}

func NewSkillHandler(s SkillService) *SkillHandler {
	return &SkillHandler{
		service: s,
	}
}

// EmployeeList godoc
// @Summary Получить список навыков сотрудника
// @Description Возвращает список навыков текущего сотрудника
// @Tags Employee
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.ListSkillsResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /employee/skills [get]
func (h *SkillHandler) EmployeeList(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	skills, err := h.service.ListByUserID(c.Request.Context(), userID, userID)
	if err != nil {
		switch {
		case errors.Is(err, skill.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	response := make([]dto.SkillResponse, 0, len(skills))

	for _, s := range skills {
		response = append(response, dto.SkillResponse{
			ID:          s.ID.String(),
			Name:        s.Name,
			Category:    s.Category,
			Description: s.Description,
			CreatedAt:   s.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, dto.ListSkillsResponse{Skills: response})
}
