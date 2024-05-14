<script lang="ts">
	import { Link } from "svelte-routing";
	import { t } from "../../../../lib/utils/i8n";
	import Card from "../../components/Card.svelte";
	import ContainersTable from "./ContainersTable.svelte";
	import SearchContainers from "./SearchContainers.svelte";
	import Button from "../../../../lib/components/Button.svelte";
</script>

<Card class="page-layout">
	<h1>{$t("containers")}</h1>
	<div class="top-bar">
		<SearchContainers />
		<Link to="create" style="display:contents">
			<Button startIcon="add">{$t("containers.create.title")}</Button>
		</Link>
	</div>
	<ContainersTable />
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
