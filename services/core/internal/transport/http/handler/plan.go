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

//go:generate mockgen -source=plan.go -destination=mocks/mock_plan_handler.go -package=mocks
type PlanService interface {
	Create(ctx context.Context, input plan.CreatePlanInput) (uuid.UUID, error)
	CreateAI(ctx context.Context, input plan.CreateAIInput) (uuid.UUID, error)
	GetByID(ctx context.Context, managerID uuid.UUID, id uuid.UUID) (*domain.Plan, error)
	ListByManager(ctx context.Context, managerID uuid.UUID) ([]*domain.Plan, error)
	GetEmployeePlan(ctx context.Context, employeeID, planID uuid.UUID) (*domain.PlanWithTasks, error)
	ListEmployeePlans(ctx context.Context, employeeID uuid.UUID) ([]*domain.PlanWithTasks, error)
	GetByIDWithTasks(ctx context.Context, managerID uuid.UUID, id uuid.UUID) (*domain.PlanWithTasks, error)
	Update(ctx context.Context, input plan.UpdatePlanInput) error
	Delete(ctx context.Context, managerID uuid.UUID, planID uuid.UUID) error
	Archive(ctx context.Context, managerID, planID uuid.UUID) error
}

type PlanHandler struct {
	service PlanService
}

func NewPlanHandler(planService PlanService) *PlanHandler {
	return &PlanHandler{
		service: planService,
	}
}

// Create godoc
// @Summary Создать план
// @Description Создает новый план развития для сотрудника
// @Tags Plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreatePlanRequest true "Данные для создания плана"
// @Success 201 {object} dto.CreatePlanResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /manager/plans [post]
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

// GetByID godoc
// @Summary Получить план по ID
// @Description Возвращает детальную информацию о плане с задачами
// @Tags Plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param plan_id path string true "ID плана"
// @Success 200 {object} dto.PlanWithTasksResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /manager/plans/{plan_id} [get]
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

	entity, err := h.service.GetByIDWithTasks(c.Request.Context(), managerID, id)
	if err != nil {
		switch {
		case errors.Is(err, plan.ErrPlanNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, plan.ErrInvalidPlanID):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, plan.ErrManagerForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	tasks := make([]dto.TaskResponse, 0, len(entity.Tasks))
	for _, t := range entity.Tasks {
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
			ID:               entity.Plan.ID.String(),
			EmployeeID:       entity.Plan.EmployeeID.String(),
			CreatedBy:        entity.Plan.CreatedBy.String(),
			Title:            entity.Plan.Title,
			Description:      entity.Plan.Description,
			GenerationStatus: string(entity.Plan.GenerationStatus),
			CreationType:     string(entity.Plan.CreationType),
			Progress:         entity.Plan.Progress,
			Status:           string(entity.Plan.Status),
			CreatedAt:        entity.Plan.CreatedAt,
			UpdatedAt:        entity.Plan.UpdatedAt,
		},
		Tasks: tasks,
	})
}

// EmployeeGetPlans godoc
// @Summary Получить список планов сотрудника
// @Description Возвращает список планов текущего сотрудника
// @Tags Employee
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.PlanWithTasksResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /employee/plans [get]
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
				ID:               p.Plan.ID.String(),
				EmployeeID:       p.Plan.EmployeeID.String(),
				CreatedBy:        p.Plan.CreatedBy.String(),
				Title:            p.Plan.Title,
				Description:      p.Plan.Description,
				GenerationStatus: string(p.Plan.GenerationStatus),
				CreationType:     string(p.Plan.CreationType),
				Progress:         p.Plan.Progress,
				Status:           string(p.Plan.Status),
				CreatedAt:        p.Plan.CreatedAt,
				UpdatedAt:        p.Plan.UpdatedAt,
			},
			Tasks: tasks,
		})
	}
	c.JSON(http.StatusOK, response)
}

