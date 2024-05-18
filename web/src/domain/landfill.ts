import type { components } from "../../api/ecomap/http";

/**
 * Landfill.
 */
export type Landfill = components["schemas"]["Landfill"];

/**
 * Paginated landfills.
 */
export type PaginatedLandfills = components["schemas"]["LandfillsPaginated"];

/**
 * Filters of landfills.
 */
export interface LandfillsFilters {
	pageIndex: number;
	location: string;
}
