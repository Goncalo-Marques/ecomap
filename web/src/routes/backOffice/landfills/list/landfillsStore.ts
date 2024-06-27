import type {
	PaginatedLandfills,
	LandfillsFilters,
} from "../../../../domain/landfill";
import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
import { DEFAULT_PAGE_SIZE } from "../../../../lib/constants/pagination";
import { createTableStore } from "../../../../lib/stores/table";
import { BackOfficeRoutes } from "../../../constants/routes";

/**
 * The search parameter names for each filter of the landfills table.
 */
const FILTERS_PARAMS_NAMES: Record<keyof LandfillsFilters, string> = {
	pageIndex: "pageIndex",
	location: "location",
};

/**
 * The initial data of the landfills table.
 */
const initialData: PaginatedLandfills = {
	landfills: [],
	total: 0,
};

/**
 * The initial filters of the landfills table.
 */
export const initialFilters: LandfillsFilters = {
	pageIndex: 0,
	location: "",
};

/**
 * Maps URL search params to landfills filters.
 * @param searchParams URL search params.
 * @returns Landfills filters.
 */
function searchParamsToFilters(
	searchParams: URLSearchParams,
): LandfillsFilters {
	let pageIndex = initialFilters.pageIndex;
	let location = initialFilters.location;

	const pageIndexParam = Number(
		searchParams.get(FILTERS_PARAMS_NAMES.pageIndex),
	);
	const locationParam = searchParams.get(FILTERS_PARAMS_NAMES.location);

	// Update page index when it's a valid number.
	if (!Number.isNaN(pageIndexParam)) {
		pageIndex = pageIndexParam;
	}

	// Update location when it's a non empty value.
	if (locationParam) {
		location = locationParam;
	}

	return {
		pageIndex,
		location,
	};
}

/**
 * Maps filters of the landfills table to URL search params.
 * @param filters Landfills filters.
 * @returns URL search params.
 */
function filtersToSearchParams(filters: LandfillsFilters): URLSearchParams {
	const { pageIndex, location } = filters;

	const searchParams = new URLSearchParams(window.location.search);
	searchParams.set(FILTERS_PARAMS_NAMES.pageIndex, pageIndex.toString());

	if (location) {
		searchParams.set(FILTERS_PARAMS_NAMES.location, location);
	} else {
		searchParams.delete(FILTERS_PARAMS_NAMES.location);
	}

	return searchParams;
}

/**
 * Retrieves landfills to be displayed in the landfills table.
 * @param filters Landfills filters.
 * @returns Landfills.
 */
async function getLandfills(
	filters: LandfillsFilters,
): Promise<PaginatedLandfills> {
	const { pageIndex, location } = filters;

	const res = await ecomapHttpClient.GET("/landfills", {
		params: {
			query: {
				offset: pageIndex * DEFAULT_PAGE_SIZE,
				limit: DEFAULT_PAGE_SIZE,
				sort: "createdAt",
				order: "desc",
				locationName: location || undefined,
			},
		},
	});

	if (res.error) {
		return { total: 0, landfills: [] };
	}

	return res.data;
}

const landfillsStore = createTableStore(
	BackOfficeRoutes.LANDFILLS,
	initialData,
	filtersToSearchParams,
	searchParamsToFilters,
	getLandfills,
);

export default landfillsStore;
