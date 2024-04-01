<script lang="ts">
	import { map, createMap } from "./mapStore";
	import { createEventDispatcher, onDestroy, onMount } from "svelte";
	import {} from "ol/ol.css";
	import type { Layer } from "ol/layer";
	import LayerItem from "./layerItem.svelte";

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

	let layers: Layer[] | undefined;

	// Mount map into div
	onMount(() => {
		createMap(lon, lat, zoom);

		$map?.setTarget(map_id);

		$map?.getLayers().on("add", () => {
			layers = $map?.getAllLayers();
		});
	});

	onDestroy(() => {
		map.set(null);
	});
</script>

<div id={map_id} class="map">
	<section class="layers" >
		<header>
			<h2>Layers</h2>
			<hr />
		</header>
		<div class="item-container">
			{#if layers}
				{#each layers as layer}
					{#if layer.get("layer-name") != "baseLayer"}
						<LayerItem {layer} />
					{/if}
				{/each}
			{/if}
		</div>
	</section>

	<!-- <a
		href="https://www.openstreetmap.org/copyright"
		target="_blank"
		rel="noopener noreferrer">© OpenStreetMap contributors</a
	> -->
</div>

<style> 
	* {
		box-sizing: border-box;
	}

	.map {
		width: 100%;
		height: 100%;
		position: relative;
	}
	
	/* a {
		position: absolute;
		color: var(--gray-950);
		background-color: var(--white);
		font-size: 12px;
		padding: 1px 0.5em;
		border-radius: 0px 5px 0px 0px;
		opacity: 0.8;
		left: 0;
		bottom: 0;
		z-index: 500;
	} */

	.layers {
		position: absolute;
		min-width: 360px;
		background-color: white;
		z-index: 999;
		top: 3em;
		left: 3em;
		border-radius: 4px;
		padding: 10px;
	}

	.item-container {
		display: flex;
		flex-direction: column-reverse;
		gap: 5px;
	}

	.item-container::after,
	.item-container::before {
		content: "";
	}
</style>
