package handler

import (
	"context"
	"core_service/internal/domain"
	"core_service/internal/pkg/jwt"
	"core_service/internal/transport/http/dto"
	"core_service/internal/transport/http/middleware"
	"core_service/internal/usecase/auth"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthService interface {
	CreateUser(ctx context.Context, u *domain.User) (uuid.UUID, error)
	Login(ctx context.Context, email, password string) (string, string, error)
	Refresh(ctx context.Context, refreshToken string) (string, string, error)
	Logout(ctx context.Context, claimsAccess *jwt.Claims, refreshToken string) error
}

type AuthHandler struct {
	service AuthService
}

func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) CreateUser(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	user := &domain.User{
		Email:        req.Email,
		PasswordHash: req.Password,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
	}

	id, err := h.service.CreateUser(c.Request.Context(), user)
	if err != nil {

		switch {

		case errors.Is(err, auth.ErrUserAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})

		case errors.Is(err, auth.ErrInvalidEmail),
			errors.Is(err, auth.ErrInvalidPassword),
			errors.Is(err, auth.ErrInvalidName):

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

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	accessToken, refreshToken, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})

		}
		return
	}
	c.JSON(http.StatusOK, dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	accessToken, refreshToken, err := h.service.Refresh(
		c.Request.Context(),
		req.RefreshToken,
	)
	if err != nil {

		switch {
		case errors.Is(err, auth.ErrInvalidCredentials),
			errors.Is(err, jwt.ErrTokenExpired),
			errors.Is(err, jwt.ErrTokenInvalid),
			errors.Is(err, auth.ErrInvalidRefreshToken):

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			return
		}
	}

	c.JSON(http.StatusOK, dto.RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.RefreshRequest

	claims, ok := middleware.GetClaims(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	err := h.service.Logout(
		c.Request.Context(),
		claims,
		req.RefreshToken,
	)

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired),
			errors.Is(err, jwt.ErrTokenInvalid),
			errors.Is(err, auth.ErrInvalidCredentials):

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
