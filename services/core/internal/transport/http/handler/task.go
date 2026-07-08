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

type TaskService interface {
	Create(ctx context.Context, input task.CreateTaskInput) (uuid.UUID, error)
	GetByID(ctx context.Context, managerID uuid.UUID, id uuid.UUID) (*domain.Task, error)
	Update(ctx context.Context, input task.UpdateTaskInput) error
	Delete(ctx context.Context, managerID uuid.UUID, taskID uuid.UUID) error
}

type TaskHandler struct {
	service TaskService
}

func NewTaskHandler(service TaskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

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

		case errors.Is(err, task.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	c.JSON(http.StatusCreated, dto.CreateTaskResponse{
		ID: id.String(),
	})
}

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
		case errors.Is(err, task.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}
	c.Status(http.StatusNoContent)
}

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
		case errors.Is(err, task.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}
	c.Status(http.StatusOK)
}
