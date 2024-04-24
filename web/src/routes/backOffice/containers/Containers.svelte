<script lang="ts">
	import { onDestroy, onMount } from "svelte";
	import { t } from "../../../lib/utils/i8n";
	import Card from "../components/Card.svelte";
	import ContainersTable from "./list/ContainersTable.svelte";
	import containersStore, { initialFilters } from "./list/containersStore";
	import { clearSearchParams } from "../../../lib/utils/url";

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

<Card class="card">
	<h1>{$t("sidebar.containers")}</h1>
	<ContainersTable />
</Card>

<style>
	:global(.card) {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: 2.5rem;
		margin: 2.5rem;
	}

	h1 {
		font: var(--text-2xl-semibold);
	}
</style>
