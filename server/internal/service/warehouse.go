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
	descriptionFailedCreateWarehouse     = "service: failed to create warehouse"
	descriptionFailedListWarehouses      = "service: failed to list warehouses"
	descriptionFailedGetWarehouseByID    = "service: failed to get warehouse by id"
	descriptionFailedPatchWarehouse      = "service: failed to patch warehouse"
	descriptionFailedDeleteWarehouseByID = "service: failed to delete warehouse by id"
)

// CreateWarehouse creates a new warehouse with the specified data.
func (s *service) CreateWarehouse(ctx context.Context, editableWarehouse domain.EditableWarehouse) (domain.Warehouse, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "CreateWarehouse"),
		slog.Int(logging.WarehouseTruckCapacity, int(editableWarehouse.TruckCapacity)),
	}

	if !editableWarehouse.TruckCapacity.Valid() {
		return domain.Warehouse{}, logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: domain.FieldTruckCapacity}, descriptionInvalidFieldValue, logAttrs...)
	}

	var geometry domain.GeoJSONGeometryPoint
	if feature, ok := editableWarehouse.GeoJSON.(domain.GeoJSONFeature); ok {
		if g, ok := feature.Geometry.(domain.GeoJSONGeometryPoint); ok {
			geometry = g
		}
	}

	var warehouse domain.Warehouse

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

		id, err := s.store.CreateWarehouse(ctx, tx, editableWarehouse, roadID, municipalityID)
		if err != nil {
			return err
		}

		warehouse, err = s.store.GetWarehouseByID(ctx, tx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return domain.Warehouse{}, logAndWrapError(ctx, err, descriptionFailedCreateWarehouse, logAttrs...)
	}

	return warehouse, nil
}

// ListWarehouses returns the warehouses with the specified filter.
func (s *service) ListWarehouses(ctx context.Context, filter domain.WarehousesPaginatedFilter) (domain.PaginatedResponse[domain.Warehouse], error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "ListWarehouses"),
	}

	if filter.Sort != nil && !filter.Sort.Valid() {
		return domain.PaginatedResponse[domain.Warehouse]{}, logInfoAndWrapError(ctx, &domain.ErrFilterValueInvalid{FilterName: domain.FieldFilterSort}, descriptionInvalidFilterValue, logAttrs...)
	}
	if !filter.Order.Valid() {
		return domain.PaginatedResponse[domain.Warehouse]{}, logInfoAndWrapError(ctx, &domain.ErrFilterValueInvalid{FilterName: domain.FieldFilterOrder}, descriptionInvalidFilterValue, logAttrs...)
	}
	if !filter.Limit.Valid() {
		return domain.PaginatedResponse[domain.Warehouse]{}, logInfoAndWrapError(ctx, &domain.ErrFilterValueInvalid{FilterName: domain.FieldFilterLimit}, descriptionInvalidFilterValue, logAttrs...)
	}
	if !filter.Offset.Valid() {
		return domain.PaginatedResponse[domain.Warehouse]{}, logInfoAndWrapError(ctx, &domain.ErrFilterValueInvalid{FilterName: domain.FieldFilterOffset}, descriptionInvalidFilterValue, logAttrs...)
	}

	var paginatedWarehouses domain.PaginatedResponse[domain.Warehouse]
	var err error

	err = s.readOnlyTx(ctx, func(tx pgx.Tx) error {
		paginatedWarehouses, err = s.store.ListWarehouses(ctx, tx, filter)
		return err
	})
	if err != nil {
		return domain.PaginatedResponse[domain.Warehouse]{}, logAndWrapError(ctx, err, descriptionFailedListWarehouses, logAttrs...)
	}

	return paginatedWarehouses, nil
}

