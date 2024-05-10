import type { components } from "../../api/ecomap/http";
import type { SortingDirection } from "../lib/components/table/types";

/**
 * Warehouse.
 */
export type Warehouse = components["schemas"]["Warehouse"];

/**
 * Sortable fields of a warehouse.
 */
export type WarehouseSortableFields = NonNullable<
	components["parameters"]["WarehouseSortQueryParam"]
>;

/**
 * Paginated warehouses.
 */
export type PaginatedWarehouses = components["schemas"]["WarehousesPaginated"];

/**
 * Filters of warehouses.
 */
export interface WarehousesFilters {
	pageIndex: number;
	sort: WarehouseSortableFields;
	order: SortingDirection;
	location: string;
}
