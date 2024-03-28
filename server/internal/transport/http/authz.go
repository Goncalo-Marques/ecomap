package http

import (
	"net/http"

	"github.com/goncalo-marques/ecomap/server/internal/authn"
	"github.com/goncalo-marques/ecomap/server/internal/authz"
)

// API roles.
const (
	roleUser          = string(authn.SubjectRoleUser)
	roleWasteOperator = string(authn.SubjectRoleWasteOperator)
	roleManager       = string(authn.SubjectRoleManager)
)

// authzRoleMap defines the authorization role map for the API.
var authzRoleMap authz.RoleMap = authz.RoleMap{
	"/api/employees/signin": {
		http.MethodPost: []string{},
	},
	"/api/employees/{employeeId}": {
		http.MethodPost: []string{roleWasteOperator, roleManager},
	},
}

// AuthzRoles defines the authorization roles for the API.
var AuthzRoles authz.Roles = authz.Roles{
	RoleMap:        authzRoleMap,
	AuthzWildcards: []string{"{employeeId}", "{userId}"},
	AdminRole:      roleManager,
}
