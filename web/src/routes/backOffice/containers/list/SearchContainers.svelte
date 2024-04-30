<script lang="ts">
	import Search from "../../../../lib/components/Search.svelte";
	import debounce from "../../../../lib/utils/debounce";
	import { t } from "../../../../lib/utils/i8n";
	import containersStore from "./containersStore";

	function handleSearchChange(e: Event) {
		const searchInput = e.target as HTMLInputElement;
		const value = searchInput.value;

		containersStore.filters.update(filters => {
			return {
				...filters,
				pageIndex: 0,
				location: value,
			};
		});
	}

	const { filters } = containersStore;
</script>

<div class="search">
	<Search
		value={$filters.location}
		placeholder={$t("containers.search")}
		onInput={debounce(handleSearchChange)}
	/>
</div>

<style>
	.search {
		width: 16rem;
	}
</style>
