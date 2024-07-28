<script lang="ts">
	import Button from "$lib/components/Button.svelte";
	import { t } from "$lib/utils/i8n";

	import { BackOfficeRoutes } from "../../constants/routes";
	import EmployeesTable from "./list/EmployeesTable.svelte";
	import SearchEmployees from "./list/SearchEmployees.svelte";
</script>

<h1 class="text-2xl font-semibold">{$t("employees")}</h1>
<div class="flex items-center justify-between">
	<SearchEmployees />
	<a href={`${BackOfficeRoutes.EMPLOYEES}/create`} class="contents">
		<Button startIcon="add">{$t("employees.create.title")}</Button>
	</a>
</div>
<EmployeesTable />
