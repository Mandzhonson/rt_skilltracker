package domain

import (
	"time"

	"github.com/google/uuid"
)

type Skill struct {
	ID          uuid.UUID
	Name        string
	Category    string
	Description *string
	CreatedAt   time.Time
}

func NewSkill(name string, category string, description *string) *Skill {
	return &Skill{
		Name:        name,
		Category:    category,
		Description: description,
	}
}
