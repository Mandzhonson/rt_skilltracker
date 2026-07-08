package domain

import (
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	TaskTodo       TaskStatus = "todo"
	TaskInProgress TaskStatus = "in_progress"
	TaskDone       TaskStatus = "done"
)

type Task struct {
	ID          uuid.UUID
	PlanID      uuid.UUID
	Title       string
	Description *string
	Position    int
	Status      TaskStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewTask(planID uuid.UUID, title string, description *string, position int) *Task {
	return &Task{
		PlanID:      planID,
		Title:       title,
		Description: description,
		Position:    position,
		Status:      TaskTodo,
	}
}
