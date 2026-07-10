package handler

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/transport/http/dto"
	"core_service/internal/transport/http/middleware"
	"core_service/internal/usecase/admin"
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
}

type AdminHandler struct {
	service AdminService
}

func NewAdminHandler(service AdminService) *AdminHandler {
	return &AdminHandler{
		service: service,
	}
}

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
			ManagerID: uuidPtrToString(userRes.ManagerID),
			CreatedAt: userRes.CreatedAt,
			UpdatedAt: userRes.UpdatedAt,
		},
	)
}

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
