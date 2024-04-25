<script lang="ts">
	import { onDestroy, onMount } from "svelte";
	import { t } from "../../../../lib/utils/i8n";
	import ContainersTable from "./ContainersTable.svelte";
	import containersStore, { initialFilters } from "./containersStore";
	import { clearSearchParams } from "../../../../lib/utils/url";

	onMount(() => {
		containersStore.fetchData();
	});

	onDestroy(() => {
		// Reset containers filters to its initial state.
		containersStore.filters.set(initialFilters);

		// Clear containers filters search params from URL.
		clearSearchParams();
	});
</script>

<h1>{$t("sidebar.containers")}</h1>
<ContainersTable />

<style>
	h1 {
		font: var(--text-2xl-semibold);
	}
</style>
