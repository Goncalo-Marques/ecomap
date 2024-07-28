import type { components } from "$api/ecomap/http";

/**
 * Route.
 */
export type Route = components["schemas"]["Route"];

/**
 * Paginated routes.
 */
export type PaginatedRoutes = components["schemas"]["RoutesPaginated"];

/**
 * Filters of routes.
 */
export interface RoutesFilters {
	pageIndex: number;
	route: string;
}
