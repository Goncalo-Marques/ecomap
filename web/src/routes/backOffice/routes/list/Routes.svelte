<script lang="ts">
	import { Link } from "svelte-routing";
	import { t } from "../../../../lib/utils/i8n";
	import Card from "../../components/Card.svelte";
	import RoutesTable from "./RoutesTable.svelte";
	import SearchRoutes from "./SearchRoutes.svelte";
	import Button from "../../../../lib/components/Button.svelte";
</script>

<Card class="m-10 flex flex-col gap-10">
	<h1 class="text-2xl font-semibold">{$t("routes")}</h1>
	<div class="flex items-center justify-between">
		<SearchRoutes />
		<Link to="create" class="contents">
			<Button startIcon="add">{$t("routes.create.title")}</Button>
		</Link>
	</div>
	<RoutesTable />
</Card>
