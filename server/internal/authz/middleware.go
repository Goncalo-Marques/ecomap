package authz

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	spec "github.com/goncalo-marques/ecomap/server/api/ecomap"
)

const (
	requestHeaderKeyAuthorization = "Authorization"
	bearerSchemaPrefix            = "Bearer "

	descriptionFailedToFindBearerSchemaPrefix = "authz: failed to find bearer schema prefix"
	descriptionFailedToParseJWT               = "authz: failed to parse jwt"
	descriptionFailedToValidateRoles          = "authz: failed to validate roles"
)

var (
	ErrAuthorizationHeaderInvalid = errors.New("invalid authorization header") // Returned when the Authorization header is invalid.
	ErrJWTInvalid                 = errors.New("invalid jwt")                  // Returned when the JWT is invalid.
	ErrRolesInvalid               = errors.New("invalid roles")                // Returned when the subject does not contain any of the required roles.
	ErrAuthorizationInvalid       = errors.New("invalid authorization")        // Returned when the subject is not the one actually contained in the path wildcard.
)

// ErrorHandlerFunc defines the function to handle an error in the middleware.
type ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)

// MiddlewareOptions defines the middleware options structure.
type MiddlewareOptions struct {
	UnauthorizedHandlerFunc ErrorHandlerFunc
	ForbiddenHandlerFunc    ErrorHandlerFunc
}

// Middleware validates the JWT in the Authorization header and ensures that the associated subject has the necessary
// roles to access an endpoint based on the configured Roles.
func (s *service) Middleware(options MiddlewareOptions) func(http.Handler) http.Handler {
	// unauthorized defines a function to handle an unauthorized HTTP response.
	unauthorized := options.UnauthorizedHandlerFunc
	if unauthorized == nil {
		unauthorized = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
	}

	// forbidden defines a function to handle a forbidden HTTP response.
	forbidden := options.ForbiddenHandlerFunc
	if forbidden == nil {
		forbidden = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusForbidden)
		}
	}

	// validateRoles defines a function to validate the subject roles to access the endpoint.
	validateRoles := func(r *http.Request, requiredRoles []string, subjectRoles []string, subject string) error {
		// Check that the subject contains at least one of the required roles.
		containedRoles := make([]string, 0, len(requiredRoles))
		for _, requiredRole := range requiredRoles {
			if !slices.Contains(subjectRoles, requiredRole) {
				continue
			}

			containedRoles = append(containedRoles, requiredRole)
		}

		if len(containedRoles) == 0 {
			return ErrRolesInvalid
		}

		// Check that the subject contains the Admin role.
		if slices.Contains(subjectRoles, s.roles.AdminRole) {
			return nil
		}

		// Check that the subject is the one actually contained in the path wildcard.
		var requiredAuthorization string
		for _, wildcard := range s.roles.AuthzWildcards {
			pathValue := r.PathValue(wildcard)
			if len(pathValue) == 0 {
				continue
			}

			requiredAuthorization = pathValue
		}

		if len(requiredAuthorization) == 0 {
			// There is no authorization wildcard.
			return nil
		}
		if requiredAuthorization != subject {
			return ErrAuthorizationInvalid
		}

		// This means that the subject has the necessary roles to access the endpoint and has the authority of the data
		// they are trying to manipulate.
		return nil
	}

	return func(nextHandler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Get required roles.
			var requiredRoles []string
			bearerAuthScopes := ctx.Value(spec.BearerAuthScopes)
			if bearerAuthScopes != nil {
				requiredRoles = bearerAuthScopes.([]string)
			}

			// Check if there are any required roles.
			if len(requiredRoles) != 0 {
				// Get token from the Authorization header.
				authorizationValue := r.Header.Get(requestHeaderKeyAuthorization)
				token, ok := strings.CutPrefix(authorizationValue, bearerSchemaPrefix)
				if !ok {
					unauthorized(w, r, fmt.Errorf("%s: %w", descriptionFailedToFindBearerSchemaPrefix, ErrAuthorizationHeaderInvalid))
					return
				}

				// Get JWT claims.
				claims, err := s.authnService.ParseJWT(token)
				if err != nil {
					unauthorized(w, r, fmt.Errorf("%s: %w", descriptionFailedToParseJWT, ErrJWTInvalid))
					return
				}

				// Evaluate roles.
				subjectRoles := make([]string, len(claims.Roles))
				for i, role := range claims.Roles {
					subjectRoles[i] = string(role)
				}

				err = validateRoles(r, requiredRoles, subjectRoles, claims.Subject)
				if err != nil {
					forbidden(w, r, fmt.Errorf("%s: %w", descriptionFailedToValidateRoles, err))
					return
				}
			}

			// Serve next handler.
			nextHandler.ServeHTTP(w, r)
		})
	}
}
