package authz

import (
	"github.com/goncalo-marques/ecomap/server/internal/authn"
)

// RoleMap defines the map of roles required to access a method in a given endpoint.
//
// Example:
//
//	var roleMap RoleMap = RoleMap{
//		"/employees/{employeeId}": {
//			http.MethodGet:    []string{"user", "admin"},
//			http.MethodDelete: []string{"admin"},
//		},
//	}
type RoleMap map[string]map[string][]string

// Roles defines the authorization roles structure.
type Roles struct {
	RoleMap        RoleMap  // Map of required roles.
	AuthzWildcards []string // Wildcards that require authorization.
	AdminRole      string   // Role that can access all identifiers (bypasses AuthzWildcards).
}

// AuthenticationService defines the authentication service interface.
type AuthenticationService interface {
	ParseJWT(tokenString string) (authn.Claims, error)
}

// service defines the authorization service structure.
type service struct {
	roles        Roles
	authnService AuthenticationService
}

// New returns a new authorization service.
func New(roles Roles, authnService AuthenticationService) *service {
	return &service{
		roles:        roles,
		authnService: authnService,
	}
}
