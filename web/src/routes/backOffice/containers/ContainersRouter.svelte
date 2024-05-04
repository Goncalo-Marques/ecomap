<script lang="ts">
	import { Route, Router } from "svelte-routing";
	import Containers from "./list/Containers.svelte";
	import Container from "./details/Container.svelte";
	import ContainerMap from "./details/ContainerMap.svelte";
	import EditContainer from "./edit/EditContainer.svelte";
</script>

<Router>
	<Route path="/:id/map" component={ContainerMap} />
	<Route path="/:id/edit" let:params>
		<EditContainer id={params.id} />
	</Route>
	<Route path="/:id" let:params>
		<Container id={params.id} />
	</Route>
	<Route path="/" component={Containers} />
</Router>
