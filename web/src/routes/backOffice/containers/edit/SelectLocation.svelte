<script lang="ts">
	import { onDestroy, onMount } from "svelte";
	import OlMap from "ol/Map";
	import Modal from "../../../../lib/clients/Modal.svelte";
	import Map from "../../../../lib/components/map/Map.svelte";
	import VectorLayer from "ol/layer/Vector";
	import { Feature, MapBrowserEvent } from "ol";
	import { Point } from "ol/geom";
	import VectorSource from "ol/source/Vector";
	import Button from "../../../../lib/components/Button.svelte";
	import type { Coordinate } from "ol/coordinate";
	import { transform } from "ol/proj";

	export let open: boolean;

	export let onCancel: () => void;

	export let onSave: (coordinate: Coordinate) => void;

	export let onOpenChange: (open: boolean) => void;

	export let coordinate: Coordinate | null = null;

	let map: OlMap;

	const containerLocationId = "container";

	const layer = new VectorLayer({
		source: new VectorSource<Feature<Point>>({ features: [] }),
		style: {
			"icon-src": "/images/pin.svg",
		},
	});

	let disabled = true;

	function handleMapClick(e: MapBrowserEvent<any>) {
		addContainerToMap(e.coordinate);
	}

	function handleSave() {
		const source = layer.getSource();
		if (!source) {
			return;
		}

		const feature = source.getFeatureById(containerLocationId);
		if (!feature) {
			return;
		}

		const point = feature.getGeometry();
		if (!point) {
			return;
		}

		const coordinates = point.getCoordinates();
		onSave(coordinates);
	}

	function removeContainerFromMap() {
		const source = layer.getSource();
		if (source) {
			source.clear();
			map.removeLayer(layer);
		}
		disabled = true;
	}

	onMount(() => {
		map.on("click", handleMapClick);
	});

	onDestroy(() => {
		map.un("click", handleMapClick);
	});

	function addContainerToMap(coordinate: Coordinate) {
		const source = layer.getSource();
		if (!source) {
			return;
		}

		source.clear();
		map.removeLayer(layer);

		const point = new Point(coordinate);
		const feature = new Feature(point);
		feature.setId(containerLocationId);

		source.addFeature(feature);
		map.addLayer(layer);

		disabled = false;
	}

	$: if (open && coordinate) {
		addContainerToMap(transform(coordinate, "EPSG:4326", "EPSG:3857"));
	}
</script>

<Modal
	{open}
	{onOpenChange}
	onClickOutside={() => {
		removeContainerFromMap();
	}}
	title="Selecionar localização"
>
	<Map bind:map mapId="select-location-map" --height="32rem" --width="60rem" />
	<svelte:fragment slot="actions">
		<Button
			variant="tertiary"
			onClick={() => {
				onCancel();
			}}>Cancelar</Button
		>
		<Button startIcon="check" {disabled} onClick={handleSave}>Guardar</Button>
	</svelte:fragment>
</Modal>
