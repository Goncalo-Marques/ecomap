package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// User errors.
var (
	ErrUserAlreadyExists = errors.New("username already exists") // Returned when a user already exists with the same username.
	ErrUserNotFound      = errors.New("user not found")          // Returned when a user is not found.
)

// EditableUser defines the editable user structure.
type EditableUser struct {
	Username  Username
	FirstName Name
	LastName  Name
}

// EditableUserWithPassword defines the editable user structure with a password.
type EditableUserWithPassword struct {
	EditableUser
	Password
}

// EditableUserPatch defines the patchable user structure.
type EditableUserPatch struct {
	Username  *Username
	FirstName *Name
	LastName  *Name
}

// User defines the user structure.
type User struct {
	EditableUser
	ID           uuid.UUID
	CreatedTime  time.Time
	ModifiedTime time.Time
}

// UserSort defines the field of the user to sort.
type UserSort string

const (
	UserSortUsername     UserSort = "username"
	UserSortFirstName    UserSort = "firstName"
	UserSortLastName     UserSort = "lastName"
	UserSortCreatedTime  UserSort = "createdTime"
	UserSortModifiedTime UserSort = "modifiedTime"
)

// Field returns the name of the field to sort by.
func (s UserSort) Field() UserSort {
	return s
}

// Valid returns true if the field is valid, false otherwise.
func (s UserSort) Valid() bool {
	switch s {
	case UserSortUsername,
		UserSortFirstName,
		UserSortLastName,
		UserSortCreatedTime,
		UserSortModifiedTime:
		return true
	default:
		return false
	}
}

// UsersFilter defines the users filter structure.
type UsersFilter struct {
	PaginatedRequest[UserSort]
	Username  *Username
	FirstName *Name
	LastName  *Name
}
