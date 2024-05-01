package http

import (
	"github.com/goncalo-marques/ecomap/server/internal/authn"
	"github.com/goncalo-marques/ecomap/server/internal/authz"
	"github.com/goncalo-marques/ecomap/server/internal/domain"
)

// AuthzRoles defines the authorization roles for the API.
var AuthzRoles authz.Roles = authz.Roles{
	AuthzWildcards: []string{domain.FieldParamEmployeeID, domain.FieldParamUserID},
	AdminRole:      string(authn.SubjectRoleManager),
}
