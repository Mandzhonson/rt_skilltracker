package handler

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/transport/http/dto"
	"core_service/internal/transport/http/middleware"
	"core_service/internal/usecase/user"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(ctx context.Context, input user.CreateUserInput) (uuid.UUID, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	UpdateProfile(ctx context.Context, upd user.UpdateProfileInput) (*domain.User, error)
	UpdatePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
	SetAvatar(ctx context.Context, input user.SetAvatarInput) error
	GetAvatar(ctx context.Context, userID uuid.UUID) (io.ReadCloser, string, error)
	DeleteAvatar(ctx context.Context, userID uuid.UUID) error
	CreateAdmin(ctx context.Context, input user.CreateUserInput) error
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

	input := user.CreateUserInput{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	id, err := h.service.CreateUser(c.Request.Context(), input)
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

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	var req dto.UpdateProfileRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	entity, err := h.service.UpdateProfile(
		c.Request.Context(),
		user.UpdateProfileInput{
			UserID:    userID,
			Email:     req.Email,
			FirstName: req.FirstName,
			LastName:  req.LastName,
		},
	)
	if err != nil {

		switch {

		case errors.Is(err, user.ErrInvalidEmail),
			errors.Is(err, user.ErrInvalidName):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, user.ErrUserAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, user.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, user.ErrNoContent):
			c.Status(http.StatusNoContent)
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
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

func (h *UserHandler) UpdatePassword(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req dto.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.service.UpdatePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		switch {
		case errors.Is(err, user.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, user.ErrInvalidPassword),
			errors.Is(err, user.ErrInvalidCredentials):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}
	c.Status(http.StatusOK)
}

func (h *UserHandler) SetAvatar(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "file is required",
		})
		return
	}
	defer file.Close()

	input := user.SetAvatarInput{
		UserID:      userID,
		File:        file,
		Size:        header.Size,
		ContentType: header.Header.Get("Content-Type"),
	}

	err = h.service.SetAvatar(c.Request.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrAvatarTooLarge),
			errors.Is(err, user.ErrInvalidAvatarFormat):

			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

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

	c.Status(http.StatusOK)
}

func (h *UserHandler) DeleteAvatar(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	err := h.service.DeleteAvatar(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrNoContent):
			c.Status(http.StatusNoContent)

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

	c.Status(http.StatusOK)
}

func (h *UserHandler) GetAvatar(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	reader, contentType, err := h.service.GetAvatar(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrAvatarNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
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
	defer reader.Close()

	c.Header("Content-Type", contentType)

	_, err = io.Copy(c.Writer, reader)
	if err != nil {
		return
	}
}
