package domain

import "time"

type User struct {
	ID           string
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	Role         string
	ManagerID    *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
