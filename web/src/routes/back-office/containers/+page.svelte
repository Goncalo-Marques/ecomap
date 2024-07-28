<script lang="ts">
	import Button from "$lib/components/Button.svelte";
	import { t } from "$lib/utils/i8n";

	import { BackOfficeRoutes } from "../../constants/routes";
	import Card from "../components/Card.svelte";
	import ContainersTable from "./list/ContainersTable.svelte";
	import SearchContainers from "./list/SearchContainers.svelte";
</script>

<Card class="m-10 flex flex-col gap-10">
	<h1 class="text-2xl font-semibold">{$t("containers")}</h1>
	<div class="flex items-center justify-between">
		<SearchContainers />
		<a href={`${BackOfficeRoutes.CONTAINERS}/create`} class="contents">
			<Button startIcon="add">{$t("containers.create.title")}</Button>
		</a>
	</div>
	<ContainersTable />
</Card>
