<script lang="ts">
	import { Link } from "svelte-routing";
	import { t } from "../../../../lib/utils/i8n";
	import Card from "../../components/Card.svelte";
	import RoutesTable from "./RoutesTable.svelte";
	import SearchRoutes from "./SearchRoutes.svelte";
	import Button from "../../../../lib/components/Button.svelte";
</script>

<Card class="page-layout">
	<h1>{$t("routes")}</h1>
	<div class="top-bar">
		<SearchRoutes />
		<Link to="create" style="display:contents">
			<Button startIcon="add">{$t("routes.create.title")}</Button>
		</Link>
	</div>
	<RoutesTable />
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
