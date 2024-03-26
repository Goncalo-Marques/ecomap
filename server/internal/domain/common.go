package domain

import (
	"errors"
)

// Common errors.
var (
	ErrCredentialsIncorrect = errors.New("incorrect credentials") // Returned when a username is not found or the password is incorrect.
)

// SignIn defines the sign-in structure.
type SignIn struct {
	Username string
	Password string
}
