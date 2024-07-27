<script lang="ts">
	import type {
		Container,
		ContainerCategory,
	} from "../../../../domain/container";
	import Button from "../../../../lib/components/Button.svelte";
	import { t } from "../../../../lib/utils/i8n";
	import DetailsFields from "../../../../lib/components/details/DetailsFields.svelte";
	import DetailsSection from "../../../../lib/components/details/DetailsSection.svelte";
	import DetailsContent from "../../../../lib/components/details/DetailsContent.svelte";
	import Map from "../../../../lib/components/map/Map.svelte";
	import OlMap from "ol/Map";
	import Select from "../../../../lib/components/Select.svelte";
	import Option from "../../../../lib/components/Option.svelte";
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
	import { categoryOptions } from "../constants/category";
	import {
		convertToMapProjection,
		convertToResourceProjection,
	} from "../../../../lib/utils/map";
	import { isValidContainerCategory } from "../utils/category";
	import { getLocationName } from "../../../../lib/utils/location";
	import { CONTAINER_ICON_SRC } from "../../../../lib/constants/map";
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
	export let onSave: (
		category: ContainerCategory,
		location: GeoJSONFeaturePoint,
	) => void;

	/**
	 * Container data.
	 * @default null
	 */
	export let container: Container | null = null;

	/**
	 * Indicates if form is being submitted.
	 */
	export let isSubmitting: boolean;

	/**
	 * The map which displays the selected container location.
	 */
	let mapPreview: OlMap;

	/**
	 * The map layer which displays the container location.
	 */
	const layer = new VectorLayer({
		source: new VectorSource<Feature<Point>>({ features: [] }),
		style: {
			"icon-src": CONTAINER_ICON_SRC,
		},
	});

	/**
	 * The select location open modal state.
	 * @default false
	 */
	let openSelectLocation = false;

	/**
	 * The selected container location coordinate.
	 */
	let selectedCoordinate = container?.geoJson.geometry.coordinates;

	/**
	 * The location name of the container.
	 */
	let locationName = container
		? getLocationName(
				container.geoJson.properties.wayName,
				container.geoJson.properties.municipalityName,
			)
		: "";

	/**
	 * Error messages of the form fields.
	 */
	let formErrorMessages = {
		category: "",
		location: "",
	};

	/**
	 * Adds the container to the map preview given a coordinate.
	 * @param coordinate Container coordinate.
	 */
	function addContainerToMap(coordinate: Coordinate) {
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
	 * @param category Category field value.
	 * @param location Location field value.
	 */
	function validateForm(category: string, location: string) {
		if (!category) {
			formErrorMessages.category = $t("error.valueMissing");
		} else {
			formErrorMessages.category = "";
		}

		if (!location) {
			formErrorMessages.location = $t("error.valueMissing");
		} else {
			formErrorMessages.location = "";
		}
	}

	/**
	 * Handles the submit event of the form.
	 * @param e Submit event.
	 */
	function handleSubmit(
		e: Event & { currentTarget: EventTarget & HTMLFormElement },
	) {
		const formData = new FormData(e.currentTarget);
		const category = formData.get("category") ?? "";
		const location = formData.get("location") ?? "";

		// Check if category and location are both strings.
		if (typeof category !== "string" || typeof location !== "string") {
			return;
		}

		validateForm(category, location);

		// Check if fields are not filled to prevent making a server request.
		if (!category || !location || !selectedCoordinate) {
			return;
		}

		if (!isValidContainerCategory(category)) {
			return;
		}

		onSave(category, {
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
	<DetailsHeader to={back} {title}>
		<Link to={back} class="contents">
			<Button variant="tertiary">{$t("cancel")}</Button>
		</Link>
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
						value={locationName}
						error={!!formErrorMessages.location}
						placeholder={$t("location.placeholder")}
						onClick={() => (openSelectLocation = true)}
					/>
				</FormControl>
				<FormControl
					label={$t("containers.category")}
					error={!!formErrorMessages.category}
					helperText={formErrorMessages.category}
				>
					<Select
						required
						name="category"
						error={!!formErrorMessages.category}
						placeholder={$t("containers.category.placeholder")}
						value={container?.category}
					>
						{#each categoryOptions as category}
							<Option value={category}>
								{$t(`containers.category.${category}`)}
							</Option>
						{/each}
					</Select>
				</FormControl>
			</DetailsFields>
		</DetailsSection>
		<DetailsSection class="flex-1" label={$t("preview")}>
			<Map
				bind:map={mapPreview}
				onInit={() => {
					const containerCoordinate = container?.geoJson.geometry.coordinates;
					if (!containerCoordinate) {
						return;
					}

					const mapCoordinate = convertToMapProjection(containerCoordinate);
					addContainerToMap(mapCoordinate);
				}}
			/>
		</DetailsSection>
		<SelectLocation
			open={openSelectLocation}
			coordinate={selectedCoordinate}
			iconSrc={CONTAINER_ICON_SRC}
			onOpenChange={open => (openSelectLocation = open)}
			onSave={(coordinate, name) => {
				addContainerToMap(coordinate);
				selectedCoordinate = convertToResourceProjection(coordinate);
				locationName = name;
			}}
		/>
	</DetailsContent>
</form>
