package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	spec "github.com/goncalo-marques/ecomap/server/api/ecomap"
	"github.com/goncalo-marques/ecomap/server/internal/domain"
	"github.com/goncalo-marques/ecomap/server/internal/logging"
)

const (
	errUserNotFound      = "user not found"
	errUserAlreadyExists = "username already exists"
)

// CreateUser handles the http request to create a user.
func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	var userPost spec.UserPost
	err = json.Unmarshal(requestBody, &userPost)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	domainEditableUser := userPostToDomainEditableUserWithPassword(userPost)
	domainUser, err := h.service.CreateUser(ctx, domainEditableUser)
	if err != nil {
		var domainErrFieldValueInvalid *domain.ErrFieldValueInvalid

		switch {
		case errors.As(err, &domainErrFieldValueInvalid):
			badRequest(w, fmt.Sprintf("%s: %s", errFieldValueInvalid, domainErrFieldValueInvalid.FieldName))
		case errors.Is(err, domain.ErrUserAlreadyExists):
			conflict(w, errUserAlreadyExists)
		default:
			internalServerError(w)
		}

		return
	}

	user := userFromDomain(domainUser)
	responseBody, err := json.Marshal(user)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusCreated, responseBody)
}

// ListUsers handles the http request to list users.
func (h *handler) ListUsers(w http.ResponseWriter, r *http.Request, params spec.ListUsersParams) {
	ctx := r.Context()

	var domainSort domain.Sort[domain.UserSort]
	if params.Sort != nil {
		domainSort = domain.UserSort(*params.Sort)
	}

	domainUsersFilter := domain.UsersFilter{
		PaginatedRequest: domain.PaginatedRequest[domain.UserSort]{
			Limit:  limitToDomain(params.Limit),
			Offset: offsetToDomain(params.Offset),
			Sort:   domainSort,
			Order:  orderToDomain((*spec.OrderQueryParam)(params.Order)),
		},
		Username:  (*domain.Username)(params.Username),
		FirstName: (*domain.Name)(params.FirstName),
		LastName:  (*domain.Name)(params.LastName),
	}

	domainPaginatedUsers, err := h.service.ListUsers(ctx, domainUsersFilter)
	if err != nil {
		var domainErrFilterValueInvalid *domain.ErrFilterValueInvalid

		switch {
		case errors.As(err, &domainErrFilterValueInvalid):
			badRequest(w, fmt.Sprintf("%s: %s", errFilterValueInvalid, domainErrFilterValueInvalid.FilterName))
		default:
			internalServerError(w)
		}

		return
	}

	users := usersFromDomain(domainPaginatedUsers.Results)
	specUsersPaginated := spec.UsersPaginated{
		Total: domainPaginatedUsers.Total,
		Users: users,
	}

	responseBody, err := json.Marshal(specUsersPaginated)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}

// GetUserByID handles the http request to get a user by ID.
func (h *handler) GetUserByID(w http.ResponseWriter, r *http.Request, userID spec.UserIdPathParam) {
	ctx := r.Context()

	domainUser, err := h.service.GetUserByID(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			notFound(w, errUserNotFound)
		default:
			internalServerError(w)
		}

		return
	}

	user := userFromDomain(domainUser)
	responseBody, err := json.Marshal(user)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}

// PatchUserByID handles the http request to modify a user by ID.
func (h *handler) PatchUserByID(w http.ResponseWriter, r *http.Request, userID spec.UserIdPathParam) {
	ctx := r.Context()

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	var userPatch spec.UserPatch
	err = json.Unmarshal(requestBody, &userPatch)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	domainEditableUser := userPatchToDomainEditableUserPatch(userPatch)
	domainUser, err := h.service.PatchUser(ctx, userID, domainEditableUser)
	if err != nil {
		var domainErrFieldValueInvalid *domain.ErrFieldValueInvalid

		switch {
		case errors.As(err, &domainErrFieldValueInvalid):
			badRequest(w, fmt.Sprintf("%s: %s", errFieldValueInvalid, domainErrFieldValueInvalid.FieldName))
		case errors.Is(err, domain.ErrUserNotFound):
			notFound(w, errUserNotFound)
		case errors.Is(err, domain.ErrUserAlreadyExists):
			conflict(w, errUserAlreadyExists)
		default:
			internalServerError(w)
		}

		return
	}

	user := userFromDomain(domainUser)
	responseBody, err := json.Marshal(user)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}

// DeleteUserByID handles the http request to delete a user by ID.
func (h *handler) DeleteUserByID(w http.ResponseWriter, r *http.Request, userID spec.UserIdPathParam) {
	// TODO: Implement this.
}

// UpdateUserPassword handles the http request to update a user password.
func (h *handler) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement this.
}

// ResetUserPassword handles the http request to reset a user password.
func (h *handler) ResetUserPassword(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement this.
}

// SignInUser handles the http request to sign in a user.
func (h *handler) SignInUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	var signIn spec.SignIn
	err = json.Unmarshal(requestBody, &signIn)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	token, err := h.service.SignInUser(ctx, domain.Username(signIn.Username), signIn.Password)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrCredentialsIncorrect):
			unauthorized(w, errCredentialsIncorrect)
		default:
			internalServerError(w)
		}

		return
	}

	jwt := jwtFromJWTToken(token)
	responseBody, err := json.Marshal(jwt)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}
