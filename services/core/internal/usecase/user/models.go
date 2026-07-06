package user

import (
	"io"

	"github.com/google/uuid"
)

type UpdateProfileInput struct {
	UserID    uuid.UUID
	Email     *string
	FirstName *string
	LastName  *string
}

type SetAvatarInput struct {
	UserID      uuid.UUID
	File        io.Reader
	Size        int64
	ContentType string
}
