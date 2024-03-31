package authn

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	jwtIssuer         = "server"
	jwtExpirationTime = time.Hour * 24

	descriptionFailedToParseJWTWithClaims = "authn: failed to parse jwt with claims"
)

// SubjectRole defines the role of the subject.
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
func (s *service) NewJWT(subject string, subjectRoles ...SubjectRole) (string, error) {
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

// ParseJWT parses the given token and returns the associated subject, or returns an error if the token is invalid.
func (s *service) ParseJWT(tokenString string) (Claims, error) {
	var claims Claims
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
		return s.jwtSigningKey, nil
	})
	if err != nil {
		return Claims{}, fmt.Errorf("%s: %w", descriptionFailedToParseJWTWithClaims, err)
	}

	return claims, nil
}
