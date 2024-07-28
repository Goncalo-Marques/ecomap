import type { components } from "$api/ecomap/http";
import type { SortingDirection } from "$lib/components/table/types";

/**
 * Truck.
 */
export type Truck = components["schemas"]["Truck"];

/**
 * Sortable fields of a truck.
 */
export type TruckSortableFields = NonNullable<
	components["parameters"]["TruckSortQueryParam"]
>;

/**
 * Paginated trucks.
 */
export type PaginatedTrucks = components["schemas"]["TrucksPaginated"];

/**
 * Filters of trucks.
 */
export interface TrucksFilters {
	pageIndex: number;
	sort: TruckSortableFields;
	order: SortingDirection;
	licensePlate: string;
}
