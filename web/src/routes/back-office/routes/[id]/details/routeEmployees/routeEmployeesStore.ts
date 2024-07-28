import type {
	PaginatedRouteEmployees,
	RouteEmployeesFilters,
} from "$domain/routeEmployee";
import ecomapHttpClient from "$lib/clients/ecomap/http";
import { DEFAULT_PAGE_SIZE } from "$lib/constants/pagination";
import { BackOfficeRoutes } from "$lib/constants/routes";
import { createTableStore } from "$lib/stores/table";

/**
 * The search parameter names for each filter of the route employees table.
 */
const FILTERS_PARAMS_NAMES: Record<keyof RouteEmployeesFilters, string> = {
	pageIndex: "pageIndex",
	routeRole: "routeRole",
};

/**
 * The initial data of the route employees table.
 */
const initialData: PaginatedRouteEmployees = {
	employees: [],
	total: 0,
};

/**
 * The initial filters of the route employees table.
 */
export const initialFilters: RouteEmployeesFilters = {
	pageIndex: 0,
	routeRole: undefined,
};

/**
 * Maps URL search params to route employees filters.
 * @param searchParams URL search params.
 * @returns Route employees filters.
 */
function searchParamsToFilters(
	searchParams: URLSearchParams,
): RouteEmployeesFilters {
	let pageIndex = initialFilters.pageIndex;
	let routeRole = initialFilters.routeRole;

	const pageIndexParam = Number(
		searchParams.get(FILTERS_PARAMS_NAMES.pageIndex),
	);
	const routeRoleParam = searchParams.get(FILTERS_PARAMS_NAMES.routeRole);

	// Update page index when it's a valid number.
	if (!Number.isNaN(pageIndexParam)) {
		pageIndex = pageIndexParam;
	}

	// Update role when it's a valid role.
	switch (routeRoleParam) {
		case "driver":
		case "collector":
			routeRole = routeRoleParam;
	}

	return {
		pageIndex,
		routeRole,
	};
}

/**
 * Maps filters of the route employees table to URL search params.
 * @param filters route employees filters.
 * @returns URL search params.
 */
function filtersToSearchParams(
	filters: RouteEmployeesFilters,
): URLSearchParams {
	const { pageIndex, routeRole } = filters;

	const searchParams = new URLSearchParams(window.location.search);
	searchParams.set(FILTERS_PARAMS_NAMES.pageIndex, pageIndex.toString());

	if (routeRole) {
		searchParams.set(FILTERS_PARAMS_NAMES.routeRole, routeRole);
	} else {
		searchParams.delete(FILTERS_PARAMS_NAMES.routeRole);
	}

	return searchParams;
}

/**
 * Retrieves route employees to be displayed in the route employees table.
 * @param routeId Route ID.
 * @param filters Route employees filters.
 * @returns Route employees.
 */
async function getRouteEmployees(
	routeId: string,
	filters: RouteEmployeesFilters,
): Promise<PaginatedRouteEmployees> {
	const { pageIndex, routeRole } = filters;

	const res = await ecomapHttpClient.GET("/routes/{routeId}/employees", {
		params: {
			path: {
				routeId,
			},
			query: {
				offset: pageIndex * DEFAULT_PAGE_SIZE,
				limit: DEFAULT_PAGE_SIZE,
				sort: "createdAt",
				order: "desc",
				routeRole,
			},
		},
	});

	if (res.error) {
		return { total: 0, employees: [] };
	}

	return res.data;
}

/**
 * Creates a route employees store.
 * @param routeId Route ID.
 * @returns Route employees store.
 */
function createRouteEmployeesStore(routeId: string) {
	const routePathname = `${BackOfficeRoutes.ROUTES}/${routeId}`;

	return createTableStore(
		routePathname,
		initialData,
		filtersToSearchParams,
		searchParamsToFilters,
		filters => getRouteEmployees(routeId, filters),
	);
}

export default createRouteEmployeesStore;
