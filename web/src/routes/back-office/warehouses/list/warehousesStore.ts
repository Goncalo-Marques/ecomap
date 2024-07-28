import type { PaginatedWarehouses, WarehousesFilters } from "$domain/warehouse";
import ecomapHttpClient from "$lib/clients/ecomap/http";
import { DEFAULT_PAGE_SIZE } from "$lib/constants/pagination";
import { BackOfficeRoutes } from "$lib/constants/routes";
import { createTableStore } from "$lib/stores/table";

/**
 * The search parameter names for each filter of the warehouses table.
 */
const FILTERS_PARAMS_NAMES: Record<keyof WarehousesFilters, string> = {
	pageIndex: "pageIndex",
	sort: "sort",
	order: "order",
	location: "location",
};

/**
 * The initial data of the warehouses table.
 */
const initialData: PaginatedWarehouses = {
	warehouses: [],
	total: 0,
};

/**
 * The initial filters of the warehouses table.
 */
export const initialFilters: WarehousesFilters = {
	pageIndex: 0,
	sort: "createdAt",
	order: "desc",
	location: "",
};

/**
 * Maps URL search params to warehouses filters.
 * @param searchParams URL search params.
 * @returns Warehouses filters.
 */
function searchParamsToFilters(
	searchParams: URLSearchParams,
): WarehousesFilters {
	let pageIndex = initialFilters.pageIndex;
	let sort = initialFilters.sort;
	let order = initialFilters.order;
	let location = initialFilters.location;

	const pageIndexParam = Number(
		searchParams.get(FILTERS_PARAMS_NAMES.pageIndex),
	);
	const sortParam = searchParams.get(FILTERS_PARAMS_NAMES.sort);
	const orderParam = searchParams.get(FILTERS_PARAMS_NAMES.order);
	const locationParam = searchParams.get(FILTERS_PARAMS_NAMES.location);

	// Update page index when it's a valid number.
	if (!Number.isNaN(pageIndexParam)) {
		pageIndex = pageIndexParam;
	}

	// Update sort when it's a valid warehouse sortable field.
	switch (sortParam) {
		case "createdAt":
		case "truckCapacity":
		case "wayName":
		case "municipalityName":
		case "modifiedAt":
			sort = sortParam;
			break;
	}

	// Update order when it's a valid direction.
	switch (orderParam) {
		case "asc":
		case "desc":
			order = orderParam;
			break;
	}

	// Update location when it's a non empty value.
	if (locationParam) {
		location = locationParam;
	}

	return {
		pageIndex,
		sort,
		order,
		location,
	};
}

/**
 * Maps filters of the warehouses table to URL search params.
 * @param filters Warehouses filters.
 * @returns URL search params.
 */
function filtersToSearchParams(filters: WarehousesFilters): URLSearchParams {
	const { pageIndex, sort, order, location } = filters;

	const searchParams = new URLSearchParams(window.location.search);
	searchParams.set(FILTERS_PARAMS_NAMES.pageIndex, pageIndex.toString());

	if (sort) {
		searchParams.set(FILTERS_PARAMS_NAMES.sort, sort);
	} else {
		searchParams.delete(FILTERS_PARAMS_NAMES.sort);
	}

	if (order) {
		searchParams.set(FILTERS_PARAMS_NAMES.order, order);
	} else {
		searchParams.delete(FILTERS_PARAMS_NAMES.order);
	}

	if (location) {
		searchParams.set(FILTERS_PARAMS_NAMES.location, location);
	} else {
		searchParams.delete(FILTERS_PARAMS_NAMES.location);
	}

	return searchParams;
}

/**
 * Retrieves warehouses to be displayed in the warehouses table.
 * @param filters Warehouses filters.
 * @returns Warehouses.
 */
async function getWarehouses(
	filters: WarehousesFilters,
): Promise<PaginatedWarehouses> {
	const { pageIndex, sort, order, location } = filters;

	const res = await ecomapHttpClient.GET("/warehouses", {
		params: {
			query: {
				offset: pageIndex * DEFAULT_PAGE_SIZE,
				limit: DEFAULT_PAGE_SIZE,
				sort,
				order,
				locationName: location || undefined,
			},
		},
	});

	if (res.error) {
		return { total: 0, warehouses: [] };
	}

	return res.data;
}

const warehousesStore = createTableStore(
	BackOfficeRoutes.WAREHOUSES,
	initialData,
	filtersToSearchParams,
	searchParamsToFilters,
	getWarehouses,
);

export default warehousesStore;
