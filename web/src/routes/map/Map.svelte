<script lang="ts">
	import Map from "../../lib/components/map/Map.svelte";
	import { onMount } from "svelte";
	import { addClusterLayer } from "../../lib/components/map/mapStore";

	onMount(() => {
		addClusterLayer("/json/contentores.geojson", "Pontos de recolha", "#15803D");
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
