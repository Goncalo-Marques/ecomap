<script lang="ts">
	import { createMap } from "./mapUtils";
	import { onMount } from "svelte";
	import "ol/ol.css";
	import type { Layer } from "ol/layer";
	import LayerItem from "./LayerItem.svelte";
	import Icon from "../Icon.svelte";
	import { t } from "../../utils/i8n";
	import Map from "ol/Map";
	import {
		DEFAULT_MAX_ZOOM,
		DEFAULT_MIN_ZOOM,
		mapLayerName,
		nameLayerKey,
	} from "../../constants/map";

	/**
	 * Zoom value for map view.
	 *
	 * @default 6.5
	 */
	export let zoom: number = 6.5;

	/**
	 * Minimum zoom value for map view.
	 *
	 * @default 2
	 */
	export let minZoom: number = DEFAULT_MIN_ZOOM;

	/**
	 * Maximum zoom value for map view.
	 *
	 * @default 18
	 */
	export let maxZoom: number = DEFAULT_MAX_ZOOM;

	/**
	 * Center latitude of map.
	 *
	 * @default 40
	 */
	export let lat: number = 40;

	/**
	 * Center longitude of map.
	 *
	 * @default -7
	 */
	export let lon: number = -7;

	/**
	 * Indicates whether layers are displayed.
	 *
	 * @default false
	 */
	export let showLayers: boolean = false;

	/**
	 * Projection used.
	 *
	 * @default "EPSG:3857"
	 */
	export let projection: string = "EPSG:3857";

	/**
	 * Open Layers map.
	 */
	export let map: Map;

	/**
	 * Map container ID.
	 *
	 * @default "map_id"
	 */
	export let mapId: string = "map_id";

	/**
	 * Callback fired when map is initialized.
	 *
	 * @default null
	 */
	export let onInit: ((map: Map) => void) | null = null;

	/**
	 * Map layers.
	 */
	let layers: Layer[] = [];

	onMount(() => {
		map = createMap({
			lon,
			lat,
			zoom,
			projection,
			maxZoom,
			minZoom,
		});

		map.setTarget(mapId);
		map.getLayers().on("add", () => {
			layers = map.getAllLayers();
		});

		onInit?.(map);
	});
</script>

<div id={mapId} class="map">
	{#if showLayers && layers.length}
		<section class="layers">
			<header>
				<Icon name="layers" />
				<h1>{$t("layers")}</h1>
			</header>
			<div class="item-container">
				{#each layers as layer}
					{#if layer.get(nameLayerKey) != mapLayerName}
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
		width: var(--width, 100%);
		height: var(--height, 100%);
		max-width: var(--max-width, 100%);
		max-height: var(--max-height, 100%);
		position: relative;
	}

	.layers {
		position: absolute;
		background-color: var(--white);
		min-width: 16rem;
		max-height: 37.5rem;
		border-radius: 0.25rem;
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
