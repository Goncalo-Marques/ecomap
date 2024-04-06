<script lang="ts">
	import { map, createMap } from "./mapStore";
	import { onDestroy, onMount } from "svelte";
	import {} from "ol/ol.css";
	import type { Layer } from "ol/layer";
	import LayerItem from "./layerItem.svelte";

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

	.layers {
		min-width: 256px;
		z-index: 999;
		position: absolute;
		background-color: var(--white);
		border-radius: 4px;
		max-height: 300px;

		bottom: 3em;
		left: 3em;
		padding: 10px;
		overflow: auto;

	}

	.item-container {
		display: flex;
		flex-direction: column;
		gap: 5px;
	}

	.item-container::after,
	.item-container::before {
		content: "";
	}
</style>
