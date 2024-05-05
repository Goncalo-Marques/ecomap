<script lang="ts">
	import { Route, Router } from "svelte-routing";
	import Employees from "./list/Employees.svelte";
	import Employee from "./details/Employee.svelte";
</script>

<Router>
	<Route path="/:id" let:params>
		<Employee id={params.id} />
	</Route>
	<Route path="/" component={Employees} />
</Router>
