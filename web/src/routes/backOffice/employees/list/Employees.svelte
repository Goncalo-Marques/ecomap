<script lang="ts">
	import { t } from "../../../../lib/utils/i8n";
	import Card from "../../components/Card.svelte";
	import EmployeesTable from "./EmployeesTable.svelte";
	import SearchEmployees from "./SearchEmployees.svelte";
</script>

<Card class="page-layout">
	<h1>{$t("sidebar.employees")}</h1>
	<SearchEmployees />
	<EmployeesTable />
</Card>

<style>
	h1 {
		font: var(--text-2xl-semibold);
	}
</style>
