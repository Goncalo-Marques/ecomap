<script lang="ts">
	import { Route, Router } from "svelte-routing";
	import Employees from "./list/Employees.svelte";
	import Employee from "./details/Employee.svelte";
	import EditEmployee from "./edit/EditEmployee.svelte";
	import CreateEmployee from "./create/CreateEmployee.svelte";
</script>

<Router>
	<Route path="/create">
		<CreateEmployee />
	</Route>
	<Route path="/:id/edit" let:params>
		<EditEmployee id={params.id} />
	</Route>
	<Route path="/:id" let:params>
		<Employee id={params.id} />
	</Route>
	<Route path="/" component={Employees} />
</Router>
