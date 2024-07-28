<script lang="ts">
	import Search from "$lib/components/Search.svelte";
	import debounce from "$lib/utils/debounce";
	import { t } from "$lib/utils/i8n";

	import routesStore from "./routesStore";

	const { filters } = routesStore;

	/**
	 * Handles the change of value of the search input.
	 * @param e Input change event.
	 */
	function handleSearchChange(e: Event) {
		const searchInput = e.target as HTMLInputElement;
		const value = searchInput.value;

		filters.update(filters => {
			return {
				...filters,
				pageIndex: 0,
				route: value,
			};
		});
	}
</script>

<Search
	class="w-64"
	value={$filters.route}
	placeholder={$t("routes.search")}
	onInput={debounce(handleSearchChange)}
/>
