package dto

import "time"

type SkillResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type ListSkillsResponse struct {
	Skills []SkillResponse `json:"skills"`
}
