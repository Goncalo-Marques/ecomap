<script lang="ts">
	import { Router, Route } from "svelte-routing";
	import SignIn from "./routes/SignIn.svelte";
	import MapTest from "./routes/MapTest.svelte";
</script>

<Router>
	<Route path="/signin" component={SignIn} />
	<Route path="/map_test" component={MapTest} />
</Router>
