<script lang="ts">
	import Map from "../../lib/components/map/Map.svelte";
	import { onDestroy, onMount } from "svelte";
	import {
		map,
		addClusterLayer,
		addVectorLayer,
	} from "../../lib/components/map/mapStore";
	onMount(() => {
		addClusterLayer("http://localhost:8000/contentores.geojson", "Pontos de recolha", "green");
		addClusterLayer("http://localhost:8000/contentores.geojson", "Pontos de recolha 1", "blue");
		addClusterLayer("http://localhost:8000/contentores.geojson", "Pontos de recolha 2", "red");
		addClusterLayer("http://localhost:8000/contentores.geojson", "Pontos de recolha 3", "lightgreen");
		addClusterLayer("http://localhost:8000/contentores.geojson", "Pontos de recolha 4", "black");
	});
</script>

<main>
	<Map /> 
</main>

<style>
	main {
		width: 100%;
	}
</style>
