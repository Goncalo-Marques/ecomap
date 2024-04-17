<script lang="ts">
	import { onMount } from "svelte";
	import Map from "../../../lib/components/map/Map.svelte";
	import type { MapHelper } from "../../../lib/components/map/mapUtils";

	let mapHelper: MapHelper;

	onMount(() => {
		mapHelper.addClusterLayer(
			"/json/containers.geojson",
			"Containers",
			"#15803D",
		);
	});
</script>

<main>
	<Map bind:mapHelper layersContainer={true} />
</main>

<style>
	main {
		width: 100%;
	}
</style>
