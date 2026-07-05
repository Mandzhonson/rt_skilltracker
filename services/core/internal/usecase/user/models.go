package user

import "github.com/google/uuid"

type UpdateProfileInput struct {
	UserID uuid.UUID

	Email     *string
	FirstName *string
	LastName  *string
}
