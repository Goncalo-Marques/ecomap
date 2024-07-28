<script lang="ts">
	import Button from "$lib/components/Button.svelte";
	import { t } from "$lib/utils/i8n";

	import { BackOfficeRoutes } from "../../constants/routes";
	import Card from "../components/Card.svelte";
	import RoutesTable from "./list/RoutesTable.svelte";
	import SearchRoutes from "./list/SearchRoutes.svelte";
</script>

<Card class="m-10 flex flex-col gap-10">
	<h1 class="text-2xl font-semibold">{$t("routes")}</h1>
	<div class="flex items-center justify-between">
		<SearchRoutes />
		<a href={`${BackOfficeRoutes.ROUTES}/create`} class="contents">
			<Button startIcon="add">{$t("routes.create.title")}</Button>
		</a>
	</div>
	<RoutesTable />
</Card>
