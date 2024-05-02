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
	import VectorSource from "ol/source/Vector";
	import Button from "../../../../lib/components/Button.svelte";

	export let open: boolean;

	export let onCancel: () => void;

	export let onSave: (coordinates: number[]) => void;

	let map: OlMap;

	let selectedCoordinates: number[];

	const layer = new VectorLayer({
		source: new VectorSource<Feature<Point>>({ features: [] }),
		style: {
			"icon-src": "/images/pin.svg",
		},
	});

	function markLocation(coordinate: Coordinate) {
		const source = layer.getSource();

		if (!source) {
			return;
		}

		source.clear();
		map.removeLayer(layer);

		const point = new Point(coordinate);
		const feature = new Feature(point);

		source.addFeature(feature);
		map.addLayer(layer);
	}

	onMount(() => {
		map.on("click", e => markLocation(e.coordinate));
	});
</script>

<Modal {open} {onClose} title="Selecionar localização">
	<Map bind:map mapId="select-location-map" --height="32rem" --width="60rem" />
	<svelte:fragment slot="actions">
		<Button variant="tertiary" onClick={onCancel}>Cancelar</Button>
		<Button
			startIcon="check"
			disabled={!!selectedCoordinates}
			onClick={() => {
				onSave(selectedCoordinates);
			}}
		>
			Guardar
		</Button>
	</svelte:fragment>
</Modal>
