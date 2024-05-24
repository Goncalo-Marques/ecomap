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
	descriptionFailedCreateContainer     = "service: failed to create container"
	descriptionFailedListContainers      = "service: failed to list containers"
	descriptionFailedGetContainerByID    = "service: failed to get container by id"
	descriptionFailedPatchContainer      = "service: failed to patch container"
	descriptionFailedDeleteContainerByID = "service: failed to delete container by id"
)

// CreateContainer creates a new container with the specified data.
func (s *service) CreateContainer(ctx context.Context, editableContainer domain.EditableContainer) (domain.Container, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "CreateContainer"),
		slog.String(logging.ContainerCategory, string(editableContainer.Category)),
	}

	if !editableContainer.Category.Valid() {
		return domain.Container{}, logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: domain.FieldCategory}, descriptionInvalidFieldValue, logAttrs...)
	}

	var geometry domain.GeoJSONGeometryPoint
	if feature, ok := editableContainer.GeoJSON.(domain.GeoJSONFeature); ok {
		if g, ok := feature.Geometry.(domain.GeoJSONGeometryPoint); ok {
			geometry = g
		}
	}

	var container domain.Container

	err := s.readWriteTx(ctx, func(tx pgx.Tx) error {
		var roadID *int
		road, err := s.store.GetRoadByGeometry(ctx, tx, geometry)
		if err != nil {
			if !errors.Is(err, domain.ErrRoadNotFound) {
				return err
			}
		} else {
			roadID = &road.ID
		}

		var municipalityID *int
		municipality, err := s.store.GetMunicipalityByGeometry(ctx, tx, geometry)
		if err != nil {
			if !errors.Is(err, domain.ErrMunicipalityNotFound) {
				return err
			}
		} else {
			municipalityID = &municipality.ID
		}

		id, err := s.store.CreateContainer(ctx, tx, editableContainer, roadID, municipalityID)
		if err != nil {
			return err
		}

		container, err = s.store.GetContainerByID(ctx, tx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return domain.Container{}, logAndWrapError(ctx, err, descriptionFailedCreateContainer, logAttrs...)
	}

	return container, nil
}

// ListContainers returns the containers with the specified filter.
func (s *service) ListContainers(ctx context.Context, filter domain.ContainersPaginatedFilter) (domain.PaginatedResponse[domain.Container], error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "ListContainers"),
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
		paginatedContainers, err = s.store.ListContainers(ctx, tx, filter)
		return err
	})
	if err != nil {
		return domain.PaginatedResponse[domain.Container]{}, logAndWrapError(ctx, err, descriptionFailedListContainers, logAttrs...)
	}

	return paginatedContainers, nil
}

// GetContainerByID returns the container with the specified identifier.
func (s *service) GetContainerByID(ctx context.Context, id uuid.UUID) (domain.Container, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "GetContainerByID"),
		slog.String(logging.ContainerID, id.String()),
	}

	var container domain.Container
	var err error

	err = s.readOnlyTx(ctx, func(tx pgx.Tx) error {
		container, err = s.store.GetContainerByID(ctx, tx, id)
		return err
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrContainerNotFound):
			return domain.Container{}, logInfoAndWrapError(ctx, err, descriptionFailedGetContainerByID, logAttrs...)
		default:
			return domain.Container{}, logAndWrapError(ctx, err, descriptionFailedGetContainerByID, logAttrs...)
		}
	}

	return container, nil
}

// PatchContainer modifies the container with the specified identifier.
func (s *service) PatchContainer(ctx context.Context, id uuid.UUID, editableContainer domain.EditableContainerPatch) (domain.Container, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "PatchContainer"),
		slog.String(logging.ContainerID, id.String()),
	}

	if editableContainer.Category != nil && !editableContainer.Category.Valid() {
		return domain.Container{}, logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: domain.FieldCategory}, descriptionInvalidFieldValue, logAttrs...)
	}

	var geometry domain.GeoJSONGeometryPoint
	if editableContainer.GeoJSON != nil {
		if feature, ok := editableContainer.GeoJSON.(domain.GeoJSONFeature); ok {
			if g, ok := feature.Geometry.(domain.GeoJSONGeometryPoint); ok {
				geometry = g
			}
		}
	}

	var container domain.Container
	var err error

	err = s.readWriteTx(ctx, func(tx pgx.Tx) error {
		var roadID *int
		var municipalityID *int

		if editableContainer.GeoJSON != nil {
			road, err := s.store.GetRoadByGeometry(ctx, tx, geometry)
			if err != nil {
				if !errors.Is(err, domain.ErrRoadNotFound) {
					return err
				}
			} else {
				roadID = &road.ID
			}

			municipality, err := s.store.GetMunicipalityByGeometry(ctx, tx, geometry)
			if err != nil {
				if !errors.Is(err, domain.ErrMunicipalityNotFound) {
					return err
				}
			} else {
				municipalityID = &municipality.ID
			}
		}

		err = s.store.PatchContainer(ctx, tx, id, editableContainer, roadID, municipalityID)
		if err != nil {
			return err
		}

		container, err = s.store.GetContainerByID(ctx, tx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrContainerNotFound):
			return domain.Container{}, logInfoAndWrapError(ctx, err, descriptionFailedPatchContainer, logAttrs...)
		default:
			return domain.Container{}, logAndWrapError(ctx, err, descriptionFailedPatchContainer, logAttrs...)
		}
	}

	return container, nil
}

// DeleteContainerByID deletes the container with the specified identifier.
func (s *service) DeleteContainerByID(ctx context.Context, id uuid.UUID) (domain.Container, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "DeleteContainerByID"),
		slog.String(logging.ContainerID, id.String()),
	}

	var container domain.Container
	var err error

	err = s.readWriteTx(ctx, func(tx pgx.Tx) error {
		container, err = s.store.GetContainerByID(ctx, tx, id)
		if err != nil {
			return err
		}

		err = s.store.DeleteContainerByID(ctx, tx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrContainerNotFound),
			errors.Is(err, domain.ErrContainerAssociatedWithUserContainerBookmark),
			errors.Is(err, domain.ErrContainerAssociatedWithRouteContainer):
			return domain.Container{}, logInfoAndWrapError(ctx, err, descriptionFailedDeleteContainerByID, logAttrs...)
		default:
			return domain.Container{}, logAndWrapError(ctx, err, descriptionFailedDeleteContainerByID, logAttrs...)
		}
	}

	return container, nil
}
