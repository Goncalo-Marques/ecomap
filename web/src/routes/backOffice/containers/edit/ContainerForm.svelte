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
	import { categoryOptions } from "../../../../lib/constants/container";
	import FormControl from "../../../../lib/components/FormControl.svelte";
	import SelectLocation from "./SelectLocation.svelte";
	import VectorLayer from "ol/layer/Vector";
	import VectorSource from "ol/source/Vector";
	import { Point } from "ol/geom";
	import { Feature } from "ol";
	import type { Coordinate } from "ol/coordinate";
	import { transform } from "ol/proj";
	import { Link } from "svelte-routing";
	import DetailsHeader from "../../../../lib/components/details/DetailsHeader.svelte";
	import type { GeoJSONFeaturePoint } from "../../../../domain/geojson";

	/**
	 * Container.
	 * @default null
	 */
	export let container: Container | null = null;

	export let to: string;

	export let title: string;

	export let onSave: (
		category: ContainerCategory,
		location: GeoJSONFeaturePoint,
	) => void;

	let map: OlMap;

	const layer = new VectorLayer({
		source: new VectorSource<Feature<Point>>({ features: [] }),
		style: {
			"icon-src": "/images/logo.svg",
		},
	});

	let openSelectLocation = false;

	let selectedCoordinate = container?.geoJson.geometry.coordinates;

	function addContainerToMap(coordinate: Coordinate) {
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

		const view = map.getView();
		view.fit(point, { maxZoom: 18 });
	}

	function getContainerCoordinate(coordinate: Coordinate | null | undefined) {
		if (!coordinate) {
			return null;
		}

		const [lon, lat] = coordinate;

		return `${lat}, ${lon}`;
	}

	/**
	 * Error messages of the form fields.
	 */
	let formErrorMessages = {
		category: "",
		location: "",
	};

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

		// Check if either fields are not filled to prevent making a server request.
		if (!category || !location || !selectedCoordinate) {
			return;
		}

		switch (category) {
			case "general":
			case "paper":
			case "plastic":
			case "metal":
			case "glass":
			case "organic":
			case "hazardous":
				onSave(category, {
					type: "Feature",
					geometry: {
						type: "Point",
						coordinates: selectedCoordinate,
					},
					properties: {},
				});
		}
	}
</script>

<form on:submit|preventDefault={handleSubmit}>
	<DetailsHeader {to} {title}>
		<Link {to} style="display:contents">
			<Button variant="tertiary">Cancelar</Button>
		</Link>
		<Button type="submit" startIcon="check">Guardar</Button>
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
						value={getContainerCoordinate(
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
		<DetailsSection class="container-map-preview" label={"Pré-visualização"}>
			<Map
				bind:map
				onInit={() => {
					if (!container) {
						return;
					}
					addContainerToMap(
						transform(
							container.geoJson.geometry.coordinates,
							"EPSG:4326",
							"EPSG:3857",
						),
					);
				}}
			/>
		</DetailsSection>
		<SelectLocation
			open={openSelectLocation}
			coordinate={selectedCoordinate}
			onOpenChange={open => {
				openSelectLocation = open;
			}}
			onCancel={() => {
				openSelectLocation = false;
			}}
			onSave={coordinates => {
				addContainerToMap(coordinates);
				selectedCoordinate = transform(coordinates, "EPSG:3857", "EPSG:4326");
				openSelectLocation = false;
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
