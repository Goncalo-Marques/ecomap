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
	descriptionFailedCreateRoute     = "service: failed to create route"
	descriptionFailedListRoutes      = "service: failed to list routes"
	descriptionFailedGetRouteByID    = "service: failed to get route by id"
	descriptionFailedPatchRoute      = "service: failed to patch route"
	descriptionFailedDeleteRouteByID = "service: failed to delete route by id"
	descriptionFailedGetRouteRoads   = "service: failed to get route roads"
)

// CreateRoute creates a new route with the specified data.
func (s *service) CreateRoute(ctx context.Context, editableRoute domain.EditableRoute) (domain.Route, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "CreateRoute"),
		slog.String(logging.RouteName, string(editableRoute.Name)),
		slog.String(logging.RouteTruckID, editableRoute.TruckID.String()),
		slog.String(logging.RouteDepartureWarehouseID, editableRoute.DepartureWarehouseID.String()),
		slog.String(logging.RouteArrivalWarehouseID, editableRoute.ArrivalWarehouseID.String()),
	}

	if !editableRoute.Name.Valid() {
		return domain.Route{}, logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: domain.FieldName}, descriptionInvalidFieldValue, logAttrs...)
	}

	var route domain.Route

	err := s.readWriteTx(ctx, func(tx pgx.Tx) error {
		id, err := s.store.CreateRoute(ctx, tx, editableRoute)
		if err != nil {
			return err
		}

		route, err = s.store.GetRouteByID(ctx, tx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTruckNotFound),
			errors.Is(err, domain.ErrRouteDepartureWarehouseNotFound),
			errors.Is(err, domain.ErrRouteArrivalWarehouseNotFound):
			return domain.Route{}, logInfoAndWrapError(ctx, err, descriptionFailedGetRouteByID, logAttrs...)
		default:
			return domain.Route{}, logAndWrapError(ctx, err, descriptionFailedCreateRoute, logAttrs...)
		}
	}

	return route, nil
}

// ListRoutes returns the routes with the specified filter.
func (s *service) ListRoutes(ctx context.Context, filter domain.RoutesPaginatedFilter) (domain.PaginatedResponse[domain.Route], error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "ListRoutes"),
	}

	if filter.Sort != nil && !filter.Sort.Valid() {
		return domain.PaginatedResponse[domain.Route]{}, logInfoAndWrapError(ctx, &domain.ErrFilterValueInvalid{FilterName: domain.FieldFilterSort}, descriptionInvalidFilterValue, logAttrs...)
	}
	if !filter.Order.Valid() {
		return domain.PaginatedResponse[domain.Route]{}, logInfoAndWrapError(ctx, &domain.ErrFilterValueInvalid{FilterName: domain.FieldFilterOrder}, descriptionInvalidFilterValue, logAttrs...)
	}
	if !filter.Limit.Valid() {
		return domain.PaginatedResponse[domain.Route]{}, logInfoAndWrapError(ctx, &domain.ErrFilterValueInvalid{FilterName: domain.FieldFilterLimit}, descriptionInvalidFilterValue, logAttrs...)
	}
	if !filter.Offset.Valid() {
		return domain.PaginatedResponse[domain.Route]{}, logInfoAndWrapError(ctx, &domain.ErrFilterValueInvalid{FilterName: domain.FieldFilterOffset}, descriptionInvalidFilterValue, logAttrs...)
	}

	var paginatedRoutes domain.PaginatedResponse[domain.Route]
	var err error

	err = s.readOnlyTx(ctx, func(tx pgx.Tx) error {
		paginatedRoutes, err = s.store.ListRoutes(ctx, tx, filter)
		return err
	})
	if err != nil {
		return domain.PaginatedResponse[domain.Route]{}, logAndWrapError(ctx, err, descriptionFailedListRoutes, logAttrs...)
	}

	return paginatedRoutes, nil
}

// GetRouteByID returns the route with the specified identifier.
func (s *service) GetRouteByID(ctx context.Context, id uuid.UUID) (domain.Route, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "GetRouteByID"),
		slog.String(logging.RouteID, id.String()),
	}

	var route domain.Route
	var err error

	err = s.readOnlyTx(ctx, func(tx pgx.Tx) error {
		route, err = s.store.GetRouteByID(ctx, tx, id)
		return err
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRouteNotFound):
			return domain.Route{}, logInfoAndWrapError(ctx, err, descriptionFailedGetRouteByID, logAttrs...)
		default:
			return domain.Route{}, logAndWrapError(ctx, err, descriptionFailedGetRouteByID, logAttrs...)
		}
	}

	return route, nil
}

