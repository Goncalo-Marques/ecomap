<script lang="ts">
	import { Link } from "svelte-routing";
	import { t } from "../../../../lib/utils/i8n";
	import Card from "../../components/Card.svelte";
	import LandfillsTable from "./LandfillsTable.svelte";
	import SearchLandfills from "./SearchLandfills.svelte";
	import Button from "../../../../lib/components/Button.svelte";
</script>

<Card class="page-layout">
	<h1>{$t("landfills")}</h1>
	<div class="top-bar">
		<SearchLandfills />
		<Link to="create" style="display:contents">
			<Button startIcon="add">{$t("landfills.create.title")}</Button>
		</Link>
	</div>
	<LandfillsTable />
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