// EmployeeGetPlan godoc
// @Summary Получить план сотрудника по ID
// @Description Возвращает детальную информацию о плане сотрудника
// @Tags Employee
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param plan_id path string true "ID плана"
// @Success 200 {object} dto.PlanWithTasksResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /employee/plans/{plan_id} [get]
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
			ID:               planEntity.Plan.ID.String(),
			EmployeeID:       planEntity.Plan.EmployeeID.String(),
			CreatedBy:        planEntity.Plan.CreatedBy.String(),
			Title:            planEntity.Plan.Title,
			Description:      planEntity.Plan.Description,
			GenerationStatus: string(planEntity.Plan.GenerationStatus),
			CreationType:     string(planEntity.Plan.CreationType),
			Progress:         planEntity.Plan.Progress,
			Status:           string(planEntity.Plan.Status),
			CreatedAt:        planEntity.Plan.CreatedAt,
			UpdatedAt:        planEntity.Plan.UpdatedAt,
		},
		Tasks: tasks,
	})
}

// ListByManager godoc
// @Summary Получить список планов менеджера
// @Description Возвращает список всех планов, созданных менеджером
// @Tags Plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.ListPlansResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /manager/plans [get]
func (h *PlanHandler) ListByManager(c *gin.Context) {
	managerID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	plans, err := h.service.ListByManager(c.Request.Context(), managerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	response := make([]dto.PlanResponse, 0, len(plans))
	for _, p := range plans {
		response = append(response, dto.PlanResponse{
			ID:               p.ID.String(),
			EmployeeID:       p.EmployeeID.String(),
			CreatedBy:        p.CreatedBy.String(),
			Title:            p.Title,
			Description:      p.Description,
			GenerationStatus: string(p.GenerationStatus),
			CreationType:     string(p.CreationType),
			Progress:         p.Progress,
			Status:           string(p.Status),
			CreatedAt:        p.CreatedAt,
			UpdatedAt:        p.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, dto.ListPlansResponse{
		Plans: response,
	})
}

// Update godoc
// @Summary Обновить план
// @Description Обновляет название и описание плана
// @Tags Plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param plan_id path string true "ID плана"
// @Param request body dto.UpdatePlanRequest true "Данные для обновления"
// @Success 200 "OK"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /manager/plans/{plan_id} [patch]
func (h *PlanHandler) Update(c *gin.Context) {
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

	var req dto.UpdatePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err = h.service.Update(c.Request.Context(), plan.UpdatePlanInput{
		PlanID:      planID,
		ManagerID:   managerID,
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		switch {
		case errors.Is(err, plan.ErrPlanNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, plan.ErrInvalidTitle):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, plan.ErrManagerForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.Status(http.StatusOK)
}

// Delete godoc
// @Summary Удалить план
// @Description Удаляет план и все связанные задачи
// @Tags Plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param plan_id path string true "ID плана"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /manager/plans/{plan_id} [delete]
func (h *PlanHandler) Delete(c *gin.Context) {
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

	err = h.service.Delete(c.Request.Context(), managerID, planID)
	if err != nil {
		switch {
		case errors.Is(err, plan.ErrPlanNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, plan.ErrManagerForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// CreateAI godoc
// @Summary Создать план с помощью ИИ
// @Description Генерирует и создает план развития с помощью ИИ
// @Tags Plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateAIPlanRequest true "Данные для генерации плана"
// @Success 201 {object} dto.CreateAIPlanResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /manager/plans/ai [post]
func (h *PlanHandler) CreateAI(c *gin.Context) {
	managerID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req dto.CreateAIPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	employeeID, err := uuid.Parse(req.EmployeeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}

	id, err := h.service.CreateAI(c.Request.Context(), plan.CreateAIInput{
		EmployeeID:  employeeID,
		CreatedBy:   managerID,
		Topic:       req.Topic,
		Description: req.Description,
	},
	)
	if err != nil {
		switch {
		case errors.Is(err, plan.ErrInvalidEmployee):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, plan.ErrInvalidCreator),
			errors.Is(err, plan.ErrEmployeeNotAssigned):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, dto.CreateAIPlanResponse{ID: id.String()})
}

// Archive godoc
// @Summary Архивировать план
// @Description Переводит план в статус archived
// @Tags Plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param plan_id path string true "ID плана"
// @Success 200 "OK"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /manager/plans/{plan_id}/archive [patch]
func (h *PlanHandler) Archive(c *gin.Context) {
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

	err = h.service.Archive(c.Request.Context(), managerID, planID)
	if err != nil {
		switch {
		case errors.Is(err, plan.ErrPlanNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, plan.ErrManagerForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.Status(http.StatusOK)
}
