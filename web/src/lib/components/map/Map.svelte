<script lang="ts">
	import { map, createMap } from "./mapStore";
	import { createEventDispatcher, onDestroy, onMount } from "svelte";
	import {} from "ol/ol.css";

	const dispatch = createEventDispatcher();

	/**
	 * Zoom value for map view
	 * @default 6.5
	 */
	export let zoom: number = 6.5;

	/**
	 * Center latitude of map
	 * @default 40
	 */
	export let lat: number = 40;

	/**
	 * Center longitude of map
	 * @default -7
	 */
	export let lon: number = -7;

	const map_id: string = "map_id";

	// Mount map into div
	onMount(() => {
		createMap(lon, lat, zoom);

		$map?.setTarget(map_id);

		$map?.getLayers().on("add", () => {
			dispatch("message", {
                layers : $map?.getAllLayers()
            });
		});
	});

	onDestroy(() => {
		map.set(null);
	});
</script>

<div id={map_id}>
	<a
		href="https://www.openstreetmap.org/copyright"
		target="_blank"
		rel="noopener noreferrer">© OpenStreetMap contributors</a
	>
</div>

<style>
	div {
		width: 100%;
		height: 100%;
		position: relative;
	}
	a {
		position: absolute;
		color: var(--gray-950);
		background-color: var(--white);
		font-size: 12px;
		padding: 1px 0.5em;
		border-radius: 0px 5px;
		opacity: 0.8;
		left: 0;
		bottom: 0;
		z-index: 999;
	}
</style>
