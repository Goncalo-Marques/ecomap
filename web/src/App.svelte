<script lang="ts">
	import { Router, Route } from "svelte-routing";
	import SignIn from "./routes/SignIn.svelte";
</script>

<Router>
	<Route path="/signin" component={SignIn} />
</Router>
