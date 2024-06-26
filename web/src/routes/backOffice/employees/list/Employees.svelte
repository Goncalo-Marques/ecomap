<script lang="ts">
	import { Link } from "svelte-routing";
	import { t } from "../../../../lib/utils/i8n";
	import Card from "../../components/Card.svelte";
	import EmployeesTable from "./EmployeesTable.svelte";
	import SearchEmployees from "./SearchEmployees.svelte";
	import Button from "../../../../lib/components/Button.svelte";
</script>

<Card class="page-layout">
	<h1>{$t("employees")}</h1>
	<div class="top-bar">
		<SearchEmployees />
		<Link to="create" style="display:contents">
			<Button startIcon="add">{$t("employees.create.title")}</Button>
		</Link>
	</div>
	<EmployeesTable />
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
