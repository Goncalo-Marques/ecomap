package domain

import (
	"errors"
	"fmt"
)

// Field constraints.
const (
	usernameMinLength = 1
	usernameMaxLength = 50

	nameMinLength = 1
	nameMaxLength = 50
)

// Common errors.
var (
	ErrCredentialsIncorrect = errors.New("incorrect credentials")  // Returned when a username is not found or the password is incorrect.
	ErrRoadNotFound         = errors.New("road not found")         // Returned when a road is not found.
	ErrMunicipalityNotFound = errors.New("municipality not found") // Returned when a municipality is not found.
)

// Returned when a field contains an invalid value.
type ErrFieldValueInvalid struct {
	FieldName string
}

func (e *ErrFieldValueInvalid) Error() string {
	return fmt.Sprintf("invalid field value: %s", e.FieldName)
}

// Username defines the username type.
type Username string

// Valid returns true if the username is valid, false otherwise.
func (u Username) Valid() bool {
	return len(u) >= usernameMinLength && len(u) <= usernameMaxLength
}

// Name defines the name type.
type Name string

// Valid returns true if the name is valid, false otherwise.
func (n Name) Valid() bool {
	return len(n) >= nameMinLength && len(n) <= nameMaxLength
}

// Password defines the password type.
type Password string

// SignIn defines the sign-in structure.
type SignIn struct {
	Username Username
	Password Password
}
