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
		var domainErrFieldInvalid *domain.ErrFieldInvalid

		switch {
		case errors.As(err, &domainErrFieldInvalid):
			badRequest(w, fmt.Sprintf("%s: %s", errFieldInvalid, domainErrFieldInvalid.FieldName))
		case errors.Is(err, domain.ErrUserAlreadyExists):
			conflict(w, errUserAlreadyExists)
		default:
			internalServerError(w)
		}

		return
	}

	user := userFromDomainUser(domainUser)
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
	// TODO: Implement this.
}

// GetUserByID handles the http request to get a user by ID.
func (h *handler) GetUserByID(w http.ResponseWriter, r *http.Request, userID spec.UserIdPathParam) {
	// TODO: Implement this.
}

// PatchUserByID handles the http request to modify a user by ID.
func (h *handler) PatchUserByID(w http.ResponseWriter, r *http.Request, userID spec.UserIdPathParam) {
	// TODO: Implement this.
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
