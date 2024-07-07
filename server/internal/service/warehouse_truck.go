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
	descriptionFailedCreateWarehouseTruck = "service: failed to create warehouse truck association"
	descriptionFailedListWarehouseTrucks  = "service: failed to list warehouse truck associations"
	descriptionFailedDeleteWarehouseTruck = "service: failed to delete warehouse truck association"
)

// CreateWarehouseTruck creates a warehouse truck association.
func (s *service) CreateWarehouseTruck(ctx context.Context, warehouseID, truckID uuid.UUID) error {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "CreateWarehouseTruck"),
		slog.String(logging.WarehouseID, warehouseID.String()),
		slog.String(logging.TruckID, truckID.String()),
	}

	err := s.readWriteTx(ctx, func(tx pgx.Tx) error {
		warehouse, err := s.store.GetWarehouseByID(ctx, tx, warehouseID)
		if err != nil {
			return err
		}

		warehouseTrucks, err := s.store.ListWarehouseTrucks(ctx, tx, warehouseID, domain.WarehouseTrucksPaginatedFilter{})
		if err != nil {
			return err
		}

		if int(warehouse.TruckCapacity) <= warehouseTrucks.Total {
			return domain.ErrWarehouseTruckCapacityMaxLimit
		}

		return s.store.CreateWarehouseTruck(ctx, tx, warehouseID, truckID)
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrWarehouseTruckCapacityMaxLimit),
			errors.Is(err, domain.ErrWarehouseTruckAlreadyExists),
			errors.Is(err, domain.ErrWarehouseNotFound),
			errors.Is(err, domain.ErrTruckNotFound):
			return logInfoAndWrapError(ctx, err, descriptionFailedCreateWarehouseTruck, logAttrs...)
		default:
			return logAndWrapError(ctx, err, descriptionFailedCreateWarehouseTruck, logAttrs...)
		}
	}

	return nil
}

// ListWarehouseTrucks returns the warehouse trucks with the specified filter.
func (s *service) ListWarehouseTrucks(ctx context.Context, warehouseID uuid.UUID, filter domain.WarehouseTrucksPaginatedFilter) (domain.PaginatedResponse[domain.Truck], error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "ListWarehouseTrucks"),
	}

	if filter.Sort != nil && !filter.Sort.Valid() {
		return domain.PaginatedResponse[domain.Truck]{}, logInfoAndWrapError(ctx, &domain.ErrFilterValueInvalid{FilterName: domain.FieldFilterSort}, descriptionInvalidFilterValue, logAttrs...)
	}
	if !filter.Order.Valid() {
		return domain.PaginatedResponse[domain.Truck]{}, logInfoAndWrapError(ctx, &domain.ErrFilterValueInvalid{FilterName: domain.FieldFilterOrder}, descriptionInvalidFilterValue, logAttrs...)
	}
	if !filter.Limit.Valid() {
		return domain.PaginatedResponse[domain.Truck]{}, logInfoAndWrapError(ctx, &domain.ErrFilterValueInvalid{FilterName: domain.FieldFilterLimit}, descriptionInvalidFilterValue, logAttrs...)
	}
	if !filter.Offset.Valid() {
		return domain.PaginatedResponse[domain.Truck]{}, logInfoAndWrapError(ctx, &domain.ErrFilterValueInvalid{FilterName: domain.FieldFilterOffset}, descriptionInvalidFilterValue, logAttrs...)
	}

	var paginatedTrucks domain.PaginatedResponse[domain.Truck]
	var err error

	err = s.readOnlyTx(ctx, func(tx pgx.Tx) error {
		paginatedTrucks, err = s.store.ListWarehouseTrucks(ctx, tx, warehouseID, filter)
		return err
	})
	if err != nil {
		return domain.PaginatedResponse[domain.Truck]{}, logAndWrapError(ctx, err, descriptionFailedListWarehouseTrucks, logAttrs...)
	}

	return paginatedTrucks, nil
}

// DeleteWarehouseTruck deletes the warehouse truck association.
func (s *service) DeleteWarehouseTruck(ctx context.Context, warehouseID, truckID uuid.UUID) error {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "DeleteWarehouseTruck"),
		slog.String(logging.WarehouseID, warehouseID.String()),
		slog.String(logging.TruckID, truckID.String()),
	}

	err := s.readWriteTx(ctx, func(tx pgx.Tx) error {
		routes, err := s.store.ListRoutes(ctx, tx, domain.RoutesPaginatedFilter{
			TruckID:              &truckID,
			DepartureWarehouseID: &warehouseID,
		})
		if err != nil {
			return err
		}

		if routes.Total > 0 {
			return domain.ErrWarehouseTruckAssociatedWithRouteDeparture
		}

		routes, err = s.store.ListRoutes(ctx, tx, domain.RoutesPaginatedFilter{
			TruckID:            &truckID,
			ArrivalWarehouseID: &warehouseID,
		})
		if err != nil {
			return err
		}

		if routes.Total > 0 {
			return domain.ErrWarehouseTruckAssociatedWithRouteArrival
		}

		return s.store.DeleteWarehouseTruck(ctx, tx, warehouseID, truckID)
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrWarehouseTruckAssociatedWithRouteDeparture),
			errors.Is(err, domain.ErrWarehouseTruckAssociatedWithRouteArrival),
			errors.Is(err, domain.ErrWarehouseTruckNotFound):
			return logInfoAndWrapError(ctx, err, descriptionFailedDeleteWarehouseTruck, logAttrs...)
		default:
			return logAndWrapError(ctx, err, descriptionFailedDeleteWarehouseTruck, logAttrs...)
		}
	}

	return nil
}
