import type { components } from "../../api/ecomap/http";
import type { SortingDirection } from "../lib/components/table/types";

/**
 * Container.
 */
export type Container = components["schemas"]["Container"];

/**
 * Sortable fields of a container.
 */
export type ContainerSortableFields = NonNullable<
	components["parameters"]["ContainerSortQueryParam"]
>;

/**
 * Paginated containers.
 */
export type PaginatedContainers = components["schemas"]["ContainersPaginated"];

/**
 * Filters of containers.
 */
export interface ContainersFilters {
	pageIndex: number;
	sort: ContainerSortableFields;
	order: SortingDirection;
	location: string;
}
