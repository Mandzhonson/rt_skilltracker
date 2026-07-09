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
	GetEmployeePlan(ctx context.Context, employeeID, planID uuid.UUID) (*domain.PlanWithTasks, error)
	ListEmployeePlans(ctx context.Context, employeeID uuid.UUID) ([]*domain.PlanWithTasks, error)
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
	id, err := uuid.Parse(c.Param("plan_id"))
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

func (h *PlanHandler) EmployeeGetPlans(c *gin.Context) {
	employeeID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	plans, err := h.service.ListEmployeePlans(c.Request.Context(), employeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	response := make([]dto.PlanWithTasksResponse, 0, len(plans))
	for _, p := range plans {
		tasks := make([]dto.TaskResponse, 0, len(p.Tasks))
		for _, t := range p.Tasks {
			tasks = append(tasks, dto.TaskResponse{
				ID:          t.ID.String(),
				PlanID:      t.PlanID.String(),
				Title:       t.Title,
				Description: t.Description,
				Position:    t.Position,
				Status:      string(t.Status),
				CreatedAt:   t.CreatedAt,
				UpdatedAt:   t.UpdatedAt,
			})
		}
		response = append(response, dto.PlanWithTasksResponse{
			Plan: dto.PlanResponse{
				ID:           p.Plan.ID.String(),
				EmployeeID:   p.Plan.EmployeeID.String(),
				CreatedBy:    p.Plan.CreatedBy.String(),
				Title:        p.Plan.Title,
				Description:  p.Plan.Description,
				CreationType: string(p.Plan.CreationType),
				Progress:     p.Plan.Progress,
				Status:       string(p.Plan.Status),
				CreatedAt:    p.Plan.CreatedAt,
				UpdatedAt:    p.Plan.UpdatedAt,
			},
			Tasks: tasks,
		})
	}
	c.JSON(http.StatusOK, response)
}

func (h *PlanHandler) EmployeeGetPlan(c *gin.Context) {
	employeeID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	planID, err := uuid.Parse(c.Param("plan_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan id"})
		return
	}

	planEntity, err := h.service.GetEmployeePlan(c.Request.Context(), employeeID, planID)

	if err != nil {
		switch {
		case errors.Is(err, plan.ErrPlanNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, plan.ErrInvalidPlanID):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, plan.ErrEmployeeForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}
	tasks := make([]dto.TaskResponse, 0, len(planEntity.Tasks))
	for _, t := range planEntity.Tasks {
		tasks = append(tasks, dto.TaskResponse{
			ID:          t.ID.String(),
			PlanID:      t.PlanID.String(),
			Title:       t.Title,
			Description: t.Description,
			Position:    t.Position,
			Status:      string(t.Status),
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, dto.PlanWithTasksResponse{
		Plan: dto.PlanResponse{
			ID:           planEntity.Plan.ID.String(),
			EmployeeID:   planEntity.Plan.EmployeeID.String(),
			CreatedBy:    planEntity.Plan.CreatedBy.String(),
			Title:        planEntity.Plan.Title,
			Description:  planEntity.Plan.Description,
			CreationType: string(planEntity.Plan.CreationType),
			Progress:     planEntity.Plan.Progress,
			Status:       string(planEntity.Plan.Status),
			CreatedAt:    planEntity.Plan.CreatedAt,
			UpdatedAt:    planEntity.Plan.UpdatedAt,
		},
		Tasks: tasks,
	})
}
