package authn

import (
	"errors"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

const (
	passwordMinLength       = 14
	passwordMaxLength       = 72
	passwordMinDigits       = 1
	passwordMinChars        = 1
	passwordMinSpecialChars = 1
)

// ValidPassword returns true if the password matches the requirements, false otherwise.
func (s *service) ValidPassword(password []byte) bool {
	// Check password length requirements.
	if len(password) < passwordMinLength || len(password) > passwordMaxLength {
		return false
	}

	// Check password complexity requirements.
	var digits int
	var chars int
	var specialChars int
	for _, r := range string(password) {
		if unicode.IsDigit(r) {
			digits++
			continue
		}
		if unicode.IsLetter(r) {
			chars++
		} else {
			specialChars++
		}
	}

	return digits >= passwordMinDigits && chars >= passwordMinChars && specialChars >= passwordMinSpecialChars
}

// HashPassword returns the hashed password.
func (s *service) HashPassword(password []byte) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	return bytes, err
}

// CheckPasswordHash returns true if the password matches the hashed password, false otherwise.
func (s *service) CheckPasswordHash(password, hash []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, password)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
