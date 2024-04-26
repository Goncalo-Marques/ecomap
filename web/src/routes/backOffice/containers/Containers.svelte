<script lang="ts">
	import { t } from "../../../lib/utils/i8n";
	import Card from "../components/Card.svelte";
	import ContainersTable from "./list/ContainersTable.svelte";
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
