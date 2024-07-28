<script lang="ts">
	import Button from "$lib/components/Button.svelte";
	import { t } from "$lib/utils/i8n";

	import { BackOfficeRoutes } from "../../constants/routes";
	import Card from "../components/Card.svelte";
	import LandfillsTable from "./list/LandfillsTable.svelte";
	import SearchLandfills from "./list/SearchLandfills.svelte";
</script>

<Card class="m-10 flex flex-col gap-10">
	<h1 class="text-2xl font-semibold">{$t("landfills")}</h1>
	<div class="flex items-center justify-between">
		<SearchLandfills />
		<a href={`${BackOfficeRoutes.LANDFILLS}/create`} class="contents">
			<Button startIcon="add">{$t("landfills.create.title")}</Button>
		</a>
	</div>
	<LandfillsTable />
</Card>
