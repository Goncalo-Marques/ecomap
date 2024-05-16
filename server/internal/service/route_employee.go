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
	descriptionFailedCreateRouteEmployee = "service: failed to create route employee association"
	descriptionFailedListRouteEmployees  = "service: failed to list route employee associations"
	descriptionFailedDeleteRouteEmployee = "service: failed to delete route employee association"
)

// CreateRouteEmployee creates a route employee association.
func (s *service) CreateRouteEmployee(ctx context.Context, routeID, employeeID uuid.UUID, editableRouteEmployee domain.EditableRouteEmployee) error {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "CreateRouteEmployee"),
		slog.String(logging.RouteID, routeID.String()),
		slog.String(logging.EmployeeID, employeeID.String()),
		slog.String(logging.RouteEmployeeEmployeeRole, string(editableRouteEmployee.RouteRole)),
	}

	if !editableRouteEmployee.RouteRole.Valid() {
		return logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: domain.FieldRouteRole}, descriptionInvalidFieldValue, logAttrs...)
	}

	err := s.readWriteTx(ctx, func(tx pgx.Tx) error {
		route, err := s.store.GetRouteByID(ctx, tx, routeID)
		if err != nil {
			return err
		}

		truck, err := s.store.GetTruckByID(ctx, tx, route.TruckID)
		if err != nil {
			return err
		}

		routeEmployees, err := s.store.ListRouteEmployees(ctx, tx, route.ID, domain.RouteEmployeesPaginatedFilter{})
		if err != nil {
			return err
		}

		if int(truck.PersonCapacity) <= routeEmployees.Total {
			return domain.ErrRouteTruckPersonCapacityMaxLimit
		}

		return s.store.CreateRouteEmployee(ctx, tx, routeID, employeeID, editableRouteEmployee)
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRouteTruckPersonCapacityMaxLimit),
			errors.Is(err, domain.ErrRouteEmployeeAlreadyExists),
			errors.Is(err, domain.ErrRouteNotFound),
			errors.Is(err, domain.ErrEmployeeNotFound):
			return logInfoAndWrapError(ctx, err, descriptionFailedCreateRouteEmployee, logAttrs...)
		default:
			return logAndWrapError(ctx, err, descriptionFailedCreateRouteEmployee, logAttrs...)
		}
	}

	return nil
}

// ListRouteEmployees returns the route employees with the specified filter.
func (s *service) ListRouteEmployees(ctx context.Context, routeID uuid.UUID, filter domain.RouteEmployeesPaginatedFilter) (domain.PaginatedResponse[domain.RouteEmployee], error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "ListRouteEmployees"),
	}

	if filter.Sort != nil && !filter.Sort.Valid() {
		return domain.PaginatedResponse[domain.RouteEmployee]{}, logInfoAndWrapError(ctx, &domain.ErrFilterValueInvalid{FilterName: domain.FieldFilterSort}, descriptionInvalidFilterValue, logAttrs...)
	}
	if !filter.Order.Valid() {
		return domain.PaginatedResponse[domain.RouteEmployee]{}, logInfoAndWrapError(ctx, &domain.ErrFilterValueInvalid{FilterName: domain.FieldFilterOrder}, descriptionInvalidFilterValue, logAttrs...)
	}
	if !filter.Limit.Valid() {
		return domain.PaginatedResponse[domain.RouteEmployee]{}, logInfoAndWrapError(ctx, &domain.ErrFilterValueInvalid{FilterName: domain.FieldFilterLimit}, descriptionInvalidFilterValue, logAttrs...)
	}
	if !filter.Offset.Valid() {
		return domain.PaginatedResponse[domain.RouteEmployee]{}, logInfoAndWrapError(ctx, &domain.ErrFilterValueInvalid{FilterName: domain.FieldFilterOffset}, descriptionInvalidFilterValue, logAttrs...)
	}

	var paginatedEmployees domain.PaginatedResponse[domain.RouteEmployee]
	var err error

	err = s.readOnlyTx(ctx, func(tx pgx.Tx) error {
		paginatedEmployees, err = s.store.ListRouteEmployees(ctx, tx, routeID, filter)
		return err
	})
	if err != nil {
		return domain.PaginatedResponse[domain.RouteEmployee]{}, logAndWrapError(ctx, err, descriptionFailedListRouteEmployees, logAttrs...)
	}

	return paginatedEmployees, nil
}

// DeleteRouteEmployee deletes the route employee association.
func (s *service) DeleteRouteEmployee(ctx context.Context, routeID, employeeID uuid.UUID) error {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "DeleteRouteEmployee"),
		slog.String(logging.RouteID, routeID.String()),
		slog.String(logging.EmployeeID, employeeID.String()),
	}

	err := s.readWriteTx(ctx, func(tx pgx.Tx) error {
		return s.store.DeleteRouteEmployee(ctx, tx, routeID, employeeID)
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRouteEmployeeNotFound):
			return logInfoAndWrapError(ctx, err, descriptionFailedDeleteRouteEmployee, logAttrs...)
		default:
			return logAndWrapError(ctx, err, descriptionFailedDeleteRouteEmployee, logAttrs...)
		}
	}

	return nil
}
