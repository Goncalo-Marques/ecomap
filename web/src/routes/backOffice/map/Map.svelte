<script lang="ts">
	import { onMount } from "svelte";
	import { addClusterLayer } from "../../../lib/components/map/mapStore";
	import Map from "../../../lib/components/map/Map.svelte";

	onMount(() => {
		addClusterLayer("/json/containers.geojson", "Containers", "#15803D");
	});
</script>

<main>
	<Map layersContainer={true} />
</main>

<style>
	main {
		width: 100%;
	}
</style>
