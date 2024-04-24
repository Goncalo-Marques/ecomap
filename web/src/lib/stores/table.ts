import { get, writable } from "svelte/store";
import type {
	FiltersToSearchParams,
	SearchParamsToFilters,
	DataFn,
	TableStore,
} from "../../domain/stores/table";
import { updateSearchParams } from "../utils/url";

/**
 * Creates a store to be used to interact with a table component.
 * @param initialData Initial data of the store.
 * @param filtersToSearchParams Mapper of filters to URL search params.
 * @param searchParamsToFilters Mapper of URL search params to filters.
 * @param dataFn Function that retrieves the data to be displayed in a table.
 * @returns Table store.
 */
export function createTableStore<TData, TFilters>(
	initialData: TData,
	filtersToSearchParams: FiltersToSearchParams<TFilters>,
	searchParamsToFilters: SearchParamsToFilters<TFilters>,
	dataFn: DataFn<TData, TFilters>,
): TableStore<TData, TFilters> {
	const data = writable(initialData);

	const loading = writable(false);

	const searchParams = new URLSearchParams(location.search);
	const filters = writable(searchParamsToFilters(searchParams));
	const { set: setFilters } = filters;

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
				setFilters(value);

				// Fetch new data after setting the filters.
				store.fetchData();
			},
			update(updater) {
				const updatedFilters = updater(get(filters));

				setFilters(updatedFilters);

				// Fetch new data after updating the filters.
				store.fetchData();
			},
		},
		async fetchData() {
			loading.set(true);

			const currentFilters = get(filters);

			// Updates the URL search params with the current filters.
			updateSearchParams(filtersToSearchParams(currentFilters));

			const responseData = await dataFn(currentFilters);
			data.set(responseData);

			loading.set(false);
		},
	};

	return store;
}
