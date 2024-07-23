<script lang="ts">
	import Button from "../../../../lib/components/Button.svelte";
	import { t } from "../../../../lib/utils/i8n";
	import DetailsFields from "../../../../lib/components/details/DetailsFields.svelte";
	import DetailsSection from "../../../../lib/components/details/DetailsSection.svelte";
	import DetailsContent from "../../../../lib/components/details/DetailsContent.svelte";
	import Input from "../../../../lib/components/Input.svelte";
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
	import type { Truck } from "../../../../domain/truck";
	import { TRUCK_ICON_SRC } from "../../../../lib/constants/map";
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
		make: string,
		model: string,
		licensePlate: string,
		personCapacity: number,
		location: GeoJSONFeaturePoint,
	) => void;

	/**
	 * Truck data.
	 * @default null
	 */
	export let truck: Truck | null = null;

	/**
	 * Indicates if form is being submitted.
	 */
	export let isSubmitting: boolean;

	/**
	 * The map which displays the selected truck location.
	 */
	let mapPreview: OlMap;

	/**
	 * The map layer which displays the truck location.
	 */
	const layer = new VectorLayer({
		source: new VectorSource<Feature<Point>>({ features: [] }),
		style: {
			"icon-src": TRUCK_ICON_SRC,
		},
	});

	/**
	 * The select location open modal state.
	 * @default false
	 */
	let openSelectLocation = false;

	/**
	 * The selected truck location coordinate.
	 */
	let selectedCoordinate = truck?.geoJson.geometry.coordinates;

	/**
	 * The location name of the truck.
	 */
	let locationName = truck
		? getLocationName(
				truck.geoJson.properties.wayName,
				truck.geoJson.properties.municipalityName,
			)
		: "";

	/**
	 * The form fields minimum and maximum value lengths.
	 */
	const formFieldsLengths = {
		make: {
			min: 0,
			max: 50,
		},
		model: {
			min: 0,
			max: 50,
		},
		licensePlate: {
			min: 0,
			max: 30,
		},
		personCapacity: {
			min: 1,
			max: 10,
		},
	};

	/**
	 * Error messages of the form fields.
	 */
	let formErrorMessages = {
		make: "",
		model: "",
		licensePlate: "",
		personCapacity: "",
		location: "",
	};

	/**
	 * Adds the truck to the map preview given a coordinate.
	 * @param coordinate Truck coordinate.
	 */
	function addTruckToMap(coordinate: Coordinate) {
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
	 * @param makeValidity Truck make field validity state.
	 * @param modelValidity Truck model field validity state.
	 * @param licensePlateValidity Truck license plate field validity state.
	 * @param personCapacityValidity Truck person capacity field validity state.
	 * @param locationValidity Truck location field validity state.
	 * @param coordinate Truck coordinate.
	 * @returns `true` if form is valid, `false` otherwise.
	 */
	function validateForm(
		makeValidity: ValidityState,
		modelValidity: ValidityState,
		licensePlateValidity: ValidityState,
		personCapacityValidity: ValidityState,
		locationValidity: ValidityState,
		coordinate: number[] | undefined,
	): coordinate is number[] {
		if (makeValidity.valueMissing) {
			formErrorMessages.make = $t("error.valueMissing");
		} else if (makeValidity.tooShort) {
			formErrorMessages.make = $t("error.tooShort", {
				minLength: formFieldsLengths.make.min,
			});
		} else if (makeValidity.tooLong) {
			formErrorMessages.make = $t("error.tooLong", {
				maxLength: formFieldsLengths.make.max,
			});
		} else {
			formErrorMessages.make = "";
		}

		if (modelValidity.valueMissing) {
			formErrorMessages.model = $t("error.valueMissing");
		} else if (modelValidity.tooShort) {
			formErrorMessages.model = $t("error.tooShort", {
				minLength: formFieldsLengths.model.min,
			});
		} else if (modelValidity.tooLong) {
			formErrorMessages.model = $t("error.tooLong", {
				maxLength: formFieldsLengths.model.max,
			});
		} else {
			formErrorMessages.model = "";
		}

		if (licensePlateValidity.valueMissing) {
			formErrorMessages.licensePlate = $t("error.valueMissing");
		} else if (licensePlateValidity.patternMismatch) {
			formErrorMessages.licensePlate = $t("error.patternMismatch");
		} else if (licensePlateValidity.tooShort) {
			formErrorMessages.licensePlate = $t("error.tooShort", {
				minLength: formFieldsLengths.licensePlate.min,
			});
		} else if (licensePlateValidity.tooLong) {
			formErrorMessages.licensePlate = $t("error.tooLong", {
				maxLength: formFieldsLengths.licensePlate.max,
			});
		} else {
			formErrorMessages.licensePlate = "";
		}

		if (
			personCapacityValidity.badInput ||
			personCapacityValidity.stepMismatch
		) {
			formErrorMessages.personCapacity = $t("error.typeMismatch.number");
		} else if (personCapacityValidity.valueMissing) {
			formErrorMessages.personCapacity = $t("error.valueMissing");
		} else if (personCapacityValidity.rangeUnderflow) {
			formErrorMessages.personCapacity = $t("error.rangeUnderflow", {
				min: formFieldsLengths.personCapacity.min,
			});
		} else if (personCapacityValidity.rangeOverflow) {
			formErrorMessages.personCapacity = $t("error.rangeOverflow", {
				max: formFieldsLengths.personCapacity.max,
			});
		} else {
			formErrorMessages.personCapacity = "";
		}

		if (locationValidity.valueMissing) {
			formErrorMessages.location = $t("error.valueMissing");
		} else {
			formErrorMessages.location = "";
		}

		return (
			!formErrorMessages.make &&
			!formErrorMessages.model &&
			!formErrorMessages.licensePlate &&
			!formErrorMessages.personCapacity &&
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

		const make = formData.get("make") ?? "";
		const model = formData.get("model") ?? "";
		const licensePlate = formData.get("licensePlate") ?? "";
		const personCapacity = formData.get("personCapacity") ?? "";
		const location = formData.get("location") ?? "";

		// Check if all fields are strings.
		if (
			typeof make !== "string" ||
			typeof model !== "string" ||
			typeof licensePlate !== "string" ||
			typeof personCapacity !== "string" ||
			typeof location !== "string"
		) {
			return;
		}

		const makeInput = form.elements.namedItem("make") as HTMLInputElement;
		const modelInput = form.elements.namedItem("model") as HTMLInputElement;
		const licensePlateInput = form.elements.namedItem(
			"licensePlate",
		) as HTMLInputElement;
		const personCapacityInput = form.elements.namedItem(
			"personCapacity",
		) as HTMLInputElement;
		const locationInput = form.elements.namedItem(
			"location",
		) as HTMLInputElement;

		// Check if form is valid to prevent making a server request.
		if (
			!validateForm(
				makeInput.validity,
				modelInput.validity,
				licensePlateInput.validity,
				personCapacityInput.validity,
				locationInput.validity,
				selectedCoordinate,
			)
		) {
			return;
		}

		onSave(make, model, licensePlate, Number(personCapacity), {
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
					label={$t("make")}
					error={!!formErrorMessages.make}
					helperText={formErrorMessages.make}
				>
					<Input
						required
						name="make"
						value={truck?.make}
						error={!!formErrorMessages.make}
						placeholder={$t("make.placeholder")}
						minLength={formFieldsLengths.make.min}
						maxLength={formFieldsLengths.make.max}
					/>
				</FormControl>
				<FormControl
					label={$t("model")}
					error={!!formErrorMessages.model}
					helperText={formErrorMessages.model}
				>
					<Input
						required
						name="model"
						value={truck?.model}
						error={!!formErrorMessages.model}
						placeholder={$t("model.placeholder")}
						minLength={formFieldsLengths.model.min}
						maxLength={formFieldsLengths.model.max}
					/>
				</FormControl>
				<FormControl
					label={$t("licensePlate")}
					error={!!formErrorMessages.licensePlate}
					helperText={formErrorMessages.licensePlate}
				>
					<Input
						required
						name="licensePlate"
						pattern={"^[0-9A-Z]{2}-[0-9A-Z]{2}-[0-9A-Z]{2}$"}
						value={truck?.licensePlate}
						error={!!formErrorMessages.licensePlate}
						placeholder={$t("licensePlate.placeholder")}
						minLength={formFieldsLengths.licensePlate.min}
						maxLength={formFieldsLengths.licensePlate.max}
					/>
				</FormControl>
				<FormControl
					label={$t("personCapacity")}
					error={!!formErrorMessages.personCapacity}
					helperText={formErrorMessages.personCapacity}
				>
					<Input
						required
						name="personCapacity"
						value={truck?.personCapacity}
						error={!!formErrorMessages.personCapacity}
						placeholder={$t("personCapacity.placeholder")}
						type="number"
						min={formFieldsLengths.personCapacity.min}
						max={formFieldsLengths.personCapacity.max}
					/>
				</FormControl>
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
		<DetailsSection class="flex-1" label={$t("preview")}>
			<Map
				bind:map={mapPreview}
				onInit={() => {
					const truckCoordinate = truck?.geoJson.geometry.coordinates;
					if (!truckCoordinate) {
						return;
					}

					const mapCoordinate = convertToMapProjection(truckCoordinate);
					addTruckToMap(mapCoordinate);
				}}
			/>
		</DetailsSection>
		<SelectLocation
			open={openSelectLocation}
			coordinate={selectedCoordinate}
			iconSrc={TRUCK_ICON_SRC}
			onOpenChange={open => (openSelectLocation = open)}
			onSave={(coordinate, name) => {
				addTruckToMap(coordinate);
				selectedCoordinate = convertToResourceProjection(coordinate);
				locationName = name;
			}}
		/>
	</DetailsContent>
</form>
