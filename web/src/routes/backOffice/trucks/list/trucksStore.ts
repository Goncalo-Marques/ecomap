import type { PaginatedTrucks, TrucksFilters } from "../../../../domain/truck";
import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
import { DEFAULT_PAGE_SIZE } from "../../../../lib/constants/pagination";
import { createTableStore } from "../../../../lib/stores/table";
import { BackOfficeRoutes } from "../../../constants/routes";

/**
 * The search parameter names for each filter of the trucks table.
 */
const FILTERS_PARAMS_NAMES: Record<keyof TrucksFilters, string> = {
	pageIndex: "pageIndex",
	sort: "sort",
	order: "order",
	licensePlate: "licensePlate",
};

/**
 * The initial data of the trucks table.
 */
const initialData: PaginatedTrucks = {
	trucks: [],
	total: 0,
};

/**
 * The initial filters of the trucks table.
 */
export const initialFilters: TrucksFilters = {
	pageIndex: 0,
	sort: "createdAt",
	order: "desc",
	licensePlate: "",
};

/**
 * Maps URL search params to trucks filters.
 * @param searchParams URL search params.
 * @returns Trucks filters.
 */
function searchParamsToFilters(searchParams: URLSearchParams): TrucksFilters {
	let pageIndex = initialFilters.pageIndex;
	let sort = initialFilters.sort;
	let order = initialFilters.order;
	let licensePlate = initialFilters.licensePlate;

	const pageIndexParam = Number(
		searchParams.get(FILTERS_PARAMS_NAMES.pageIndex),
	);
	const sortParam = searchParams.get(FILTERS_PARAMS_NAMES.sort);
	const orderParam = searchParams.get(FILTERS_PARAMS_NAMES.order);
	const licensePlateParam = searchParams.get(FILTERS_PARAMS_NAMES.licensePlate);

	// Update page index when it's a valid number.
	if (!Number.isNaN(pageIndexParam)) {
		pageIndex = pageIndexParam;
	}

	// Update sort when it's a valid truck sortable field.
	switch (sortParam) {
		case "licensePlate":
		case "createdAt":
		case "make":
		case "model":
		case "personCapacity":
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

	// Update license plate when it's a non empty value.
	if (licensePlateParam) {
		licensePlate = licensePlateParam;
	}

	return {
		pageIndex,
		sort,
		order,
		licensePlate,
	};
}

/**
 * Maps filters of the trucks table to URL search params.
 * @param filters Trucks filters.
 * @returns URL search params.
 */
function filtersToSearchParams(filters: TrucksFilters): URLSearchParams {
	const { pageIndex, sort, order, licensePlate } = filters;

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

	if (licensePlate) {
		searchParams.set(FILTERS_PARAMS_NAMES.licensePlate, licensePlate);
	} else {
		searchParams.delete(FILTERS_PARAMS_NAMES.licensePlate);
	}

	return searchParams;
}

/**
 * Retrieves trucks to be displayed in the trucks table.
 * @param filters Trucks filters.
 * @returns Trucks.
 */
async function getTrucks(filters: TrucksFilters): Promise<PaginatedTrucks> {
	const { pageIndex, sort, order, licensePlate } = filters;

	const res = await ecomapHttpClient.GET("/trucks", {
		params: {
			query: {
				offset: pageIndex * DEFAULT_PAGE_SIZE,
				limit: DEFAULT_PAGE_SIZE,
				sort,
				order,
				licensePlate: licensePlate || undefined,
			},
		},
	});

	if (res.error) {
		return { total: 0, trucks: [] };
	}

	return res.data;
}

const trucksStore = createTableStore(
	BackOfficeRoutes.TRUCKS,
	initialData,
	filtersToSearchParams,
	searchParamsToFilters,
	getTrucks,
);

export default trucksStore;
