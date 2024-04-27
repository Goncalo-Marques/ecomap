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
	descriptionFailedCreateContainer  = "service: failed to create container"
	descriptionFailedGetContainerByID = "service: failed to get container by id"
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
