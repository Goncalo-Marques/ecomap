import { derived, get, type Writable, writable } from "svelte/store";

import { page } from "$app/stores";
import type {
	DataFn,
	FiltersToSearchParams,
	SearchParamsToFilters,
	TableStore,
} from "$domain/stores/table";

import { updateSearchParams } from "../utils/url";

/**
 * Fetches data to be displayed on a table.
 * @param loading Loading store.
 * @param data Data store.
 * @param filters Filters.
 * @param dataFn Function that retrieves the data to be displayed in a table.
 */
async function fetchData<TData, TFilters>(
	loading: Writable<boolean>,
	data: Writable<TData>,
	filters: TFilters,
	dataFn: DataFn<TData, TFilters>,
) {
	loading.set(true);

	const responseData = await dataFn(filters);
	data.set(responseData);

	loading.set(false);
}

/**
 * Creates a store to be used to interact with a table component.
 * @param pathname Pathname of the location where the table is used. Used to determine when to request new data.
 * @param initialData Initial data of the store.
 * @param filtersToSearchParams Mapper of filters to URL search params.
 * @param searchParamsToFilters Mapper of URL search params to filters.
 * @param dataFn Function that retrieves the data to be displayed in a table.
 * @returns Table store.
 */
export function createTableStore<TData, TFilters>(
	pathname: string,
	initialData: TData,
	filtersToSearchParams: FiltersToSearchParams<TFilters>,
	searchParamsToFilters: SearchParamsToFilters<TFilters>,
	dataFn: DataFn<TData, TFilters>,
): TableStore<TData, TFilters> {
	const data = writable(initialData);

	const loading = writable(false);

	const filters = derived(page, () => {
		const searchParams = new URLSearchParams(location.search);
		const filtersFromSearchParams = searchParamsToFilters(searchParams);

		if (location.pathname === pathname) {
			fetchData(loading, data, filtersFromSearchParams, dataFn);
		} else {
			// Reset data store to its initial state when the user exits the pathname that the store belongs to.
			data.set(initialData);
		}

		return filtersFromSearchParams;
	});

	const store: TableStore<TData, TFilters> = {
		data: {
			subscribe: data.subscribe,
		},
		loading: {
			subscribe: loading.subscribe,
		},
		filters: {
			subscribe: filters.subscribe,
			set(value) {
				// Update the URL search params.
				updateSearchParams(filtersToSearchParams(value));
			},
			update(updater) {
				const updatedFilters = updater(get(filters));

				// Update the URL search params.
				updateSearchParams(filtersToSearchParams(updatedFilters));
			},
		},
	};

	return store;
}