// GetWarehouseByID returns the warehouse with the specified identifier.
func (s *service) GetWarehouseByID(ctx context.Context, id uuid.UUID) (domain.Warehouse, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "GetWarehouseByID"),
		slog.String(logging.WarehouseID, id.String()),
	}

	var warehouse domain.Warehouse
	var err error

	err = s.readOnlyTx(ctx, func(tx pgx.Tx) error {
		warehouse, err = s.store.GetWarehouseByID(ctx, tx, id)
		return err
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrWarehouseNotFound):
			return domain.Warehouse{}, logInfoAndWrapError(ctx, err, descriptionFailedGetWarehouseByID, logAttrs...)
		default:
			return domain.Warehouse{}, logAndWrapError(ctx, err, descriptionFailedGetWarehouseByID, logAttrs...)
		}
	}

	return warehouse, nil
}

// PatchWarehouse modifies the warehouse with the specified identifier.
func (s *service) PatchWarehouse(ctx context.Context, id uuid.UUID, editableWarehouse domain.EditableWarehousePatch) (domain.Warehouse, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "PatchWarehouse"),
		slog.String(logging.WarehouseID, id.String()),
	}

	if editableWarehouse.TruckCapacity != nil && !editableWarehouse.TruckCapacity.Valid() {
		return domain.Warehouse{}, logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: domain.FieldTruckCapacity}, descriptionInvalidFieldValue, logAttrs...)
	}

	var geometry domain.GeoJSONGeometryPoint
	if editableWarehouse.GeoJSON != nil {
		if feature, ok := editableWarehouse.GeoJSON.(domain.GeoJSONFeature); ok {
			if g, ok := feature.Geometry.(domain.GeoJSONGeometryPoint); ok {
				geometry = g
			}
		}
	}

	var warehouse domain.Warehouse

	err := s.readWriteTx(ctx, func(tx pgx.Tx) error {
		if editableWarehouse.TruckCapacity != nil {
			trucks, err := s.store.ListWarehouseTrucks(ctx, tx, id)
			if err != nil {
				return err
			}

			if int(*editableWarehouse.TruckCapacity) < len(trucks) {
				return domain.ErrWarehouseTruckCapacityMinLimit
			}
		}

		var roadID *int
		var municipalityID *int

		if editableWarehouse.GeoJSON != nil {
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

		err := s.store.PatchWarehouse(ctx, tx, id, editableWarehouse, roadID, municipalityID)
		if err != nil {
			return err
		}

		warehouse, err = s.store.GetWarehouseByID(ctx, tx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrWarehouseTruckCapacityMinLimit),
			errors.Is(err, domain.ErrWarehouseNotFound):
			return domain.Warehouse{}, logInfoAndWrapError(ctx, err, descriptionFailedPatchWarehouse, logAttrs...)
		default:
			return domain.Warehouse{}, logAndWrapError(ctx, err, descriptionFailedPatchWarehouse, logAttrs...)
		}
	}

	return warehouse, nil
}

// DeleteWarehouseByID deletes the warehouse with the specified identifier.
func (s *service) DeleteWarehouseByID(ctx context.Context, id uuid.UUID) (domain.Warehouse, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "DeleteWarehouseByID"),
		slog.String(logging.WarehouseID, id.String()),
	}

	var warehouse domain.Warehouse
	var err error

	err = s.readWriteTx(ctx, func(tx pgx.Tx) error {
		warehouse, err = s.store.GetWarehouseByID(ctx, tx, id)
		if err != nil {
			return err
		}

		err = s.store.DeleteWarehouseByID(ctx, tx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrWarehouseNotFound),
			errors.Is(err, domain.ErrWarehouseAssociatedWithRouteDeparture),
			errors.Is(err, domain.ErrWarehouseAssociatedWithRouteArrival):
			return domain.Warehouse{}, logInfoAndWrapError(ctx, err, descriptionFailedDeleteWarehouseByID, logAttrs...)
		default:
			return domain.Warehouse{}, logAndWrapError(ctx, err, descriptionFailedDeleteWarehouseByID, logAttrs...)
		}
	}

	return warehouse, nil
}
