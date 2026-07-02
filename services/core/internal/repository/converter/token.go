package converter

import (
	"core_service/internal/domain"
	"core_service/internal/repository/model"
	"time"
)

func ToRefreshTokenModel(d domain.RefreshToken) model.RefreshTokenModel {
	return model.RefreshTokenModel{
		UserID:    d.UserID,
		JTI:       d.JTI,
		TokenHash: d.TokenHash,
		ExpiresAt: d.ExpiresAt,
		Revoked:   false,
		CreatedAt: time.Now(),
	}
}

func ToRefreshTokenEntity(m model.RefreshTokenModel) domain.RefreshToken {
	return domain.RefreshToken{
		UserID:    m.UserID,
		JTI:       m.JTI,
		TokenHash: m.TokenHash,
		ExpiresAt: m.ExpiresAt,
	}
}
