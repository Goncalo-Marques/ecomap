package http

import (
	"github.com/goncalo-marques/ecomap/server/internal/authn"
	"github.com/goncalo-marques/ecomap/server/internal/authz"
)

// API roles.
const (
	roleUser          = string(authn.SubjectRoleUser)
	roleWasteOperator = string(authn.SubjectRoleWasteOperator)
	roleManager       = string(authn.SubjectRoleManager)
)

// AuthzRoles defines the authorization roles for the API.
var AuthzRoles authz.Roles = authz.Roles{
	AuthzWildcards: []string{"employeeId", "userId"},
	AdminRole:      roleManager,
}
