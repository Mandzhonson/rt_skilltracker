package dto

import "time"

type CreateTaskRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description"`
}

type CreateTaskResponse struct {
	ID string `json:"id"`
}

type TaskResponse struct {
	ID          string    `json:"id"`
	PlanID      string    `json:"plan_id"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	Position    int       `json:"position"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateTaskRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

type UpdateTaskStatusRequest struct {
	Status string `json:"status" binding:"required"`
}
