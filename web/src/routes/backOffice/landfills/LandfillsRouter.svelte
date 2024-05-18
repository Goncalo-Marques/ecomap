<script lang="ts">
	import { Route, Router } from "svelte-routing";
	import Landfills from "./list/Landfills.svelte";
	import Landfill from "./details/Landfill.svelte";
	import LandfillMap from "./details/LandfillMap.svelte";
	import CreateLandfill from "./create/CreateLandfill.svelte";
	import EditLandfill from "./edit/EditLandfill.svelte";
</script>

<Router>
	<Route path="/create" component={CreateLandfill} />
	<Route path="/:id/map" component={LandfillMap} />
	<Route path="/:id/edit" component={EditLandfill} />
	<Route path="/:id" component={Landfill} />
	<Route path="/" component={Landfills} />
</Router>
