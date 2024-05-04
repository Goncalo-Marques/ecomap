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
	import Input from "../../../../lib/components/Input.svelte";
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
	import { formatContainerCoordinate } from "../utils/location";

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
	 * The map which displays the selected container location.
	 */
	let mapPreview: OlMap;

	/**
	 * The map layer which displays the container location.
	 */
	const layer = new VectorLayer({
		source: new VectorSource<Feature<Point>>({ features: [] }),
		style: {
			"icon-src": "/images/logo.svg",
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
			formErrorMessages.category = $t("error.requiredField");
		} else {
			formErrorMessages.category = "";
		}

		if (!location) {
			formErrorMessages.location = $t("error.requiredField");
		} else {
			formErrorMessages.location = "";
		}
	}

	/**
	 * Handles the submit event of the form.
	 * @param e Submit event.
	 */
	function handleSubmit(e: SubmitEvent) {
		const formData = new FormData(e.currentTarget as HTMLFormElement);
		const category = formData.get("category");
		const location = formData.get("location");

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

<form on:submit|preventDefault={handleSubmit}>
	<DetailsHeader to={back} {title}>
		<Link to={back} style="display:contents">
			<Button variant="tertiary">{$t("cancel")}</Button>
		</Link>
		<Button type="submit" startIcon="check">{$t("save")}</Button>
	</DetailsHeader>
	<DetailsContent>
		<DetailsSection label={$t("generalInfo")}>
			<DetailsFields>
				<FormControl label={$t("containers.category")}>
					<Select name="category" value={container?.category}>
						{#each categoryOptions as category}
							<Option value={category}>
								{$t(`containers.category.${category}`)}
							</Option>
						{/each}
					</Select>
				</FormControl>
				<FormControl
					label={$t("containers.location")}
					error={!!formErrorMessages.location}
					helperText={formErrorMessages.location}
				>
					<Input
						readonly
						name="location"
						value={formatContainerCoordinate(
							selectedCoordinate ?? container?.geoJson.geometry.coordinates,
						)}
						error={!!formErrorMessages.location}
						placeholder={$t("containers.location")}
						endIcon="location_on"
						onClick={() => (openSelectLocation = true)}
					/>
				</FormControl>
			</DetailsFields>
		</DetailsSection>
		<DetailsSection class="container-map-preview" label={$t("preview")}>
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
			onOpenChange={open => (openSelectLocation = open)}
			onSave={coordinate => {
				addContainerToMap(coordinate);
				selectedCoordinate = convertToResourceProjection(coordinate);
			}}
		/>
	</DetailsContent>
</form>

<style>
	:global(.container-map-preview) {
		flex: 1;
	}

	form {
		display: contents;
	}
</style>
