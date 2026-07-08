package handler

import "github.com/google/uuid"

func uuidPtrToString(id *uuid.UUID) *string {
	if id == nil {
		return nil
	}

	value := id.String()
	return &value
}
