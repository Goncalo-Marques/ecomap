<script lang="ts">
	import { Feature } from "ol";
	import type { Coordinate } from "ol/coordinate";
	import { Point } from "ol/geom";
	import VectorLayer from "ol/layer/Vector";
	import OlMap from "ol/Map";
	import VectorSource from "ol/source/Vector";

	import type { GeoJSONFeaturePoint } from "$domain/geojson";
	import type { Warehouse } from "$domain/warehouse";
	import Button from "$lib/components/Button.svelte";
	import DetailsContent from "$lib/components/details/DetailsContent.svelte";
	import DetailsFields from "$lib/components/details/DetailsFields.svelte";
	import DetailsHeader from "$lib/components/details/DetailsHeader.svelte";
	import DetailsSection from "$lib/components/details/DetailsSection.svelte";
	import FormControl from "$lib/components/FormControl.svelte";
	import Input from "$lib/components/Input.svelte";
	import LocationInput from "$lib/components/LocationInput.svelte";
	import Map from "$lib/components/map/Map.svelte";
	import SelectLocation from "$lib/components/SelectLocation.svelte";
	import { WAREHOUSE_ICON_SRC } from "$lib/constants/map";
	import { t } from "$lib/utils/i8n";
	import { getLocationName } from "$lib/utils/location";
	import {
		convertToMapProjection,
		convertToResourceProjection,
	} from "$lib/utils/map";

	/**
	 * The back route.
	 */
	export let back: string;

	/**
	 * The title in the form.
	 */
	export let title: string;

	/**
	 * Callback fired when save action is triggered.
	 */
	export let onSave: (
		truckCapacity: number,
		location: GeoJSONFeaturePoint,
	) => void;

	/**
	 * Warehouse data.
	 * @default null
	 */
	export let warehouse: Warehouse | null = null;

	/**
	 * Indicates if form is being submitted.
	 */
	export let isSubmitting: boolean;

	/**
	 * The minimum valid capacity for the truck capacity field.
	 */
	const TRUCK_CAPACITY_MIN_VALUE = 0;

	/**
	 * The maximum valid capacity for the truck capacity field.
	 */
	const TRUCK_CAPACITY_MAX_VALUE = 99;

	/**
	 * The map which displays the selected warehouse location.
	 */
	let mapPreview: OlMap;

	/**
	 * The map layer which displays the warehouse location.
	 */
	const layer = new VectorLayer({
		source: new VectorSource<Feature<Point>>({ features: [] }),
		style: {
			"icon-src": WAREHOUSE_ICON_SRC,
		},
	});

	/**
	 * The select location open modal state.
	 * @default false
	 */
	let openSelectLocation = false;

	/**
	 * The selected warehouse location coordinate.
	 */
	let selectedCoordinate = warehouse?.geoJson.geometry.coordinates;

	/**
	 * The location name of the warehouse.
	 */
	let locationName = warehouse
		? getLocationName(
				warehouse.geoJson.properties.wayName,
				warehouse.geoJson.properties.municipalityName,
			)
		: "";

	let truckCapacity = warehouse?.truckCapacity;

	/**
	 * Error messages of the form fields.
	 */
	let formErrorMessages = {
		truckCapacity: "",
		location: "",
	};

	/**
	 * Adds the warehouse to the map preview given a coordinate.
	 * @param coordinate Warehouse coordinate.
	 */
	function addWarehouseToMap(coordinate: Coordinate) {
		const source = layer.getSource();
		if (!source) {
			return;
		}

		source.clear();
		mapPreview.removeLayer(layer);

		const point = new Point(coordinate);
		const feature = new Feature(point);

		source.addFeature(feature);
		mapPreview.addLayer(layer);

		const view = mapPreview.getView();
		view.fit(point);
	}

	/**
	 * Validates the form and sets error messages on the form fields
	 * if they contain any errors.
	 * @param truckCapacityValidity Truck capacity field validity state.
	 * @param locationValidity Location field validity state.
	 * @param coordinate Warehouse coordinate.
	 * @returns `true` if form is valid, `false` otherwise.
	 */
	function validateForm(
		truckCapacityValidity: ValidityState,
		locationValidity: ValidityState,
		coordinate: number[] | undefined,
	): coordinate is number[] {
		if (truckCapacityValidity.badInput || truckCapacityValidity.stepMismatch) {
			formErrorMessages.truckCapacity = $t("error.typeMismatch.number");
		} else if (truckCapacityValidity.valueMissing) {
			formErrorMessages.truckCapacity = $t("error.valueMissing");
		} else if (truckCapacityValidity.rangeUnderflow) {
			formErrorMessages.truckCapacity = $t("error.rangeUnderflow", {
				min: TRUCK_CAPACITY_MIN_VALUE,
			});
		} else if (truckCapacityValidity.rangeOverflow) {
			formErrorMessages.truckCapacity = $t("error.rangeOverflow", {
				max: TRUCK_CAPACITY_MAX_VALUE,
			});
		} else {
			formErrorMessages.truckCapacity = "";
		}

		if (locationValidity.valueMissing) {
			formErrorMessages.location = $t("error.valueMissing");
		} else {
			formErrorMessages.location = "";
		}

		return (
			!formErrorMessages.truckCapacity &&
			!formErrorMessages.location &&
			!!coordinate
		);
	}

	/**
	 * Handles the submit event of the form.
	 * @param e Submit event.
	 */
	function handleSubmit(
		e: Event & { currentTarget: EventTarget & HTMLFormElement },
	) {
		const form = e.currentTarget;
		const formData = new FormData(form);

		const truckCapacity = formData.get("truckCapacity") ?? "";
		const location = formData.get("location") ?? "";

		// Check if truck capacity and location are both strings.
		if (typeof truckCapacity !== "string" || typeof location !== "string") {
			return;
		}

		const truckCapacityInput = form.elements.namedItem(
			"truckCapacity",
		) as HTMLInputElement;
		const locationInput = form.elements.namedItem(
			"location",
		) as HTMLInputElement;

		// Check if form is valid to prevent making a server request.
		if (
			!validateForm(
				truckCapacityInput.validity,
				locationInput.validity,
				selectedCoordinate,
			)
		) {
			return;
		}

		onSave(Number(truckCapacity), {
			type: "Feature",
			geometry: {
				type: "Point",
				coordinates: selectedCoordinate,
			},
			properties: {},
		});
	}
