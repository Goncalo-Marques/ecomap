<script lang="ts">
	import Search from "../../../../lib/components/Search.svelte";
	import debounce from "../../../../lib/utils/debounce";
	import { t } from "../../../../lib/utils/i8n";
	import warehousesStore from "./warehousesStore";

	const { filters } = warehousesStore;

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
				location: value,
			};
		});
	}
</script>

<div class="search">
	<Search
		value={$filters.location}
		placeholder={$t("location.search")}
		onInput={debounce(handleSearchChange)}
	/>
</div>

<style>
	.search {
		width: 16rem;
	}
</style>
