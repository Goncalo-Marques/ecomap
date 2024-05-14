<script lang="ts">
	import { Route, Router } from "svelte-routing";

	import Employees from "./list/Employees.svelte";
</script>

<Router>
	<Route path="/" component={Employees} />
</Router>
