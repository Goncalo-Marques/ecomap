package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/goncalo-marques/ecomap/server/internal/domain"
)

const (
	constraintRoutesEmployeesEmployeeIDFkey = "routes_employees_employee_id_fkey"
	constraintRoutesEmployeesRouteIDFkey    = "routes_employees_route_id_fkey"
)

// ListRouteEmployees executes a query to return the route employees for the specified filter.
func (s *store) ListRouteEmployees(ctx context.Context, tx pgx.Tx, routeID uuid.UUID, filter domain.RouteEmployeesPaginatedFilter) (domain.PaginatedResponse[domain.RouteEmployee], error) {
	var filterFields []string
	var argsWhere []any

	// Append the optional fields to filter.
	filterFields = append(filterFields, "re.route_id::text")
	argsWhere = append(argsWhere, routeID)
	if filter.RouteRole != nil {
		filterFields = append(filterFields, "re.employee_role::text")
		argsWhere = append(argsWhere, routeEmployeeRoleFromDomain(*filter.RouteRole))
	}

	sqlWhere := listSQLWhere(filterFields, nil)

	// Get the total number of rows for the given filter.
	var total int
	row := tx.QueryRow(ctx, `
		SELECT count(re.route_id) 
		FROM routes_employees AS re
		INNER JOIN employees AS e ON re.employee_id = e.id
		LEFT JOIN road_network AS rn ON e.road_id = rn.id
		LEFT JOIN municipalities AS m ON e.municipality_id = m.id
	`+sqlWhere,
		argsWhere...,
	)

	err := row.Scan(&total)
	if err != nil {
		return domain.PaginatedResponse[domain.RouteEmployee]{}, fmt.Errorf("%s: %w", descriptionFailedScanRow, err)
	}

	// Append the field to sort, if provided.
	var domainSortField domain.RouteEmployeePaginatedSort
	if filter.Sort != nil {
		domainSortField = filter.Sort.Field()
	}

	sortField := "re.created_at"
	switch domainSortField {
	case domain.RouteEmployeePaginatedSortRouteRole:
		sortField = "re.employee_role"
	case domain.RouteEmployeePaginatedSortCreatedAt:
		sortField = "re.created_at"
	}

	// Get rows for the given filter.
	rows, err := tx.Query(ctx, `
		SELECT re.employee_role,
			e.id, e.username, e.first_name, e.last_name, e.role, e.date_of_birth, e.phone_number, ST_AsGeoJSON(e.geom)::jsonb, rn.osm_name, m.name, e.schedule_start, e.schedule_end, e.created_at, e.modified_at 
		FROM routes_employees AS re
		INNER JOIN employees AS e ON re.employee_id = e.id
		LEFT JOIN road_network AS rn ON e.road_id = rn.id
		LEFT JOIN municipalities AS m ON e.municipality_id = m.id
	`+sqlWhere+listSQLOrder(sortField, filter.Order)+listSQLLimitOffset(filter.Limit, filter.Offset),
		argsWhere...,
	)
	if err != nil {
		return domain.PaginatedResponse[domain.RouteEmployee]{}, fmt.Errorf("%s: %w", descriptionFailedQuery, err)
	}
	defer rows.Close()

	routeEmployees, err := getRouteEmployeesFromRows(rows)
	if err != nil {
		return domain.PaginatedResponse[domain.RouteEmployee]{}, fmt.Errorf("%s: %w", descriptionFailedScanRows, err)
	}

	return domain.PaginatedResponse[domain.RouteEmployee]{
		Total:   total,
		Results: routeEmployees,
	}, nil
}

// routeEmployeeRoleFromDomain returns a store route employee role based on the domain model.
func routeEmployeeRoleFromDomain(role domain.RouteEmployeeRole) string {
	switch role {
	case domain.RouteEmployeeRoleDriver:
		return "driver"
	case domain.RouteEmployeeRoleCollector:
		return "collector"
	default:
		return string(role)
	}
}

// routeEmployeeRoleToDomain returns a domain route employee role based on the store model.
func routeEmployeeRoleToDomain(role string) domain.RouteEmployeeRole {
	switch role {
	case "driver":
		return domain.RouteEmployeeRoleDriver
	case "collector":
		return domain.RouteEmployeeRoleCollector
	default:
		return domain.RouteEmployeeRole(role)
	}
}

// getRouteEmployeeFromRow returns the route employee by scanning the given row.
func getRouteEmployeeFromRow(row pgx.Row) (domain.RouteEmployee, error) {
	var routeEmployee domain.RouteEmployee
	var routeRole string
	var role string
	var geoJSONPoint domain.GeoJSONGeometryPoint
	var wayName *string
	var municipalityName *string

	err := row.Scan(
		&routeRole,
		&routeEmployee.ID,
		&routeEmployee.Username,
		&routeEmployee.FirstName,
		&routeEmployee.LastName,
		&role,
		&routeEmployee.DateOfBirth,
		&routeEmployee.PhoneNumber,
		&geoJSONPoint,
		&wayName,
		&municipalityName,
		&routeEmployee.ScheduleStart,
		&routeEmployee.ScheduleEnd,
		&routeEmployee.CreatedAt,
		&routeEmployee.ModifiedAt,
	)
	if err != nil {
		return domain.RouteEmployee{}, err
	}

	routeEmployee.Role = employeeRoleToDomain(role)
	routeEmployee.RouteRole = routeEmployeeRoleToDomain(routeRole)

	geoJSONProperties := make(domain.GeoJSONFeatureProperties)
	if wayName != nil {
		geoJSONProperties.SetWayName(*wayName)
	}
	if municipalityName != nil {
		geoJSONProperties.SetMunicipalityName(*municipalityName)
	}

	routeEmployee.GeoJSON = domain.GeoJSONFeature{
		Geometry:   geoJSONPoint,
		Properties: geoJSONProperties,
	}

	return routeEmployee, nil
}

// getRouteEmployeesFromRows returns the route employees by scanning the given rows.
func getRouteEmployeesFromRows(rows pgx.Rows) ([]domain.RouteEmployee, error) {
	var routeEmployees []domain.RouteEmployee
	for rows.Next() {
		routeEmployee, err := getRouteEmployeeFromRow(rows)
		if err != nil {
			return nil, err
		}

		routeEmployees = append(routeEmployees, routeEmployee)
	}

	return routeEmployees, nil
}
