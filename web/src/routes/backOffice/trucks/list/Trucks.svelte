<script lang="ts">
	import { Link } from "svelte-routing";
	import { t } from "../../../../lib/utils/i8n";
	import Card from "../../components/Card.svelte";
	import TrucksTable from "./TrucksTable.svelte";
	import SearchTrucks from "./SearchTrucks.svelte";
	import Button from "../../../../lib/components/Button.svelte";
</script>

<Card class="page-layout">
	<h1>{$t("trucks")}</h1>
	<div class="top-bar">
		<SearchTrucks />
		<Link to="create" style="display:contents">
			<Button startIcon="add">{$t("trucks.create.title")}</Button>
		</Link>
	</div>
	<TrucksTable />
</Card>

<style>
	h1 {
		font: var(--text-2xl-semibold);
	}

	.top-bar {
		display: flex;
		align-items: center;
		justify-content: space-between;
	}
</style>
