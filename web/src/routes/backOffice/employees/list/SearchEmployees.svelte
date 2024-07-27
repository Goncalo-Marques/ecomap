<script lang="ts">
	import Search from "../../../../lib/components/Search.svelte";
	import debounce from "../../../../lib/utils/debounce";
	import { t } from "../../../../lib/utils/i8n";
	import employeesStore from "./employeesStore";

	const { filters } = employeesStore;

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
				username: value,
			};
		});
	}
</script>

<Search
	class="w-80"
	value={$filters.username}
	placeholder={$t("employees.search")}
	onInput={debounce(handleSearchChange)}
/>
