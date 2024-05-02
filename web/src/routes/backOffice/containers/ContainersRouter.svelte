<script lang="ts">
	import { Route, Router } from "svelte-routing";
	import Containers from "./list/Containers.svelte";
	import Container from "./details/Container.svelte";
	import ContainerMap from "./details/ContainerMap.svelte";
</script>

<Router>
	<Route path="/:id/map" component={ContainerMap} />
	<Route path="/:id" component={Container} />
	<Route path="/" component={Containers} />
</Router>
