import type { Readable } from "svelte/motion";
import type { Writable } from "svelte/store";

/**
 * Represents a table store.
 * @template TData Type of the data.
 * @template TFilters Type of the filters.
 */
export interface TableStore<TData, TFilters> {
	/**
	 * Data store.
	 * Used to store the data that is displayed on a table.
	 */
	data: Readable<TData>;

	/**
	 * Loading store.
	 * Indicates if data is being loaded in.
	 */
	loading: Readable<boolean>;

	/**
	 * Filters store.
	 * Used to store the filters of a table.
	 */
	filters: Writable<TFilters>;

	/**
	 * Requests new data.
	 */
	fetchData(): Promise<void>;
}

/**
 * Represents a mapper of filters to URL search params.
 * @template TFilters Type of the filters.
 * @param filters Filters.
 * @returns Search params used in the URL.
 */
export type FiltersToSearchParams<TFilters> = (
	filters: TFilters,
) => URLSearchParams;

/**
 * Represents a mapper of URL search params to filters.
 * @template TFilters Type of the filters.
 * @param searchParams Search params in the URL.
 * @returns Filters.
 */
export type SearchParamsToFilters<TFilters> = (
	searchParams: URLSearchParams,
) => TFilters;

/**
 * Represents the function that retrieves data.
 * @template TData Type of the data.
 * @template TFilters Type of the filters.
 * @returns Data to be displayed in a table.
 */
export type DataFn<TData, TFilters> = (filters: TFilters) => Promise<TData>;
