<script lang="ts">
	import OlMap from "ol/Map";
	import Modal from "../../../../lib/clients/Modal.svelte";
	import { onMount } from "svelte";
	import Map from "../../../../lib/components/map/Map.svelte";
	import type { Coordinate } from "ol/coordinate";
	import { MapHelper } from "../../../../lib/components/map/mapUtils";
	import VectorLayer from "ol/layer/Vector";
	import { Feature } from "ol";
	import { Point } from "ol/geom";

	export let open: boolean;

	export let onClose: () => void;

	let map: OlMap;

	let mapHelper: MapHelper;

	let selectedCoordinate: Coordinate;

	function markLocation(coordinate: Coordinate) {
		const point = new Point(coordinate);
		const feature = new Feature({ name: "location", geometry: point });
		feature.setId("location");
		map.s;
	}

	onMount(() => {
		mapHelper = new MapHelper(map);
		map.on("click", e => markLocation(e.coordinate));
	});
</script>

<Modal {open} {onClose} title="Selecionar localização">
	<Map bind:map mapId="select-location-map" --height="32rem" --width="60rem" />
</Modal>
