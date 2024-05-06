import type { components } from "../../api/ecomap/http";

/**
 * Warehouse.
 */
export type Warehouse = components["schemas"]["Warehouse"];

/**
 * Paginated warehouses.
 */
export type PaginatedWarehouses = components["schemas"]["WarehousesPaginated"];

/**
 * Filters of warehouses.
 */
export interface WarehousesFilters {
	pageIndex: number;
	location: string;
}
