<script lang="ts">
	import { Route, Router } from "svelte-routing";
	import Card from "./../components/Card.svelte";
	import Containers from "./list/Containers.svelte";
	import Container from "./details/Container.svelte";
</script>

<Router>
	<Card class="page-layout">
		<Route path="/:id/edit" let:params>
			<Container id={params.id} mode="edit" />
		</Route>
		<Route path="/:id" let:params>
			<Container id={params.id} mode="view" />
		</Route>
		<Route path="/" component={Containers} />
	</Card>
</Router>
