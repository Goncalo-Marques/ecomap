<script lang="ts">
	import { onMount } from "svelte";
	import Map from "../../../lib/components/map/Map.svelte";
	import { MapHelper } from "../../../lib/components/map/mapUtils";
	import OlMap from "ol/Map";

	let map: OlMap;

	onMount(() => {
		const mapHelper = new MapHelper(map);

		mapHelper.addClusterLayer(
			"/json/containers.geojson",
			"Containers",
			"#15803D",
		);
	});
</script>

<main>
	<Map bind:map layersContainer={true} />
</main>

<style>
	main {
		width: 100%;
	}
</style>
