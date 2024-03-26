package authn

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	jwtIssuer         = "server"
	jwtExpirationTime = time.Hour * 24
)

var ErrJWTInvalid = errors.New("invalid jwt") // Returned when the JWT is invalid.

// SubjectRole defines the role of subject.
type SubjectRole string

const (
	SubjectRoleUser          SubjectRole = "user"
	SubjectRoleWasteOperator SubjectRole = "waste_operator"
	SubjectRoleManager       SubjectRole = "manager"
)

// Claims defines the JSON Web Token claims structure.
type Claims struct {
	jwt.RegisteredClaims
	Roles []SubjectRole `json:"roles,omitempty"`
}

// NewJWT returns a new signed JSON Web Token with an expiration time of 24 hours and the specified claims.
func (s *service) NewJWT(subject string, subjectRoles []SubjectRole) (string, error) {
	expiresAt := time.Now().Add(jwtExpirationTime).UTC()
	issuedAt := time.Now().UTC()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jwtIssuer,
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
		},
		Roles: subjectRoles,
	})

	tokenString, err := token.SignedString(s.jwtSigningKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseJWT parses the given token and returns the associated subject.
// If the token is invalid, ErrJWTInvalid is returned.
func (s *service) ParseJWT(tokenString string) (Claims, error) {
	var claims Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
		return s.jwtSigningKey, nil
	})
	if err != nil {
		return Claims{}, err
	}

	if !token.Valid {
		return Claims{}, ErrJWTInvalid
	}

	return claims, nil
}
