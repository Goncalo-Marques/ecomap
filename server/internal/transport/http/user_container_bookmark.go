package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	spec "github.com/goncalo-marques/ecomap/server/api/ecomap"
	"github.com/goncalo-marques/ecomap/server/internal/domain"
	"github.com/goncalo-marques/ecomap/server/internal/logging"
)

const (
	errUserContainerBookmarkAlreadyExists = "user container bookmark already exists"
	errUserContainerBookmarkNotFound      = "user container bookmark does not exist"
)

// ListUserContainerBookmarks handles the http request to list user container bookmarks.
func (h *handler) ListUserContainerBookmarks(w http.ResponseWriter, r *http.Request, userID spec.UserIdPathParam, params spec.ListUserContainerBookmarksParams) {
	ctx := r.Context()

	domainUserContainerBookmarksFilter := listUserContainerBookmarksParamsToDomain(params)
	domainPaginatedContainers, err := h.service.ListUserContainerBookmarks(ctx, userID, domainUserContainerBookmarksFilter)
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

	containersPaginated, err := containersPaginatedFromDomain(domainPaginatedContainers)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMapResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	responseBody, err := json.Marshal(containersPaginated)
	if err != nil {
		logging.Logger.ErrorContext(ctx, descriptionFailedToMarshalResponseBody, logging.Error(err))
		internalServerError(w)
		return
	}

	writeResponseJSON(w, http.StatusOK, responseBody)
}

// CreateUserContainerBookmark handles the http request to create a user container bookmark.
func (h *handler) CreateUserContainerBookmark(w http.ResponseWriter, r *http.Request, userID spec.UserIdPathParam, containerID spec.ContainerIdPathParam) {
	ctx := r.Context()

	err := h.service.CreateUserContainerBookmark(ctx, userID, containerID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			notFound(w, errUserNotFound)
		case errors.Is(err, domain.ErrContainerNotFound):
			notFound(w, errContainerNotFound)
		case errors.Is(err, domain.ErrUserContainerBookmarkAlreadyExists):
			conflict(w, errUserContainerBookmarkAlreadyExists)
		default:
			internalServerError(w)
		}

		return
	}

	writeResponseJSON(w, http.StatusNoContent, nil)
}

// DeleteUserContainerBookmark handles the http request to delete a user container bookmark.
func (h *handler) DeleteUserContainerBookmark(w http.ResponseWriter, r *http.Request, userID spec.UserIdPathParam, containerID spec.ContainerIdPathParam) {
	ctx := r.Context()

	err := h.service.DeleteUserContainerBookmark(ctx, userID, containerID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserContainerBookmarkNotFound):
			conflict(w, errUserContainerBookmarkNotFound)
		default:
			internalServerError(w)
		}

		return
	}

	writeResponseJSON(w, http.StatusNoContent, nil)
}

// listUserContainerBookmarksParamsToDomain returns a domain user container bookmarks paginated filter based on the
// standardized list containers parameters.
func listUserContainerBookmarksParamsToDomain(params spec.ListUserContainerBookmarksParams) domain.UserContainerBookmarksPaginatedFilter {
	domainSort := domain.UserContainerBookmarkPaginatedSortCreatedAt
	if params.Sort != nil {
		switch *params.Sort {
		case spec.ListUserContainerBookmarksParamsSortContainerCategory:
			domainSort = domain.UserContainerBookmarkPaginatedSortContainerCategory
		case spec.ListUserContainerBookmarksParamsSortContainerWayName:
			domainSort = domain.UserContainerBookmarkPaginatedSortContainerWayName
		case spec.ListUserContainerBookmarksParamsSortContainerMunicipalityName:
			domainSort = domain.UserContainerBookmarkPaginatedSortContainerMunicipalityName
		case spec.ListUserContainerBookmarksParamsSortCreatedAt:
			domainSort = domain.UserContainerBookmarkPaginatedSortCreatedAt
		default:
			domainSort = domain.UserContainerBookmarkPaginatedSort(*params.Sort)
		}
	}

	var domainContainerCategory *domain.ContainerCategory
	if params.ContainerCategory != nil {
		category := containerCategoryToDomain(*params.ContainerCategory)
		domainContainerCategory = &category
	}

	return domain.UserContainerBookmarksPaginatedFilter{
		PaginatedRequest: paginatedRequestToDomain(
			domainSort,
			(*spec.OrderQueryParam)(params.Order),
			params.Limit,
			params.Offset,
		),
		ContainerCategory: domainContainerCategory,
		LocationName:      params.LocationName,
	}
}
