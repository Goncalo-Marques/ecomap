package http

import (
	"net/http"

	spec "github.com/goncalo-marques/ecomap/server/api/ecomap"
)

// ListUserContainerBookmarks handles the http request to list user container bookmarks.
func (h *handler) ListUserContainerBookmarks(w http.ResponseWriter, r *http.Request, userID spec.UserIdPathParam, params spec.ListUserContainerBookmarksParams) {
	w.WriteHeader(http.StatusNotFound)
}

// CreateUserContainerBookmark handles the http request to create a user container bookmark.
func (h *handler) CreateUserContainerBookmark(w http.ResponseWriter, r *http.Request, userID spec.UserIdPathParam, containerID spec.ContainerIdPathParam) {
	w.WriteHeader(http.StatusNotFound)
}

// DeleteUserContainerBookmark handles the http request to delete a user container bookmark.
func (h *handler) DeleteUserContainerBookmark(w http.ResponseWriter, r *http.Request, userID spec.UserIdPathParam, containerID spec.ContainerIdPathParam) {
	w.WriteHeader(http.StatusNotFound)
}
