<script lang="ts">
	import { onMount } from "svelte";
	import OlMap from "ol/Map";
	import GeoJSON from "ol/format/GeoJSON";
	import Map from "../../../lib/components/map/Map.svelte";
	import { MapHelper } from "../../../lib/components/map/mapUtils";

	let map: OlMap;

	onMount(() => {
		const mapHelper = new MapHelper(map);

		mapHelper.addClusterLayer(
			{
				url: "/json/containers.geojson",
				format: new GeoJSON(),
			},
			"Containers",
			"#15803D",
		);
	});
</script>

<main>
	<Map bind:map showLayers />
</main>

<style>
	main {
		width: 100%;
	}
</style>
