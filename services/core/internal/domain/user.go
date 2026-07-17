package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleEmployee Role = "employee"
	RoleManager  Role = "manager"
	RoleAdmin    Role = "admin"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	AvatarKey    *string
	Role         Role
	Position     string
	ManagerID    *uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) IsManager() bool {
	return u.Role == RoleManager
}

func (u *User) IsEmployee() bool {
	return u.Role == RoleEmployee
}

func (u *User) AssignManager(managerID uuid.UUID) {
	u.ManagerID = &managerID
}

func (u *User) RemoveManager() {
	u.ManagerID = nil
}

func (u *User) SetRole(role Role) {
	u.Role = role
}

func (r Role) IsValid() bool {
	switch r {
	case RoleEmployee,
		RoleManager,
		RoleAdmin:
		return true
	default:
		return false
	}
}

func NewEmployee(email, passwordHash, firstName, lastName, position string) *User {
	return &User{
		Email:        email,
		PasswordHash: passwordHash,
		FirstName:    firstName,
		LastName:     lastName,
		Role:         RoleEmployee,
		Position:     position,
	}
}

func NewManager(email, passwordHash, firstName, lastName string) *User {
	return &User{
		Email:        email,
		PasswordHash: passwordHash,
		FirstName:    firstName,
		LastName:     lastName,
		Role:         RoleManager,
	}
}

func NewAdmin(email, passwordHash, firstName, lastName string) *User {
	return &User{
		Email:        email,
		PasswordHash: passwordHash,
		FirstName:    firstName,
		LastName:     lastName,
		Role:         RoleAdmin,
	}
}

type UpdateUserProfile struct {
	Email     *string
	FirstName *string
	LastName  *string
}
