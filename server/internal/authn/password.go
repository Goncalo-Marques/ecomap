package authn

import (
	"errors"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

const (
	passwordMinLength       = 14
	passwordMinDigits       = 1
	passwordMinSpecialChars = 1
)

// ValidPassword returns true if the password matches the requirements, false otherwise.
func (s *service) ValidPassword(password string) bool {
	var digits int
	var specialChars int
	for _, r := range password {
		if unicode.IsDigit(r) {
			digits++
			continue
		}
		if !unicode.IsLetter(r) {
			specialChars++
			continue
		}
	}

	return len(password) >= passwordMinLength && digits >= passwordMinDigits && specialChars >= passwordMinSpecialChars
}

// HashPassword returns the hashed password.
func (s *service) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash returns true if the password matches the hashed password, false otherwise.
func (s *service) CheckPasswordHash(password, hash string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
