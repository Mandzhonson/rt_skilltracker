package handler

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/transport/http/dto"
	"core_service/internal/transport/http/middleware"
	"core_service/internal/usecase/admin"
	"core_service/internal/usecase/task"
	"core_service/internal/usecase/user"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminService interface {
	ListUsers(ctx context.Context, input admin.ListUsersInput) ([]*domain.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error)
	UpdateRole(ctx context.Context, input admin.UpdateRoleInput) error
	AssignManager(ctx context.Context, input admin.AssignManagerInput) error
	RemoveManager(ctx context.Context, userID uuid.UUID) error
	ListEmployeesByManager(ctx context.Context, managerID uuid.UUID) ([]*domain.User, error)
	GetUserAvatar(ctx context.Context, userID uuid.UUID) (io.ReadCloser, string, error)
	UpdatePosition(ctx context.Context, userID uuid.UUID, position string) error
	GetEmployeeProfileForAdmin(ctx context.Context, employeeID uuid.UUID) (*domain.EmployeeProfile, error)
	AdminGetPlan(ctx context.Context, planID uuid.UUID) (*domain.PlanWithTasks, error)
}

type AdminHandler struct {
	service AdminService
}

func NewAdminHandler(service AdminService) *AdminHandler {
	return &AdminHandler{
		service: service,
	}
}

// ListUsers godoc
// @Summary Получить список пользователей
// @Description Возвращает список всех пользователей с пагинацией и фильтрацией
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество на странице" default(20)
// @Param role query string false "Фильтр по роли" Enums(admin, manager, employee)
// @Param search query string false "Поиск по email или имени"
// @Success 200 {object} dto.ListUsersResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/users [get]
func (h *AdminHandler) ListUsers(c *gin.Context) {

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	var role *domain.Role

	if value := c.Query("role"); value != "" {
		r := domain.Role(value)

		if !r.IsValid() {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid role",
			})
			return
		}

		role = &r
	}
	var search *string

	if value := c.Query("search"); value != "" {
		search = &value
	}

	users, err := h.service.ListUsers(
		c.Request.Context(),
		admin.ListUsersInput{
			Page:   page,
			Limit:  limit,
			Role:   role,
			Search: search,
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	response := dto.ListUsersResponse{
		Users: make([]dto.UserResponse, 0, len(users)),
	}

	for _, user := range users {
		response.Users = append(
			response.Users,
			dto.UserResponse{
				ID:        user.ID.String(),
				Email:     user.Email,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Role:      string(user.Role),
				ManagerID: uuidPtrToString(user.ManagerID),
			},
		)
	}

	c.JSON(http.StatusOK, response)
}

// GetUser godoc
// @Summary Получить пользователя по ID
// @Description Возвращает детальную информацию о пользователе
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID пользователя"
// @Success 200 {object} dto.UserDetailsResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/users/{id} [get]
func (h *AdminHandler) GetUser(c *gin.Context) {

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	userRes, err := h.service.GetUser(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(
		http.StatusOK,
		dto.UserDetailsResponse{
			ID:        userRes.ID.String(),
			Email:     userRes.Email,
			FirstName: userRes.FirstName,
			LastName:  userRes.LastName,
			Role:      string(userRes.Role),
			Position:  userRes.Position,
			ManagerID: uuidPtrToString(userRes.ManagerID),
			CreatedAt: userRes.CreatedAt,
			UpdatedAt: userRes.UpdatedAt,
		},
	)
}

// UpdateRole godoc
// @Summary Обновить роль пользователя
// @Description Изменяет роль пользователя (только для админов)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID пользователя"
// @Param request body dto.UpdateRoleRequest true "Новая роль"
// @Success 200 "OK"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse "Нельзя удалить последнего админа"
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/users/{id}/role [patch]
func (h *AdminHandler) UpdateRole(c *gin.Context) {
	actorID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	err = h.service.UpdateRole(
		c.Request.Context(),
		admin.UpdateRoleInput{
			ActorID: actorID,
			UserID:  userID,
			Role:    domain.Role(req.Role),
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, admin.ErrInvalidRole):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, user.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, admin.ErrChangeOwnRole):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, admin.ErrLastAdminProtected):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.Status(http.StatusOK)
}

// AssignManager godoc
// @Summary Назначить менеджера сотруднику
// @Description Назначает менеджера для указанного сотрудника
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID сотрудника"
// @Param request body dto.AssignManagerRequest true "ID менеджера"
// @Success 200 "OK"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/users/{id}/manager [patch]
func (h *AdminHandler) AssignManager(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	var req dto.AssignManagerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	managerID, err := uuid.Parse(req.ManagerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid manager id"})
		return
	}

	err = h.service.AssignManager(c.Request.Context(), admin.AssignManagerInput{UserID: userID, ManagerID: managerID})

	if err != nil {

		switch {

		case errors.Is(err, admin.ErrAssignYourself),
			errors.Is(err, admin.ErrInvalidManager),
			errors.Is(err, admin.ErrManagerCycle):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, user.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}

		return
	}

	c.Status(http.StatusOK)
}

// RemoveManager godoc
// @Summary Удалить менеджера у сотрудника
// @Description Удаляет назначенного менеджера у сотрудника
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID сотрудника"
// @Success 200 "OK"
// @Failure 204 "No Content (менеджер не был назначен)"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/users/{id}/manager [delete]
func (h *AdminHandler) RemoveManager(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	err = h.service.RemoveManager(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		case errors.Is(err, admin.ErrManagerNotAssigned):
			c.Status(http.StatusNoContent)

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.Status(http.StatusOK)
}

// ListEmployeesByManager godoc
// @Summary Получить список сотрудников менеджера
// @Description Возвращает список сотрудников, закрепленных за менеджером
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID менеджера"
// @Success 200 {array} dto.UserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/managers/{id}/employees [get]
func (h *AdminHandler) ListEmployeesByManager(c *gin.Context) {
	managerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid manager id"})
		return
	}

	users, err := h.service.ListEmployeesByManager(c.Request.Context(), managerID)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		case errors.Is(err, admin.ErrInvalidManager):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	resp := make([]dto.UserResponse, 0, len(users))

	for _, u := range users {
		resp = append(resp, dto.UserResponse{
			ID:        u.ID.String(),
			Email:     u.Email,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Role:      string(u.Role),
		})
	}

	c.JSON(http.StatusOK, resp)
}

// GetUserAvatar godoc
// @Summary Получить аватар пользователя
// @Description Возвращает аватар пользователя
// @Tags Admin
// @Accept json
// @Produce image/jpeg, image/png, image/gif
// @Security BearerAuth
// @Param id path string true "ID пользователя"
// @Success 200 {file} file "Аватар пользователя"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/users/{id}/avatar [get]
func (h *AdminHandler) GetUserAvatar(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	reader, contentType, err := h.service.GetUserAvatar(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrAvatarNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "avatar not found"})
		case errors.Is(err, user.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}
	defer reader.Close()

	c.Header("Content-Type", contentType)
	_, err = io.Copy(c.Writer, reader)
	if err != nil {
		return
	}
}

