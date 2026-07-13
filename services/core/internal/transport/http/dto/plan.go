package dto

import "time"

type CreatePlanRequest struct {
	EmployeeID  string  `json:"employee_id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
}

type CreatePlanResponse struct {
	ID string `json:"id"`
}

type PlanResponse struct {
	ID           string    `json:"id"`
	EmployeeID   string    `json:"employee_id"`
	CreatedBy    string    `json:"created_by"`
	Title        string    `json:"title"`
	Description  *string   `json:"description"`
	CreationType string    `json:"creation_type"`
	Progress     int       `json:"progress"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type PlanWithTasksResponse struct {
	Plan  PlanResponse   `json:"plan"`
	Tasks []TaskResponse `json:"tasks"`
}

type UpdatePlanRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description"`
}

type CreateAIPlanRequest struct {
	EmployeeID  string `json:"employee_id" binding:"required"`
	Topic       string `json:"topic" binding:"required"`
	Description string `json:"description"`
}

type CreateAIPlanResponse struct {
	ID string `json:"id"`
}
