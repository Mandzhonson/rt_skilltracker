package handler

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/transport/http/dto"
	"core_service/internal/transport/http/middleware"
	"core_service/internal/usecase/user"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(ctx context.Context, u *domain.User) (uuid.UUID, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error)
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{
		service: userService,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	entity := &domain.User{
		Email:        req.Email,
		PasswordHash: req.Password,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
	}

	id, err := h.service.CreateUser(c.Request.Context(), entity)
	if err != nil {

		switch {

		case errors.Is(err, user.ErrUserAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, user.ErrInvalidEmail),
			errors.Is(err, user.ErrInvalidPassword),
			errors.Is(err, user.ErrInvalidName):

			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}

		return
	}

	c.JSON(http.StatusCreated, dto.RegisterResponse{
		ID: id.String(),
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	claims, ok := middleware.GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}

	entity, err := h.service.GetProfile(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, dto.ProfileResponse{
		ID:        entity.ID.String(),
		Email:     entity.Email,
		FirstName: entity.FirstName,
		LastName:  entity.LastName,
		Role:      string(entity.Role),
	})
}
