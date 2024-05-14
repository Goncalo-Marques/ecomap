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
	descriptionFailedCreateRouteContainer = "service: failed to create route container association"
	descriptionFailedListRouteContainers  = "service: failed to list route container associations"
	descriptionFailedDeleteRouteContainer = "service: failed to delete route container association"
)

// CreateRouteContainer creates a route container association.
func (s *service) CreateRouteContainer(ctx context.Context, routeID, containerID uuid.UUID) error {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "CreateRouteContainer"),
		slog.String(logging.RouteID, routeID.String()),
		slog.String(logging.ContainerID, containerID.String()),
	}

	err := s.readWriteTx(ctx, func(tx pgx.Tx) error {
		return s.store.CreateRouteContainer(ctx, tx, routeID, containerID)
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRouteContainerAlreadyExists),
			errors.Is(err, domain.ErrRouteNotFound),
			errors.Is(err, domain.ErrContainerNotFound):
			return logInfoAndWrapError(ctx, err, descriptionFailedCreateRouteContainer, logAttrs...)
		default:
			return logAndWrapError(ctx, err, descriptionFailedCreateRouteContainer, logAttrs...)
		}
	}

	return nil
}

// ListRouteContainers returns the route containers with the specified filter.
func (s *service) ListRouteContainers(ctx context.Context, routeID uuid.UUID, filter domain.RouteContainersPaginatedFilter) (domain.PaginatedResponse[domain.Container], error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "ListRouteContainers"),
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
		paginatedContainers, err = s.store.ListRouteContainers(ctx, tx, routeID, filter)
		return err
	})
	if err != nil {
		return domain.PaginatedResponse[domain.Container]{}, logAndWrapError(ctx, err, descriptionFailedListRouteContainers, logAttrs...)
	}

	return paginatedContainers, nil
}

// DeleteRouteContainer deletes the route container association.
func (s *service) DeleteRouteContainer(ctx context.Context, routeID, containerID uuid.UUID) error {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "DeleteRouteContainer"),
		slog.String(logging.RouteID, routeID.String()),
		slog.String(logging.ContainerID, containerID.String()),
	}

	err := s.readWriteTx(ctx, func(tx pgx.Tx) error {
		return s.store.DeleteRouteContainer(ctx, tx, routeID, containerID)
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRouteContainerNotFound):
			return logInfoAndWrapError(ctx, err, descriptionFailedDeleteRouteContainer, logAttrs...)
		default:
			return logAndWrapError(ctx, err, descriptionFailedDeleteRouteContainer, logAttrs...)
		}
	}

	return nil
}
