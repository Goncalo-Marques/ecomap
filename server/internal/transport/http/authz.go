package http

import (
	"github.com/goncalo-marques/ecomap/server/internal/authn"
	"github.com/goncalo-marques/ecomap/server/internal/authz"
)

// AuthzRoles defines the authorization roles for the API.
var AuthzRoles authz.Roles = authz.Roles{
	AuthzWildcards: []string{"employeeId", "userId"},
	AdminRole:      string(authn.SubjectRoleManager),
}
