<script lang="ts">
	import { Route, Router } from "svelte-routing";
	import Trucks from "./list/Trucks.svelte";
	import Truck from "./details/Truck.svelte";
	import TruckMap from "./details/TruckMap.svelte";
	import CreateTruck from "./create/CreateTruck.svelte";
	import EditTruck from "./edit/EditTruck.svelte";
</script>

<Router>
	<Route path="/create" component={CreateTruck} />
	<Route path="/:id/map" component={TruckMap} />
	<Route path="/:id/edit" component={EditTruck} />
	<Route path="/:id" component={Truck} />
	<Route path="/" component={Trucks} />
</Router>
