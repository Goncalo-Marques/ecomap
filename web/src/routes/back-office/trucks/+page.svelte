<script lang="ts">
	import Button from "$lib/components/Button.svelte";
	import { BackOfficeRoutes } from "$lib/constants/routes";
	import { t } from "$lib/utils/i8n";

	import Card from "../components/Card.svelte";
	import SearchTrucks from "./list/SearchTrucks.svelte";
	import TrucksTable from "./list/TrucksTable.svelte";
</script>

<Card class="m-10 flex flex-col gap-10">
	<h1 class="text-2xl font-semibold">{$t("trucks")}</h1>
	<div class="flex items-center justify-between">
		<SearchTrucks />
		<a href={`${BackOfficeRoutes.TRUCKS}/create`} class="contents">
			<Button startIcon="add">{$t("trucks.create.title")}</Button>
		</a>
	</div>
	<TrucksTable />
</Card>