// PatchRoute modifies the route with the specified identifier.
func (s *service) PatchRoute(ctx context.Context, id uuid.UUID, editableRoute domain.EditableRoutePatch) (domain.Route, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "PatchRoute"),
		slog.String(logging.RouteID, id.String()),
	}

	if editableRoute.Name != nil && !editableRoute.Name.Valid() {
		return domain.Route{}, logInfoAndWrapError(ctx, &domain.ErrFieldValueInvalid{FieldName: domain.FieldName}, descriptionInvalidFieldValue, logAttrs...)
	}

	var route domain.Route
	var err error

	err = s.readWriteTx(ctx, func(tx pgx.Tx) error {
		if editableRoute.TruckID != nil {
			truck, err := s.store.GetTruckByID(ctx, tx, *editableRoute.TruckID)
			if err != nil {
				return err
			}

			routeEmployees, err := s.store.ListRouteEmployees(ctx, tx, id, domain.RouteEmployeesPaginatedFilter{})
			if err != nil {
				return err
			}

			if int(truck.PersonCapacity) < routeEmployees.Total {
				return domain.ErrRouteTruckPersonCapacityMinLimit
			}
		}

		err = s.store.PatchRoute(ctx, tx, id, editableRoute)
		if err != nil {
			return err
		}

		route, err = s.store.GetRouteByID(ctx, tx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRouteNotFound),
			errors.Is(err, domain.ErrTruckNotFound),
			errors.Is(err, domain.ErrRouteDepartureWarehouseNotFound),
			errors.Is(err, domain.ErrRouteArrivalWarehouseNotFound),
			errors.Is(err, domain.ErrRouteTruckPersonCapacityMinLimit):
			return domain.Route{}, logInfoAndWrapError(ctx, err, descriptionFailedPatchRoute, logAttrs...)
		default:
			return domain.Route{}, logAndWrapError(ctx, err, descriptionFailedPatchRoute, logAttrs...)
		}
	}

	return route, nil
}

// DeleteRouteByID deletes the route with the specified identifier.
func (s *service) DeleteRouteByID(ctx context.Context, id uuid.UUID) (domain.Route, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "DeleteRouteByID"),
		slog.String(logging.RouteID, id.String()),
	}

	var route domain.Route
	var err error

	err = s.readWriteTx(ctx, func(tx pgx.Tx) error {
		route, err = s.store.GetRouteByID(ctx, tx, id)
		if err != nil {
			return err
		}

		err = s.store.DeleteRouteByID(ctx, tx, id)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRouteNotFound),
			errors.Is(err, domain.ErrRouteAssociatedWithRouteContainer),
			errors.Is(err, domain.ErrRouteAssociatedWithRouteEmployee):
			return domain.Route{}, logInfoAndWrapError(ctx, err, descriptionFailedDeleteRouteByID, logAttrs...)
		default:
			return domain.Route{}, logAndWrapError(ctx, err, descriptionFailedDeleteRouteByID, logAttrs...)
		}
	}

	return route, nil
}

