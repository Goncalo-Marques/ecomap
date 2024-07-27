<script lang="ts">
	import Search from "../../../../lib/components/Search.svelte";
	import debounce from "../../../../lib/utils/debounce";
	import { t } from "../../../../lib/utils/i8n";
	import trucksStore from "./trucksStore";

	const { filters } = trucksStore;

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
				licensePlate: value,
			};
		});
	}
</script>

<Search
	class="w-64"
	value={$filters.licensePlate}
	placeholder={$t("licensePlate.search")}
	onInput={debounce(handleSearchChange)}
/>
