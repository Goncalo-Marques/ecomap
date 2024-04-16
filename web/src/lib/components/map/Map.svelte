<script lang="ts">
	import { map, createMap } from "./mapStore";
	import { onDestroy, onMount } from "svelte";
	import {} from "ol/ol.css";
	import type { Layer } from "ol/layer";
	import LayerItem from "./layerItem.svelte";
	import Icon from "../Icon.svelte";
	import { t } from "../../utils/i8n";

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

	/**
	 * Show/Hide layers container
	 * @default false
	 */
	export let layersContainer: boolean = false;

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
	{#if layersContainer && layers}
		<section class="layers">
			<header>
				<Icon name="layers" />
				<h1>{$t("map.layers")}</h1>
			</header>
			<div class="item-container">
				{#each layers as layer}
					{#if layer.get("layer-name") != "baseLayer"}
						<LayerItem {layer} />
					{/if}
				{/each}
			</div>
		</section>
	{/if}
</div>

<style>
	* {
		box-sizing: border-box;
	}

	header {
		display: flex;
		align-items: center;
		gap: 4px;
	}

	h1 {
		font: var(--text-base-semibold);
	}

	.map {
		width: 100%;
		height: 100%;
		position: relative;
	}

	.layers {
		position: absolute;
		background-color: var(--white);
		min-width: 256px;
		max-height: 300px;
		border-radius: 4px;
		padding: 10px;
		overflow: auto;
		bottom: 3em;
		left: 3em;
		z-index: 999;
	}

	.item-container {
		display: flex;
		flex-direction: column;
		gap: 8px;
		margin-top: 12px;
	}
</style>
