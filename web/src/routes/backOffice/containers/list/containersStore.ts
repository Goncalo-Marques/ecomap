import type {
	ContainersFilters,
	PaginatedContainers,
} from "../../../../domain/container";
import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
import { DEFAULT_PAGE_SIZE } from "../../../../lib/constants/pagination";
import { createTableStore } from "../../../../lib/stores/table";
import { BackOfficeRoutes } from "../../../constants/routes";

/**
 * The search parameter names for each filter of the containers table.
 */
const FILTERS_PARAMS_NAMES = {
	pageIndex: "page-index",
	sort: "sort",
	order: "order",
	location: "location",
	category: "category",
} as const;

/**
 * The initial data of the containers table.
 */
const initialData: PaginatedContainers = {
	containers: [],
	total: 0,
};

/**
 * The initial filters of the containers table.
 */
export const initialFilters: ContainersFilters = {
	pageIndex: 0,
	sort: "category",
	order: "asc",
	location: "",
	category: undefined,
};

/**
 * Maps URL search params to containers filters.
 * @param searchParams URL search params.
 * @returns Containers filters.
 */
function searchParamsToFilters(
	searchParams: URLSearchParams,
): ContainersFilters {
	let pageIndex = initialFilters.pageIndex;
	let sortingField = initialFilters.sort;
	let sortingDirection = initialFilters.order;
	let location = initialFilters.location;
	let category = initialFilters.category;

	const pageIndexParam = Number(
		searchParams.get(FILTERS_PARAMS_NAMES.pageIndex),
	);
	const sortParam = searchParams.get(FILTERS_PARAMS_NAMES.sort);
	const orderParam = searchParams.get(FILTERS_PARAMS_NAMES.order);
	const locationParam = searchParams.get(FILTERS_PARAMS_NAMES.location);
	const categoryParam = searchParams.get(FILTERS_PARAMS_NAMES.category);

	// Update page index when it's is a valid number.
	if (!Number.isNaN(pageIndexParam)) {
		pageIndex = pageIndexParam;
	}

	// Update sorting field when it's a valid sort for containers.
	switch (sortParam) {
		case "category":
		case "createdAt":
		case "modifiedAt":
			sortingField = sortParam;
	}

	// Update sorting direction when it's a valid direction for containers.
	switch (orderParam) {
		case "asc":
		case "desc":
			sortingDirection = orderParam;
	}

	if (locationParam) {
		location = locationParam;
	}

	switch (categoryParam) {
		case "general":
		case "paper":
		case "plastic":
		case "metal":
		case "glass":
		case "organic":
		case "hazardous":
			category = categoryParam;
	}

	return {
		pageIndex,
		location,
		category,
		sort: sortingField,
		order: sortingDirection,
	};
}

/**
 * Maps filters of the containers table to URL search params.
 * @param filters Containers filters.
 * @returns URL search params.
 */
function filtersToSearchParams(filters: ContainersFilters): URLSearchParams {
	const { pageIndex, sort, order, location, category } = filters;

	const searchParams = new URLSearchParams(window.location.search);

	searchParams.set(FILTERS_PARAMS_NAMES.pageIndex, pageIndex.toString());
	searchParams.set(FILTERS_PARAMS_NAMES.sort, sort);

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

	if (category) {
		searchParams.set(FILTERS_PARAMS_NAMES.category, category);
	} else {
		searchParams.delete(FILTERS_PARAMS_NAMES.category);
	}

	return searchParams;
}

/**
 * Retrieves containers to be displayed in the containers table.
 * @param filters Containers filters.
 * @returns Containers.
 */
async function getContainers(
	filters: ContainersFilters,
): Promise<PaginatedContainers> {
	const { pageIndex, sort, order, location, category } = filters;

	const res = await ecomapHttpClient.GET("/containers", {
		params: {
			query: {
				offset: pageIndex * DEFAULT_PAGE_SIZE,
				limit: DEFAULT_PAGE_SIZE,
				sort,
				order,
				category,
				wayName: location,
				municipalityName: location,
				logicalOperator: "or",
			},
		},
	});

	if (res.error) {
		return { total: 0, containers: [] };
	}

	return res.data;
}

const containersStore = createTableStore(
	BackOfficeRoutes.CONTAINERS,
	initialData,
	filtersToSearchParams,
	searchParamsToFilters,
	getContainers,
);

export default containersStore;
