<script lang="ts">
	import { Route, Router } from "svelte-routing";
	import Card from "./../components/Card.svelte";
	import Containers from "./list/Containers.svelte";
	import Container from "./details/Container.svelte";
</script>

<Router>
	<Card class="page-layout">
		<Route path="/:id" let:params>
			<Container id={params.id} />
		</Route>
		<Route path="/" component={Containers} />
	</Card>
</Router>
