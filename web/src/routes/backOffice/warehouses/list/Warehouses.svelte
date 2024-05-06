<script lang="ts">
	import { Link } from "svelte-routing";
	import { t } from "../../../../lib/utils/i8n";
	import Card from "../../components/Card.svelte";
	import WarehousesTable from "./WarehousesTable.svelte";
	import SearchWarehouses from "./SearchWarehouses.svelte";
	import Button from "../../../../lib/components/Button.svelte";
</script>

<Card class="page-layout">
	<h1>{$t("warehouses")}</h1>
	<div class="top-bar">
		<SearchWarehouses />
		<Link to="create" style="display:contents">
			<Button startIcon="add">{$t("warehouses.create.title")}</Button>
		</Link>
	</div>
	<WarehousesTable />
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
