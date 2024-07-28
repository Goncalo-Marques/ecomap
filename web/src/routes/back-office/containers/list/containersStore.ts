import type { ContainersFilters, PaginatedContainers } from "$domain/container";
import ecomapHttpClient from "$lib/clients/ecomap/http";
import { DEFAULT_PAGE_SIZE } from "$lib/constants/pagination";
import { BackOfficeRoutes } from "$lib/constants/routes";
import { createTableStore } from "$lib/stores/table";

/**
 * The search parameter names for each filter of the containers table.
 */
const FILTERS_PARAMS_NAMES: Record<keyof ContainersFilters, string> = {
	pageIndex: "pageIndex",
	location: "location",
	category: "category",
};

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
	let location = initialFilters.location;
	let category = initialFilters.category;

	const pageIndexParam = Number(
		searchParams.get(FILTERS_PARAMS_NAMES.pageIndex),
	);
	const locationParam = searchParams.get(FILTERS_PARAMS_NAMES.location);
	const categoryParam = searchParams.get(FILTERS_PARAMS_NAMES.category);

	// Update page index when it's a valid number.
	if (!Number.isNaN(pageIndexParam)) {
		pageIndex = pageIndexParam;
	}

	// Update location when it's a non empty value.
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
	};
}

/**
 * Maps filters of the containers table to URL search params.
 * @param filters Containers filters.
 * @returns URL search params.
 */
function filtersToSearchParams(filters: ContainersFilters): URLSearchParams {
	const { pageIndex, location, category } = filters;

	const searchParams = new URLSearchParams(window.location.search);
	searchParams.set(FILTERS_PARAMS_NAMES.pageIndex, pageIndex.toString());

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
	const { pageIndex, location, category } = filters;

	const res = await ecomapHttpClient.GET("/containers", {
		params: {
			query: {
				offset: pageIndex * DEFAULT_PAGE_SIZE,
				limit: DEFAULT_PAGE_SIZE,
				sort: "createdAt",
				order: "desc",
				category,
				locationName: location || undefined,
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