</script>

<form novalidate class="contents" on:submit|preventDefault={handleSubmit}>
	<DetailsHeader href={back} {title}>
		<a href={back} class="contents">
			<Button variant="tertiary">{$t("cancel")}</Button>
		</a>
		<Button type="submit" startIcon="check" disabled={isSubmitting}>
			{$t("save")}
		</Button>
	</DetailsHeader>
	<DetailsContent>
		<DetailsSection label={$t("generalInfo")}>
			<DetailsFields>
				<FormControl
					label={$t("location")}
					error={!!formErrorMessages.location}
					helperText={formErrorMessages.location}
				>
					<LocationInput
						required
						name="location"
						placeholder={$t("location.placeholder")}
						value={locationName}
						error={!!formErrorMessages.location}
						onClick={() => (openSelectLocation = true)}
					/>
				</FormControl>
				<FormControl
					label={$t("truckCapacity")}
					error={!!formErrorMessages.truckCapacity}
					helperText={formErrorMessages.truckCapacity}
				>
					<Input
						required
						name="truckCapacity"
						value={truckCapacity}
						error={!!formErrorMessages.truckCapacity}
						placeholder={$t("truckCapacity.placeholder")}
						type="number"
						min={TRUCK_CAPACITY_MIN_VALUE}
						max={TRUCK_CAPACITY_MAX_VALUE}
					/>
				</FormControl>
			</DetailsFields>
		</DetailsSection>
		<DetailsSection class="flex-1" label={$t("preview")}>
			<Map
				bind:map={mapPreview}
				onInit={() => {
					const warehouseCoordinate = warehouse?.geoJson.geometry.coordinates;
					if (!warehouseCoordinate) {
						return;
					}

					const mapCoordinate = convertToMapProjection(warehouseCoordinate);
					addWarehouseToMap(mapCoordinate);
				}}
			/>
		</DetailsSection>
		<SelectLocation
			open={openSelectLocation}
			coordinate={selectedCoordinate}
			iconSrc={WAREHOUSE_ICON_SRC}
			onOpenChange={open => (openSelectLocation = open)}
			onSave={(coordinate, name) => {
				addWarehouseToMap(coordinate);
				selectedCoordinate = convertToResourceProjection(coordinate);
				locationName = name;
			}}
		/>
	</DetailsContent>
</form>
