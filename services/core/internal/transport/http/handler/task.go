package handler

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/transport/http/dto"
	"core_service/internal/transport/http/middleware"
	"core_service/internal/usecase/task"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//go:generate mockgen -source=task.go -destination=mocks/mock_task_handler.go -package=mocks
type TaskService interface {
	Create(ctx context.Context, input task.CreateTaskInput) (uuid.UUID, error)
	GetByID(ctx context.Context, managerID uuid.UUID, id uuid.UUID) (*domain.Task, error)
	Update(ctx context.Context, input task.UpdateTaskInput) error
	UpdateStatus(ctx context.Context, input task.UpdateTaskStatusInput) error
	Delete(ctx context.Context, managerID uuid.UUID, taskID uuid.UUID) error
	CompleteTestingTask(ctx context.Context, planID uuid.UUID, userID uuid.UUID) error
}

type TaskHandler struct {
	service TaskService
}

func NewTaskHandler(service TaskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

// Create godoc
// @Summary Создать задачу
// @Description Создает новую задачу в плане
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param plan_id path string true "ID плана"
// @Param request body dto.CreateTaskRequest true "Данные для создания задачи"
// @Success 201 {object} dto.CreateTaskResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse "План архивирован"
// @Failure 500 {object} dto.ErrorResponse
// @Router /manager/plans/{plan_id}/tasks [post]
func (h *TaskHandler) Create(c *gin.Context) {
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

	var req dto.CreateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	id, err := h.service.Create(
		c.Request.Context(),
		task.CreateTaskInput{
			PlanID:      planID,
			CreatedBy:   managerID,
			Title:       req.Title,
			Description: req.Description,
		},
	)

	if err != nil {
		switch {
		case errors.Is(err, task.ErrInvalidTitle):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		case errors.Is(err, task.ErrPlanNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		case errors.Is(err, task.ErrManagerForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})

		case errors.Is(err, task.ErrPlanArchived):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	c.JSON(http.StatusCreated, dto.CreateTaskResponse{
		ID: id.String(),
	})
}

// GetByID godoc
// @Summary Получить задачу по ID
// @Description Возвращает детальную информацию о задаче
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param task_id path string true "ID задачи"
// @Success 200 {object} dto.TaskResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /manager/tasks/{task_id} [get]
func (h *TaskHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("task_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid task id",
		})
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
		case errors.Is(err, task.ErrTaskNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		case errors.Is(err, task.ErrInvalidTaskID):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	c.JSON(http.StatusOK, dto.TaskResponse{
		ID:          entity.ID.String(),
		PlanID:      entity.PlanID.String(),
		Title:       entity.Title,
		Description: entity.Description,
		Position:    entity.Position,
		Status:      string(entity.Status),
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	})
}

// Delete godoc
// @Summary Удалить задачу
// @Description Удаляет задачу
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param task_id path string true "ID задачи"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /manager/tasks/{task_id} [delete]
func (h *TaskHandler) Delete(c *gin.Context) {
	managerID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	taskID, err := uuid.Parse(c.Param("task_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	err = h.service.Delete(c.Request.Context(), managerID, taskID)
	if err != nil {
		switch {
		case errors.Is(err, task.ErrTaskNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, task.ErrManagerForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}
	c.Status(http.StatusNoContent)
}

// Update godoc
// @Summary Обновить задачу
// @Description Обновляет название и описание задачи
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param task_id path string true "ID задачи"
// @Param request body dto.UpdateTaskRequest true "Данные для обновления"
// @Success 200 "OK"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /manager/tasks/{task_id} [patch]
func (h *TaskHandler) Update(c *gin.Context) {
	managerID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	taskID, err := uuid.Parse(c.Param("task_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	var req dto.UpdateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err = h.service.Update(c.Request.Context(), task.UpdateTaskInput{
		TaskID:      taskID,
		ManagerID:   managerID,
		Title:       req.Title,
		Description: req.Description,
	},
	)

	if err != nil {
		switch {
		case errors.Is(err, task.ErrInvalidTitle):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, task.ErrTaskNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, task.ErrManagerForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, task.ErrInvalidUpdate):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}
	c.Status(http.StatusOK)
}

// UpdateStatus godoc
// @Summary Обновить статус задачи
// @Description Изменяет статус задачи (для сотрудника)
// @Tags Employee
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param task_id path string true "ID задачи"
// @Param request body dto.UpdateTaskStatusRequest true "Новый статус"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /employee/tasks/{task_id}/status [patch]
func (h *TaskHandler) UpdateStatus(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	taskID, err := uuid.Parse(c.Param("task_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	var req dto.UpdateTaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err = h.service.UpdateStatus(
		c.Request.Context(),
		task.UpdateTaskStatusInput{
			TaskID: taskID,
			UserID: userID,
			Status: domain.TaskStatus(req.Status),
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, task.ErrInvalidTaskID),
			errors.Is(err, task.ErrInvalidUserID),
			errors.Is(err, task.ErrInvalidStatus):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, task.ErrTaskNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, task.ErrEmployeeForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}
	c.Status(http.StatusNoContent)
}
