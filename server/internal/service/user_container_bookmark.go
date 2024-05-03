package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/goncalo-marques/ecomap/server/internal/domain"
	"github.com/goncalo-marques/ecomap/server/internal/logging"
)

const (
	descriptionFailedCreateUserContainerBookmark = "service: failed to create user container bookmark"
	descriptionFailedListUserContainerBookmarks  = "service: failed to list user container bookmarks"
	descriptionFailedDeleteUserContainerBookmark = "service: failed to delete user container bookmark"
)

// CreateUserContainerBookmark creates a new user container bookmark.
func (s *service) CreateUserContainerBookmark(ctx context.Context, userID, containerID uuid.UUID) error {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "CreateUserContainerBookmark"),
		slog.String(logging.UserID, userID.String()),
		slog.String(logging.ContainerID, containerID.String()),
	}

	err := s.readWriteTx(ctx, func(tx pgx.Tx) error {
		return s.store.CreateUserContainerBookmark(ctx, tx, userID, containerID)
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserContainerBookmarkAlreadyExists),
			errors.Is(err, domain.ErrUserNotFound),
			errors.Is(err, domain.ErrContainerNotFound):
			return logInfoAndWrapError(ctx, err, descriptionFailedCreateUserContainerBookmark, logAttrs...)
		default:
			return logAndWrapError(ctx, err, descriptionFailedCreateUserContainerBookmark, logAttrs...)
		}
	}

	return nil
}

// ListUserContainerBookmarks returns the user container bookmarks with the specified filter.
func (s *service) ListUserContainerBookmarks(ctx context.Context, userID uuid.UUID, filter domain.UserContainerBookmarksPaginatedFilter) (domain.PaginatedResponse[domain.Container], error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "ListUserContainerBookmarks"),
	}

	if filter.Sort != nil && !filter.Sort.Valid() {
		return domain.PaginatedResponse[domain.Container]{}, logInfoAndWrapError(ctx, &domain.ErrFilterValueInvalid{FilterName: domain.FieldFilterSort}, descriptionInvalidFilterValue, logAttrs...)
	}
	if !filter.Order.Valid() {
		return domain.PaginatedResponse[domain.Container]{}, logInfoAndWrapError(ctx, &domain.ErrFilterValueInvalid{FilterName: domain.FieldFilterOrder}, descriptionInvalidFilterValue, logAttrs...)
	}
	if !filter.Limit.Valid() {
		return domain.PaginatedResponse[domain.Container]{}, logInfoAndWrapError(ctx, &domain.ErrFilterValueInvalid{FilterName: domain.FieldFilterLimit}, descriptionInvalidFilterValue, logAttrs...)
	}
	if !filter.Offset.Valid() {
		return domain.PaginatedResponse[domain.Container]{}, logInfoAndWrapError(ctx, &domain.ErrFilterValueInvalid{FilterName: domain.FieldFilterOffset}, descriptionInvalidFilterValue, logAttrs...)
	}

	var paginatedContainers domain.PaginatedResponse[domain.Container]
	var err error

	err = s.readOnlyTx(ctx, func(tx pgx.Tx) error {
		paginatedContainers, err = s.store.ListUserContainerBookmarks(ctx, tx, userID, filter)
		return err
	})
	if err != nil {
		return domain.PaginatedResponse[domain.Container]{}, logAndWrapError(ctx, err, descriptionFailedListUserContainerBookmarks, logAttrs...)
	}

	return paginatedContainers, nil
}

// DeleteUserContainerBookmark deletes the user container bookmark.
func (s *service) DeleteUserContainerBookmark(ctx context.Context, userID, containerID uuid.UUID) error {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "DeleteUserContainerBookmark"),
		slog.String(logging.UserID, userID.String()),
		slog.String(logging.ContainerID, containerID.String()),
	}

	err := s.readWriteTx(ctx, func(tx pgx.Tx) error {
		return s.store.DeleteUserContainerBookmark(ctx, tx, userID, containerID)
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserContainerBookmarkNotFound):
			return logInfoAndWrapError(ctx, err, descriptionFailedDeleteUserContainerBookmark, logAttrs...)
		default:
			return logAndWrapError(ctx, err, descriptionFailedDeleteUserContainerBookmark, logAttrs...)
		}
	}

	return nil
}
