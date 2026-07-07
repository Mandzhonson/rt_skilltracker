package dto

import "time"

type ListUsersResponse struct {
	Users []UserResponse `json:"users"`
}

type UserResponse struct {
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Role      string  `json:"role"`
	ManagerID *string `json:"manager_id,omitempty"`
}

type UserDetailsResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`

	Role string `json:"role"`

	ManagerID *string `json:"manager_id,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

type AssignManagerRequest struct {
	ManagerID string `json:"manager_id" binding:"required,uuid"`
}
