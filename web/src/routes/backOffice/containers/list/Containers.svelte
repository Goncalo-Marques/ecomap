<script lang="ts">
	import { t } from "../../../../lib/utils/i8n";
	import Card from "../../components/Card.svelte";
	import ContainersTable from "./ContainersTable.svelte";
	import SearchContainers from "./SearchContainers.svelte";
</script>

<Card class="page-layout">
	<h1>{$t("sidebar.containers")}</h1>
	<SearchContainers />
	<ContainersTable />
</Card>

<style>
	h1 {
		font: var(--text-2xl-semibold);
	}
</style>