// GetRouteRoads returns the route roads using the TSP and A* algorithms. The route starts at the departure warehouse,
// passes through all the route containers, and before terminating at the arrival warehouse, passes through the nearest
// landfill to the arrival warehouse.
func (s *service) GetRouteRoads(ctx context.Context, id uuid.UUID) (domain.GeoJSON, error) {
	logAttrs := []any{
		slog.String(logging.ServiceMethod, "GetRouteRoads"),
		slog.String(logging.RouteID, id.String()),
	}

	var roadsGeometry []domain.GeoJSONGeometryLineString

	err := s.readOnlyTx(ctx, func(tx pgx.Tx) error {
		route, err := s.store.GetRouteByID(ctx, tx, id)
		if err != nil {
			return err
		}

		routeContainerRoads, err := s.store.GetContainerRoadsByRouteID(ctx, tx, route.ID)
		if err != nil {
			return err
		}

		// Early return when the route does not contain any containers associated.
		if len(routeContainerRoads) == 0 {
			return nil
		}

		arrivalWarehouse, err := s.store.GetWarehouseByID(ctx, tx, route.ArrivalWarehouse.ID)
		if err != nil {
			return err
		}

		// Get essential roads.
		var departureRoad *domain.Road
		tempDepartureRoad, err := s.store.GetRoadByWarehouseID(ctx, tx, route.DepartureWarehouse.ID)
		if err != nil {
			if !errors.Is(err, domain.ErrRoadNotFound) {
				return err
			}
		} else {
			departureRoad = &tempDepartureRoad
		}

		var arrivalRoad *domain.Road
		tempArrivalRoad, err := s.store.GetRoadByWarehouseID(ctx, tx, route.ArrivalWarehouse.ID)
		if err != nil {
			if !errors.Is(err, domain.ErrRoadNotFound) {
				return err
			}
		} else {
			arrivalRoad = &tempArrivalRoad
		}

		var arrivalWarehouseGeometry domain.GeoJSONGeometryPoint
		if feature, ok := arrivalWarehouse.GeoJSON.(domain.GeoJSONFeature); ok {
			if g, ok := feature.Geometry.(domain.GeoJSONGeometryPoint); ok {
				arrivalWarehouseGeometry = g
			}
		}

		var landfillID *uuid.UUID
		landfill, err := s.store.GetLandfillClosestGeometry(ctx, tx, arrivalWarehouseGeometry)
		if err != nil {
			if !errors.Is(err, domain.ErrLandfillNotFound) {
				return err
			}
		} else {
			landfillID = &landfill.ID
		}

		var landfillRoad *domain.Road
		if landfillID != nil {
			tempLandfillRoad, err := s.store.GetRoadByLandfillID(ctx, tx, *landfillID)
			if err != nil {
				if !errors.Is(err, domain.ErrRoadNotFound) {
					return err
				}
			} else {
				landfillRoad = &tempLandfillRoad
			}
		}

		// Compute the TSP for the route container vertices, starting at the departure warehouse and ending at the
		// closest landfill to the arrival warehouse.
		vertexIDs := make([]int, 0, len(routeContainerRoads)+2)

		for _, road := range routeContainerRoads {
			if road.Source == nil {
				continue
			}

			vertexIDs = append(vertexIDs, *road.Source)
		}
		if departureRoad != nil && departureRoad.Source != nil {
			vertexIDs = append(vertexIDs, *departureRoad.Source)
		}
		if landfillRoad != nil && landfillRoad.Source != nil {
			vertexIDs = append(vertexIDs, *landfillRoad.Source)
		}

		departureVertexID := vertexIDs[0]
		if departureRoad != nil && departureRoad.Source != nil {
			departureVertexID = *departureRoad.Source
		}

		landfillVertexID := vertexIDs[0]
		if landfillRoad != nil && landfillRoad.Source != nil {
			landfillVertexID = *landfillRoad.Source
		}

		seqVertexIDs, err := s.store.GetRoadVerticesTSP(ctx, tx, vertexIDs, departureVertexID, landfillVertexID, true)
		if err != nil {
			return err
		}

		arrivalVertexID := departureVertexID
		if arrivalRoad != nil && arrivalRoad.Source != nil {
			arrivalVertexID = *arrivalRoad.Source
		}

		// Replace the last vertex with the last actual point of the route.
		// This is a blind replacement because the last vertex is always the same as the first.
		if len(seqVertexIDs) != 0 {
			seqVertexIDs[len(seqVertexIDs)-1] = arrivalVertexID
		}

		// Compute the roads based on the sequential vertices.
		roadsGeometry, err = s.store.GetRoadsGeometryAStar(ctx, tx, seqVertexIDs, true)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRouteNotFound):
			return nil, logInfoAndWrapError(ctx, err, descriptionFailedGetRouteRoads, logAttrs...)
		default:
			return nil, logAndWrapError(ctx, err, descriptionFailedGetRouteRoads, logAttrs...)
		}
	}

	geoJSONFeature := make([]domain.GeoJSONFeature, len(roadsGeometry))
	for i, geometry := range roadsGeometry {
		geoJSONFeature[i] = domain.GeoJSONFeature{
			Geometry: geometry,
		}
	}

	return domain.GeoJSONFeatureCollection{
		Features: geoJSONFeature,
	}, nil
}
