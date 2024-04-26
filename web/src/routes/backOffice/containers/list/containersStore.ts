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
	pageIndex: "pageIndex",
	sort: "sort",
	order: "order",
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

	const pageIndexParam = Number(
		searchParams.get(FILTERS_PARAMS_NAMES.pageIndex),
	);
	const sortParam = searchParams.get(FILTERS_PARAMS_NAMES.sort);
	const orderParam = searchParams.get(FILTERS_PARAMS_NAMES.order);

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
			break;
		default:
			break;
	}

	// Update sorting direction when it's a valid direction for containers.
	switch (orderParam) {
		case "asc":
		case "desc":
			sortingDirection = orderParam;
			break;
		default:
			break;
	}

	return {
		pageIndex,
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
	const { pageIndex, sort, order } = filters;

	const searchParams = new URLSearchParams(location.search);

	searchParams.set(FILTERS_PARAMS_NAMES.pageIndex, pageIndex.toString());
	searchParams.set(FILTERS_PARAMS_NAMES.sort, sort);

	if (order) {
		searchParams.set(FILTERS_PARAMS_NAMES.order, order);
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
	const { pageIndex, sort, order } = filters;

	const res = await ecomapHttpClient.GET("/containers", {
		params: {
			query: {
				offset: pageIndex * DEFAULT_PAGE_SIZE,
				limit: DEFAULT_PAGE_SIZE,
				sort,
				order,
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
