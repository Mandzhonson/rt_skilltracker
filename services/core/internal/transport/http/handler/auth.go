package handler

import (
	"context"
	"core_service/internal/pkg/jwt"
	"core_service/internal/transport/http/dto"
	"core_service/internal/transport/http/middleware"
	"core_service/internal/usecase/auth"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:generate mockgen -source=auth.go -destination=mocks/mock_auth_handler.go -package=mocks
type AuthService interface {
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

// Login godoc
// @Summary Авторизация
// @Description Авторизация пользователя по email и паролю
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Данные для входа"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/login [post]
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
			slog.Error("login", slog.Any("error", err.Error()))
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

// Refresh godoc
// @Summary Обновить токены
// @Description Обновляет access и refresh токены
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshRequest true "Refresh токен"
// @Success 200 {object} dto.RefreshResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/refresh [post]
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

// Logout godoc
// @Summary Выход из системы
// @Description Инвалидирует refresh токен
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.RefreshRequest true "Refresh токен"
// @Success 204 "No Content"
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/logout [post]
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
