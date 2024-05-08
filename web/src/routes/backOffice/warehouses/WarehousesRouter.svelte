<script lang="ts">
	import { Route, Router } from "svelte-routing";
	import Warehouses from "./list/Warehouses.svelte";
	import Warehouse from "./details/Warehouse.svelte";
	import WarehouseMap from "./details/WarehouseMap.svelte";
	import CreateWarehouse from "./create/CreateWarehouse.svelte";
	import EditWarehouse from "./edit/EditWarehouse.svelte";
</script>

<Router>
	<Route path="/create" component={CreateWarehouse} />
	<Route path="/:id/map" component={WarehouseMap} />
	<Route path="/:id/edit" component={EditWarehouse} />
	<Route path="/:id" component={Warehouse} />
	<Route path="/" component={Warehouses} />
</Router>
