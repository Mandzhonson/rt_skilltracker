package handler

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/transport/http/dto"
	"core_service/internal/transport/http/middleware"
	"core_service/internal/usecase/plan"
	"core_service/internal/usecase/user"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PlanService interface {
	Create(ctx context.Context, input plan.CreatePlanInput) (uuid.UUID, error)
	GetByID(ctx context.Context, managerID uuid.UUID, id uuid.UUID) (*domain.Plan, error)
}

type PlanHandler struct {
	service PlanService
}

func NewPlanHandler(planService PlanService) *PlanHandler {
	return &PlanHandler{
		service: planService,
	}
}

func (h *PlanHandler) Create(c *gin.Context) {
	managerID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.CreatePlanRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	employeeID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}

	id, err := h.service.Create(
		c.Request.Context(),
		plan.CreatePlanInput{
			EmployeeID:  employeeID,
			CreatedBy:   managerID,
			Title:       req.Title,
			Description: req.Description,
		},
	)

	if err != nil {
		switch {
		case errors.Is(err, plan.ErrInvalidTitle),
			errors.Is(err, plan.ErrInvalidEmployee),
			errors.Is(err, plan.ErrInvalidCreator):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		case errors.Is(err, plan.ErrEmployeeNotAssigned):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, user.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	c.JSON(http.StatusCreated, dto.CreatePlanResponse{ID: id.String()})
}

func (h *PlanHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan id"})
		return
	}
	managerID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	entity, err := h.service.GetByID(c.Request.Context(), managerID, id)
	if err != nil {
		switch {
		case errors.Is(err, plan.ErrPlanNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, plan.ErrInvalidPlanID):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}
	c.JSON(http.StatusOK, dto.PlanResponse{
		ID:           entity.ID.String(),
		EmployeeID:   entity.EmployeeID.String(),
		CreatedBy:    entity.CreatedBy.String(),
		Title:        entity.Title,
		Description:  entity.Description,
		CreationType: string(entity.CreationType),
		Progress:     entity.Progress,
		Status:       string(entity.Status),
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    entity.UpdatedAt,
	})
}
