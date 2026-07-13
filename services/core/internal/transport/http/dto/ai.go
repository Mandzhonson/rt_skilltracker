package dto

type GeneratePlanRequest struct {
	Topic       string `json:"topic" binding:"required"`
	Description string `json:"description"`
}
