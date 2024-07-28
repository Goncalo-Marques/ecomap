<script lang="ts">
	import Button from "$lib/components/Button.svelte";
	import { BackOfficeRoutes } from "$lib/constants/routes";
	import { t } from "$lib/utils/i8n";

	import Card from "../components/Card.svelte";
	import SearchWarehouses from "./list/SearchWarehouses.svelte";
	import WarehousesTable from "./list/WarehousesTable.svelte";
</script>

<Card class="m-10 flex flex-col gap-10">
	<h1 class="text-2xl font-semibold">{$t("warehouses")}</h1>
	<div class="flex items-center justify-between">
		<SearchWarehouses />
		<a href={`${BackOfficeRoutes.WAREHOUSES}/create`} class="contents">
			<Button startIcon="add">{$t("warehouses.create.title")}</Button>
		</a>
	</div>
	<WarehousesTable />
</Card>
