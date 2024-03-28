package authz

import "net/http"

// ErrorHandlerFunc defines the function to handle an error in the middleware.
type ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)

// MiddlewareOptions defines the middleware options structure.
type MiddlewareOptions struct {
	UnauthorizedHandlerFunc        ErrorHandlerFunc
	ForbiddenHandlerFunc           ErrorHandlerFunc
	InternalServerErrorHandlerFunc ErrorHandlerFunc
}

// Middleware validates the JWT in the Authorization header and ensures that the associated subject has the necessary
// roles to access an endpoint based on the configured Roles.
func (s *service) Middleware(options MiddlewareOptions) func(http.Handler) http.Handler {
	unauthorized := options.UnauthorizedHandlerFunc
	if unauthorized == nil {
		unauthorized = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}
	}

	forbidden := options.ForbiddenHandlerFunc
	if forbidden == nil {
		forbidden = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusForbidden)
		}
	}

	internalServerError := options.InternalServerErrorHandlerFunc
	if internalServerError == nil {
		internalServerError = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	return func(nextHandler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: use option errors
			// TODO: Authorization: Bearer <token>
			// TODO: validate jwt
			// TODO: fail if role map do not contain pattern
			// TODO: don't validate roles for empty
			// TODO: guarantee contain necessary roles
			// TODO: admin early return in validation
			// TODO: guarantee non-admin contains the same id as the one defined

			// Serve next handler.
			nextHandler.ServeHTTP(w, r)
		})
	}
}
