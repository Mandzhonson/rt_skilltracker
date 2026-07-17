package jwt

import (
	"core_service/internal/config"
	"core_service/internal/domain"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrTokenExpired = errors.New("token expired")
	ErrTokenInvalid = errors.New("token invalid")
)

type JWTService struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewJWTService(cfg config.JWTConfig) *JWTService {
	return &JWTService{
		accessSecret:  []byte(cfg.AccessSecret),
		refreshSecret: []byte(cfg.RefreshSecret),
		accessTTL:     cfg.AccessTTL,
		refreshTTL:    cfg.RefreshTTL,
	}
}

type Claims struct {
	jwt.RegisteredClaims
	Role domain.Role `json:"role,omitempty"`
}

func (c *JWTService) GenerateAccessToken(userID uuid.UUID, role domain.Role) (string, error) {
	jti := uuid.NewString()

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			ID:        jti,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(c.accessTTL)),
		},
		Role: role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(c.accessSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (c *JWTService) GenerateRefreshToken(userID uuid.UUID) (string, string, error) {
	now := time.Now()
	jti := uuid.NewString()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			ID:        jti,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(c.refreshTTL)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(c.refreshSecret)
	if err != nil {
		return "", "", err
	}
	return tokenString, jti, nil
}

func (c *JWTService) ParseAccessToken(token string) (*Claims, error) {
	return c.parse(token, c.accessSecret)
}

func (c *JWTService) ParseRefreshToken(token string) (*Claims, error) {
	return c.parse(token, c.refreshSecret)
}

func (c *JWTService) parse(tokenString string, secret []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrTokenInvalid
		}
		return secret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}
	return claims, nil
}

func (j *JWTService) RefreshTTL() time.Duration {
	return j.refreshTTL
}
