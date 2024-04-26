<script lang="ts">
	import { t } from "../../../../lib/utils/i8n";
	import ContainersTable from "./ContainersTable.svelte";
</script>

<h1>{$t("sidebar.containers")}</h1>
<ContainersTable />

<style>
	h1 {
		font: var(--text-2xl-semibold);
	}
</style>
