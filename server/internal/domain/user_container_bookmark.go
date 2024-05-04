package domain

import "errors"

// User container bookmark errors.
var (
	ErrUserContainerBookmarkAlreadyExists = errors.New("user container bookmark already exists") // Returned when a user container bookmark already exists.
	ErrUserContainerBookmarkNotFound      = errors.New("user container bookmark not found")      // Returned when a user container bookmark is not found.
)

// UserContainerBookmarkPaginatedSort defines the field of the user container bookmark to sort.
type UserContainerBookmarkPaginatedSort string

const (
	UserContainerBookmarkPaginatedSortContainerCategory UserContainerBookmarkPaginatedSort = "containerCategory"
	UserContainerBookmarkPaginatedSortWayName           UserContainerBookmarkPaginatedSort = "wayName"
	UserContainerBookmarkPaginatedSortMunicipalityName  UserContainerBookmarkPaginatedSort = "municipalityName"
	UserContainerBookmarkPaginatedSortCreatedAt         UserContainerBookmarkPaginatedSort = "createdAt"
)

// Field returns the name of the field to sort by.
func (s UserContainerBookmarkPaginatedSort) Field() UserContainerBookmarkPaginatedSort {
	return s
}

// Valid returns true if the field is valid, false otherwise.
func (s UserContainerBookmarkPaginatedSort) Valid() bool {
	switch s {
	case UserContainerBookmarkPaginatedSortContainerCategory,
		UserContainerBookmarkPaginatedSortWayName,
		UserContainerBookmarkPaginatedSortMunicipalityName,
		UserContainerBookmarkPaginatedSortCreatedAt:
		return true
	default:
		return false
	}
}

// UserContainerBookmarksPaginatedFilter defines the user container bookmarks filter structure.
type UserContainerBookmarksPaginatedFilter struct {
	PaginatedRequest[UserContainerBookmarkPaginatedSort]
	ContainerCategory *ContainerCategory
	LocationName      *string
}
