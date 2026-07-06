package converter

import (
	"core_service/internal/domain"
	"core_service/internal/repository/model"
)

func ToUserEntity(m *model.UserModel) *domain.User {
	return &domain.User{
		ID:           m.ID,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
		FirstName:    m.FirstName,
		LastName:     m.LastName,
		AvatarKey:    m.AvatarKey,
		Role:         m.Role,
		ManagerID:    m.ManagerID,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func ToUserModel(m *domain.User) *model.UserModel {
	return &model.UserModel{
		ID:           m.ID,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
		FirstName:    m.FirstName,
		LastName:     m.LastName,
		AvatarKey:    m.AvatarKey,
		Role:         m.Role,
		ManagerID:    m.ManagerID,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}
