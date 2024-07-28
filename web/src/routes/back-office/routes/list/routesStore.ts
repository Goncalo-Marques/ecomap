import type { PaginatedRoutes, RoutesFilters } from "$domain/route";
import ecomapHttpClient from "$lib/clients/ecomap/http";
import { DEFAULT_PAGE_SIZE } from "$lib/constants/pagination";
import { BackOfficeRoutes } from "$lib/constants/routes";
import { createTableStore } from "$lib/stores/table";

/**
 * The search parameter names for each filter of the routes table.
 */
const FILTERS_PARAMS_NAMES: Record<keyof RoutesFilters, string> = {
	pageIndex: "pageIndex",
	route: "route",
};

/**
 * The initial data of the routes table.
 */
const initialData: PaginatedRoutes = {
	routes: [],
	total: 0,
};

/**
 * The initial filters of the routes table.
 */
export const initialFilters: RoutesFilters = {
	pageIndex: 0,
	route: "",
};

/**
 * Maps URL search params to routes filters.
 * @param searchParams URL search params.
 * @returns Routes filters.
 */
function searchParamsToFilters(searchParams: URLSearchParams): RoutesFilters {
	let pageIndex = initialFilters.pageIndex;
	let route = initialFilters.route;

	const pageIndexParam = Number(
		searchParams.get(FILTERS_PARAMS_NAMES.pageIndex),
	);
	const routeParam = searchParams.get(FILTERS_PARAMS_NAMES.route);

	// Update page index when it's a valid number.
	if (!Number.isNaN(pageIndexParam)) {
		pageIndex = pageIndexParam;
	}

	// Update route when it's a non empty value.
	if (routeParam) {
		route = routeParam;
	}

	return {
		pageIndex,
		route,
	};
}

/**
 * Maps filters of the routes table to URL search params.
 * @param filters Routes filters.
 * @returns URL search params.
 */
function filtersToSearchParams(filters: RoutesFilters): URLSearchParams {
	const { pageIndex, route } = filters;

	const searchParams = new URLSearchParams(window.location.search);
	searchParams.set(FILTERS_PARAMS_NAMES.pageIndex, pageIndex.toString());

	if (route) {
		searchParams.set(FILTERS_PARAMS_NAMES.route, route);
	} else {
		searchParams.delete(FILTERS_PARAMS_NAMES.route);
	}

	return searchParams;
}

/**
 * Retrieves routes to be displayed in the routes table.
 * @param filters Routes filters.
 * @returns Routes.
 */
async function getRoutes(filters: RoutesFilters): Promise<PaginatedRoutes> {
	const { pageIndex, route } = filters;

	const res = await ecomapHttpClient.GET("/routes", {
		params: {
			query: {
				offset: pageIndex * DEFAULT_PAGE_SIZE,
				limit: DEFAULT_PAGE_SIZE,
				sort: "createdAt",
				order: "desc",
				name: route || undefined,
			},
		},
	});

	if (res.error) {
		return { total: 0, routes: [] };
	}

	return res.data;
}

const routesStore = createTableStore(
	BackOfficeRoutes.ROUTES,
	initialData,
	filtersToSearchParams,
	searchParamsToFilters,
	getRoutes,
);

export default routesStore;
