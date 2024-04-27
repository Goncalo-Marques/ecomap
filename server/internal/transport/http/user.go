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

	domainEditableUser := userPostToDomain(userPost)
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

	domainUsersFilter := listUsersParamsToDomain(params)
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

	usersPaginated := usersPaginatedFromDomain(domainPaginatedUsers)
	responseBody, err := json.Marshal(usersPaginated)
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

	domainEditableUser := userPatchToDomain(userPatch)
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
	ctx := r.Context()

	domainUser, err := h.service.DeleteUserByID(ctx, userID)
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

// UpdateUserPassword handles the http request to update a user password.
func (h *handler) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	var passwordChange spec.PasswordChange
	err = json.Unmarshal(requestBody, &passwordChange)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	err = h.service.UpdateUserPassword(ctx, domain.Username(passwordChange.Username), domain.Password(passwordChange.OldPassword), domain.Password(passwordChange.NewPassword))
	if err != nil {
		var domainErrFieldValueInvalid *domain.ErrFieldValueInvalid

		switch {
		case errors.As(err, &domainErrFieldValueInvalid):
			badRequest(w, fmt.Sprintf("%s: %s", errFieldValueInvalid, domainErrFieldValueInvalid.FieldName))
		case errors.Is(err, domain.ErrCredentialsIncorrect):
			unauthorized(w, errCredentialsIncorrect)
		default:
			internalServerError(w)
		}

		return
	}

	writeResponseJSON(w, http.StatusNoContent, nil)
}

// ResetUserPassword handles the http request to reset a user password.
func (h *handler) ResetUserPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	var passwordReset spec.PasswordReset
	err = json.Unmarshal(requestBody, &passwordReset)
	if err != nil {
		badRequest(w, errRequestBodyInvalid)
		return
	}

	err = h.service.ResetUserPassword(ctx, domain.Username(passwordReset.Username), domain.Password(passwordReset.NewPassword))
	if err != nil {
		var domainErrFieldValueInvalid *domain.ErrFieldValueInvalid

		switch {
		case errors.As(err, &domainErrFieldValueInvalid):
			badRequest(w, fmt.Sprintf("%s: %s", errFieldValueInvalid, domainErrFieldValueInvalid.FieldName))
		case errors.Is(err, domain.ErrUserNotFound):
			notFound(w, errUserNotFound)
		default:
			internalServerError(w)
		}

		return
	}

	writeResponseJSON(w, http.StatusNoContent, nil)
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

	token, err := h.service.SignInUser(ctx, domain.Username(signIn.Username), domain.Password(signIn.Password))
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

// userPostToDomain returns a domain editable user with password based on the standardized user post.
func userPostToDomain(userPost spec.UserPost) domain.EditableUserWithPassword {
	return domain.EditableUserWithPassword{
		EditableUser: domain.EditableUser{
			Username:  domain.Username(userPost.Username),
			FirstName: domain.Name(userPost.FirstName),
			LastName:  domain.Name(userPost.LastName),
		},
		Password: domain.Password(userPost.Password),
	}
}

// userPatchToDomain returns a domain patchable user based on the standardized user patch.
func userPatchToDomain(userPatch spec.UserPatch) domain.EditableUserPatch {
	return domain.EditableUserPatch{
		Username:  (*domain.Username)(userPatch.Username),
		FirstName: (*domain.Name)(userPatch.FirstName),
		LastName:  (*domain.Name)(userPatch.LastName),
	}
}

// listUsersParamsToDomain returns a domain users paginated filter based on the standardized list users parameters.
func listUsersParamsToDomain(params spec.ListUsersParams) domain.UsersPaginatedFilter {
	domainSort := domain.UserPaginatedSortCreatedAt
	if params.Sort != nil {
		switch *params.Sort {
		case spec.ListUsersParamsSortUsername:
			domainSort = domain.UserPaginatedSortUsername
		case spec.ListUsersParamsSortFirstName:
			domainSort = domain.UserPaginatedSortFirstName
		case spec.ListUsersParamsSortLastName:
			domainSort = domain.UserPaginatedSortLastName
		case spec.ListUsersParamsSortCreatedAt:
			domainSort = domain.UserPaginatedSortCreatedAt
		case spec.ListUsersParamsSortModifiedAt:
			domainSort = domain.UserPaginatedSortModifiedAt
		default:
			domainSort = domain.UserPaginatedSort(*params.Sort)
		}
	}

	return domain.UsersPaginatedFilter{
		PaginatedRequest: paginatedRequestToDomain(
			(*spec.LogicalOperatorQueryParam)(params.LogicalOperator),
			domainSort,
			(*spec.OrderQueryParam)(params.Order),
			params.Limit,
			params.Offset,
		),
		Username:  (*domain.Username)(params.Username),
		FirstName: (*domain.Name)(params.FirstName),
		LastName:  (*domain.Name)(params.LastName),
	}
}

// userFromDomain returns a standardized user based on the domain model.
func userFromDomain(user domain.User) spec.User {
	return spec.User{
		Id:         user.ID,
		Username:   string(user.Username),
		FirstName:  string(user.FirstName),
		LastName:   string(user.LastName),
		CreatedAt:  user.CreatedAt,
		ModifiedAt: user.ModifiedAt,
	}
}

// usersFromDomain returns standardized users based on the domain model.
func usersFromDomain(users []domain.User) []spec.User {
	specUsers := make([]spec.User, len(users))
	for i, user := range users {
		specUsers[i] = userFromDomain(user)
	}

	return specUsers
}

// usersPaginatedFromDomain returns a standardized users paginated response based on the domain model.
func usersPaginatedFromDomain(paginatedResponse domain.PaginatedResponse[domain.User]) spec.UsersPaginated {
	return spec.UsersPaginated{
		Total: paginatedResponse.Total,
		Users: usersFromDomain(paginatedResponse.Results),
	}
}
