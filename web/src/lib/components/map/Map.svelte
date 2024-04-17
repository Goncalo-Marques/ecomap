<script lang="ts">
	import { MapHelper, createMap } from "./mapUtils";
	import { onMount } from "svelte";
	import {} from "ol/ol.css";
	import type { Layer } from "ol/layer";
	import LayerItem from "./LayerItem.svelte";
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

	/**
	 * MapHelper Object, contains OpenLayers Map object and helper methods
	 */
	export let mapHelper: MapHelper;

	const map_id: string = "map_id";

	let layers: Layer[] | undefined;

	onMount(() => {
		mapHelper = createMap(lon, lat, zoom);

		mapHelper.map.setTarget(map_id);

		mapHelper.map.getLayers().on("add", () => {
			layers = mapHelper.map.getAllLayers();
		});
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
	header {
		display: flex;
		align-items: center;
		gap: 0.25rem;
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
		min-width: 16rem;
		max-height: 37.5rem;
		border-radius: 4px;
		padding: 0.625rem;
		overflow: auto;
		bottom: 3rem;
		left: 3rem;
		z-index: 999;
	}

	.item-container {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		margin-top: 0.75rem;
	}
</style>