// UpdatePosition godoc
// @Summary Обновить должность пользователя
// @Description Изменяет должность пользователя
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID пользователя"
// @Param request body dto.UpdatePositionRequest true "Новая должность"
// @Success 200 "OK"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/users/{id}/position [patch]
func (h *AdminHandler) UpdatePosition(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req dto.UpdatePositionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err = h.service.UpdatePosition(c.Request.Context(), id, req.Position)
	if err != nil {
		switch {
		case errors.Is(err, admin.ErrInvalidPosition):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, user.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}
	c.Status(http.StatusOK)
}

// GetEmployeeProfile godoc
// @Summary Получить профиль сотрудника для админа
// @Description Возвращает расширенный профиль сотрудника с навыками и планами
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID сотрудника"
// @Success 200 {object} dto.EmployeeProfileResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/users/{id}/profile [get]
func (h *AdminHandler) GetEmployeeProfile(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	profile, err := h.service.GetEmployeeProfileForAdmin(
		c.Request.Context(),
		id,
	)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		case errors.Is(err, admin.ErrInvalidEmployee):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	plans := make([]dto.PlanResponse, 0, len(profile.Plans))
	for _, p := range profile.Plans {
		plans = append(plans, dto.PlanResponse{
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

	skills := make([]dto.SkillResponse, 0, len(profile.Skills))
	for _, s := range profile.Skills {
		skills = append(skills, dto.SkillResponse{
			ID:          s.ID.String(),
			Name:        s.Name,
			Category:    s.Category,
			Description: s.Description,
		})
	}

	c.JSON(http.StatusOK, dto.EmployeeProfileResponse{
		User: dto.UserResponse{
			ID:        profile.User.ID.String(),
			Email:     profile.User.Email,
			FirstName: profile.User.FirstName,
			LastName:  profile.User.LastName,
			Role:      string(profile.User.Role),
			Position:  profile.User.Position,
		},
		Skills: skills,
		Plans:  plans,
	})
}

// GetPlan godoc
// @Summary Получить план по ID (админ)
// @Description Возвращает детальную информацию о плане с задачами
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param plan_id path string true "ID плана"
// @Success 200 {object} dto.PlanWithTasksResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/plans/{plan_id} [get]
func (h *AdminHandler) GetPlan(c *gin.Context) {
	planID, err := uuid.Parse(c.Param("plan_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan id"})
		return
	}

	plan, err := h.service.AdminGetPlan(c.Request.Context(), planID)
	if err != nil {
		switch {
		case errors.Is(err, task.ErrPlanNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	tasks := make([]dto.TaskResponse, 0, len(plan.Tasks))
	for _, task := range plan.Tasks {
		tasks = append(tasks, dto.TaskResponse{
			ID:          task.ID.String(),
			Title:       task.Title,
			Description: task.Description,
			Position:    task.Position,
			Status:      string(task.Status),
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, dto.PlanWithTasksResponse{
		Plan: dto.PlanResponse{
			ID:               plan.Plan.ID.String(),
			EmployeeID:       plan.Plan.EmployeeID.String(),
			CreatedBy:        plan.Plan.CreatedBy.String(),
			Title:            plan.Plan.Title,
			Description:      plan.Plan.Description,
			GenerationStatus: string(plan.Plan.GenerationStatus),
			CreationType:     string(plan.Plan.CreationType),
			Progress:         plan.Plan.Progress,
			Status:           string(plan.Plan.Status),
			CreatedAt:        plan.Plan.CreatedAt,
			UpdatedAt:        plan.Plan.UpdatedAt,
		},
		Tasks: tasks,
	})
}
