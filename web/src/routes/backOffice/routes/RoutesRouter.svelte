<script lang="ts">
	import { Route as RouterRoute, Router } from "svelte-routing";
	import Routes from "./list/Routes.svelte";
	import Route from "./details/Route.svelte";
	import RouteMap from "./details/RouteMap.svelte";
	import CreateRoute from "./create/CreateRoute.svelte";
	import EditRoute from "./edit/EditRoute.svelte";
</script>

<Router>
	<RouterRoute path="/create" component={CreateRoute} />
	<RouterRoute path="/:id/edit" component={EditRoute} />
	<RouterRoute path="/:id/map" component={RouteMap} />
	<RouterRoute path="/:id" component={Route} />
	<RouterRoute path="/" component={Routes} />
</Router>
