<script lang="ts">
	import { onDestroy, onMount } from "svelte";
	import OlMap from "ol/Map";
	import VectorLayer from "ol/layer/Vector";
	import { Feature, MapBrowserEvent } from "ol";
	import { Point } from "ol/geom";
	import VectorSource from "ol/source/Vector";
	import Modal from "../../../../lib/clients/Modal.svelte";
	import Button from "../../../../lib/components/Button.svelte";
	import Map from "../../../../lib/components/map/Map.svelte";
	import { t } from "../../../../lib/utils/i8n";
	import type { Container } from "../../../../domain/container";
	import { iconStyle, selectedIconStyle } from "./constants";

	/**
	 * Indicates if the modal is open.
	 */
	export let open: boolean;

	/**
	 * Callback fired when open state of the modal changes.
	 * @param open New open state modal.
	 */
	export let onOpenChange: (open: boolean) => void;

	/**
	 * Callback fired when save action is triggered.
	 * @param addedContainers Container added.
	 * @param deletedContainers Container removed.
	 */
	export let onSave: (
		addedContainers: Container[],
		deletedContainers: Container[],
	) => void;

	/**
	 * Added containers.
	 * @default null
	 */
	export let addedContainers: Container[] | null = null;

	/**
	 * Deleted containers.
	 * @default null
	 */
	export let deletedContainers: Container[] | null = null;

	/**
	 * Callback fired when cancel action is triggered.
	 * @default null
	 */
	export let onCancel: (() => void) | null = null;

	/**
	 * Open Layers map.
	 */
	let map: OlMap;

	/**
	 * ID of the selected location.
	 */
	const featureId = "location";

	/**
	 * The map layer which displays the selected location.
	 */
	const layer = new VectorLayer({
		source: new VectorSource<Feature<Point>>({ features: [] }),
		style(feature) {
			if (feature.get("selected")) {
				return selectedIconStyle;
			}

			return iconStyle;
		},
	});

	/**
	 * Indicates if save action is disabled.
	 * @default true
	 */
	let isSaveActionDisabled = true;

	/**
	 * Handles the click event on the map.
	 * @param e Click event.
	 */
	function handleMapClick(e: MapBrowserEvent<UIEvent>) {}

	/**
	 * Handles the save action of the modal.
	 */
	async function handleSave() {
		onOpenChange(false);
	}

	/**
	 * Handles cancel action.
	 */
	function handleCancel() {
		onOpenChange(false);
		onCancel?.();
	}

	onMount(() => {
		map.on("click", handleMapClick);
	});

	onDestroy(() => {
		map.un("click", handleMapClick);
	});

	// Adds/removes selected location when modal is opened.
	$: if (open) {
	}
</script>

<Modal {open} {onOpenChange} title={$t("routes.selectContainers")}>
	<Map
		bind:map
		mapId="select-containers-map"
		--height="32rem"
		--width="60rem"
	/>
	<svelte:fragment slot="actions">
		<Button variant="tertiary" onClick={handleCancel}>{$t("cancel")}</Button>
		<Button
			startIcon="check"
			disabled={isSaveActionDisabled}
			onClick={handleSave}
		>
			{$t("save")}
		</Button>
	</svelte:fragment>
</Modal>
