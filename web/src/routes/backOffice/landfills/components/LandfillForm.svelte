<script lang="ts">
	import Button from "../../../../lib/components/Button.svelte";
	import { t } from "../../../../lib/utils/i8n";
	import DetailsFields from "../../../../lib/components/details/DetailsFields.svelte";
	import DetailsSection from "../../../../lib/components/details/DetailsSection.svelte";
	import DetailsContent from "../../../../lib/components/details/DetailsContent.svelte";
	import Map from "../../../../lib/components/map/Map.svelte";
	import OlMap from "ol/Map";
	import FormControl from "../../../../lib/components/FormControl.svelte";
	import SelectLocation from "../../../../lib/components/SelectLocation.svelte";
	import VectorLayer from "ol/layer/Vector";
	import VectorSource from "ol/source/Vector";
	import { Point } from "ol/geom";
	import { Feature } from "ol";
	import type { Coordinate } from "ol/coordinate";
	import { Link } from "svelte-routing";
	import DetailsHeader from "../../../../lib/components/details/DetailsHeader.svelte";
	import type { GeoJSONFeaturePoint } from "../../../../domain/geojson";
	import {
		convertToMapProjection,
		convertToResourceProjection,
	} from "../../../../lib/utils/map";
	import { getLocationName } from "../../../../lib/utils/location";
	import type { Landfill } from "../../../../domain/landfill";
	import { LANDFILL_ICON_SRC } from "../../../../lib/constants/map";
	import LocationInput from "../../../../lib/components/LocationInput.svelte";

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
	export let onSave: (location: GeoJSONFeaturePoint) => void;

	/**
	 * Landfill data.
	 * @default null
	 */
	export let landfill: Landfill | null = null;

	/**
	 * The map which displays the selected landfill location.
	 */
	let mapPreview: OlMap;

	/**
	 * The map layer which displays the landfill location.
	 */
	const layer = new VectorLayer({
		source: new VectorSource<Feature<Point>>({ features: [] }),
		style: {
			"icon-src": LANDFILL_ICON_SRC,
		},
	});

	/**
	 * The select location open modal state.
	 * @default false
	 */
	let openSelectLocation = false;

	/**
	 * The selected landfill location coordinate.
	 */
	let selectedCoordinate = landfill?.geoJson.geometry.coordinates;

	/**
	 * The location name of the landfill.
	 */
	let locationName = landfill
		? getLocationName(
				landfill.geoJson.properties.wayName,
				landfill.geoJson.properties.municipalityName,
			)
		: "";

	/**
	 * Error messages of the form fields.
	 */
	let formErrorMessages = {
		location: "",
	};

	/**
	 * Adds the landfill to the map preview given a coordinate.
	 * @param coordinate Landfill coordinate.
	 */
	function addLandfillToMap(coordinate: Coordinate) {
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
	 * @param locationValidity Location field validity state.
	 * @param coordinate Landfill coordinate.
	 * @returns `true` if form is valid, `false` otherwise.
	 */
	function validateForm(
		locationValidity: ValidityState,
		coordinate: number[] | undefined,
	): coordinate is number[] {
		if (locationValidity.valueMissing) {
			formErrorMessages.location = $t("error.valueMissing");
		} else {
			formErrorMessages.location = "";
		}

		return !formErrorMessages.location && !!coordinate;
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

		const location = formData.get("location") ?? "";

		// Check if location is a string.
		if (typeof location !== "string") {
			return;
		}

		const locationInput = form.elements.namedItem(
			"location",
		) as HTMLInputElement;

		// Check if form is valid to prevent making a server request.
		if (!validateForm(locationInput.validity, selectedCoordinate)) {
			return;
		}

		onSave({
			type: "Feature",
			geometry: {
				type: "Point",
				coordinates: selectedCoordinate,
			},
			properties: {},
		});
	}
</script>

<form novalidate on:submit|preventDefault={handleSubmit}>
	<DetailsHeader to={back} {title}>
		<Link to={back} style="display:contents">
			<Button variant="tertiary">{$t("cancel")}</Button>
		</Link>
		<Button type="submit" startIcon="check">{$t("save")}</Button>
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
			</DetailsFields>
		</DetailsSection>
		<DetailsSection class="landfill-map-preview" label={$t("preview")}>
			<Map
				bind:map={mapPreview}
				onInit={() => {
					const landfillCoordinate = landfill?.geoJson.geometry.coordinates;
					if (!landfillCoordinate) {
						return;
					}

					const mapCoordinate = convertToMapProjection(landfillCoordinate);
					addLandfillToMap(mapCoordinate);
				}}
			/>
		</DetailsSection>
		<SelectLocation
			open={openSelectLocation}
			coordinate={selectedCoordinate}
			iconSrc={LANDFILL_ICON_SRC}
			onOpenChange={open => (openSelectLocation = open)}
			onSave={(coordinate, name) => {
				addLandfillToMap(coordinate);
				selectedCoordinate = convertToResourceProjection(coordinate);
				locationName = name;
			}}
		/>
	</DetailsContent>
</form>

<style>
	:global(.landfill-map-preview) {
		flex: 1;
	}

	form {
		display: contents;
	}
</style>
