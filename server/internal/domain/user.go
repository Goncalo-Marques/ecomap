package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// User errors.
var (
	ErrUserNotFound = errors.New("user not found") // Returned when a user is not found.
)

// EditableUser defines the editable user structure.
type EditableUser struct {
	FirstName string
	LastName  string
}

// User defines the user structure.
type User struct {
	EditableUser
	ID           uuid.UUID
	CreatedTime  time.Time
	ModifiedTime time.Time
}
