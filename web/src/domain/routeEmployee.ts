import type { components } from "$api/ecomap/http";

/**
 * Route employee.
 */
export type RouteEmployee = components["schemas"]["RouteEmployee"];

/**
 * Paginated route employees.
 */
export type PaginatedRouteEmployees =
	components["schemas"]["RouteEmployeesPaginated"];

/**
 * Filters of route employees.
 */
export interface RouteEmployeesFilters {
	pageIndex: number;
	routeRole?: RouteEmployee["routeRole"];
}

/**
 * Role of a route employee.
 */
export type RouteEmployeeRole = components["schemas"]["RouteEmployeeRole"];

/**
 * Represents a selected employee for a route.
 */
export interface SelectedRouteEmployee {
	id: string;
	routeRole: RouteEmployeeRole;
}

/**
 * Represents the selected employees that are added or deleted from a route.
 */
export interface SelectedRouteEmployees {
	added: SelectedRouteEmployee[];
	deleted: SelectedRouteEmployee[];
}
