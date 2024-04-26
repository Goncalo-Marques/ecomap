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
	ID         uuid.UUID
	CreatedAt  time.Time
	ModifiedAt time.Time
}

// UserPaginatedSort defines the field of the user to sort.
type UserPaginatedSort string

const (
	UserPaginatedSortUsername   UserPaginatedSort = "Username"
	UserPaginatedSortFirstName  UserPaginatedSort = "FirstName"
	UserPaginatedSortLastName   UserPaginatedSort = "LastName"
	UserPaginatedSortCreatedAt  UserPaginatedSort = "CreatedAt"
	UserPaginatedSortModifiedAt UserPaginatedSort = "ModifiedAt"
)

// Field returns the name of the field to sort by.
func (s UserPaginatedSort) Field() UserPaginatedSort {
	return s
}

// Valid returns true if the field is valid, false otherwise.
func (s UserPaginatedSort) Valid() bool {
	switch s {
	case UserPaginatedSortUsername,
		UserPaginatedSortFirstName,
		UserPaginatedSortLastName,
		UserPaginatedSortCreatedAt,
		UserPaginatedSortModifiedAt:
		return true
	default:
		return false
	}
}

// UsersPaginatedFilter defines the users filter structure.
type UsersPaginatedFilter struct {
	PaginatedRequest[UserPaginatedSort]
	Username  *Username
	FirstName *Name
	LastName  *Name
}
