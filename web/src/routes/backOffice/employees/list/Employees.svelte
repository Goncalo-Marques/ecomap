<script lang="ts">
	import { Link } from "svelte-routing";
	import { t } from "../../../../lib/utils/i8n";
	import Card from "../../components/Card.svelte";
	import EmployeesTable from "./EmployeesTable.svelte";
	import SearchEmployees from "./SearchEmployees.svelte";
	import Button from "../../../../lib/components/Button.svelte";
</script>

<Card class="m-10 flex flex-col gap-10">
	<h1 class="text-2xl font-semibold">{$t("employees")}</h1>
	<div class="flex items-center justify-between">
		<SearchEmployees />
		<Link to="create" class="contents">
			<Button startIcon="add">{$t("employees.create.title")}</Button>
		</Link>
	</div>
	<EmployeesTable />
</Card>
