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
	 * A space-separated list of the classes of the element.
	 *
	 * @default ""
	 */
	let className: string = "";
	export { className as class };

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

<div
	id={mapId}
	class={`relative h-[var(--height,100%)] max-h-[var(--max-height,100%)] w-[var(--width,100%)] max-w-[var(--max-width,100%)] ${className}`}
>
	{#if showLayers && layers.length}
		<section
			class="absolute bottom-12 left-12 z-50 max-h-[37.5rem] min-w-64 overflow-auto rounded bg-white p-[0.625rem]"
		>
			<header class="flex items-center gap-1">
				<Icon name="layers" />
				<h1 class="font-semibold">{$t("layers")}</h1>
			</header>
			<div class="mt-3 flex flex-col gap-2">
				{#each layers as layer}
					{#if layer.get(nameLayerKey) != mapLayerName}
						<LayerItem {layer} />
					{/if}
				{/each}
			</div>
		</section>
	{/if}
</div>
