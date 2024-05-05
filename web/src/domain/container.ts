import type { components } from "../../api/ecomap/http";

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
	location: string;
	category?: Container["category"];
}

/**
 * Category of a container.
 */
export type ContainerCategory = components["schemas"]["ContainerCategory"];
