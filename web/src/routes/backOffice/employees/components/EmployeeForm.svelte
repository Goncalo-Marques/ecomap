<script lang="ts">
	import Button from "../../../../lib/components/Button.svelte";
	import { t } from "../../../../lib/utils/i8n";
	import DetailsFields from "../../../../lib/components/details/DetailsFields.svelte";
	import DetailsSection from "../../../../lib/components/details/DetailsSection.svelte";
	import DetailsContent from "../../../../lib/components/details/DetailsContent.svelte";
	import Input from "../../../../lib/components/Input.svelte";
	import FormControl from "../../../../lib/components/FormControl.svelte";
	import { Link } from "svelte-routing";
	import DetailsHeader from "../../../../lib/components/details/DetailsHeader.svelte";
	import type { GeoJSONFeaturePoint } from "../../../../domain/geojson";
	import { getLocationName } from "../../../../lib/utils/location";
	import LocationInput from "../../../../lib/components/LocationInput.svelte";
	import type { Employee } from "../../../../domain/employees";
	import SelectLocation from "../../../../lib/components/SelectLocation.svelte";
	import { convertToResourceProjection } from "../../../../lib/utils/map";

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
		username: string,
		firstName: string,
		lastName: string,
		dateOfBirth: string,
		phoneNumber: string,
		location: GeoJSONFeaturePoint,
		scheduleStart: string,
		scheduleEnd: string,
	) => void;

	/**
	 * Employee data.
	 * @default null
	 */
	export let employee: Employee | null = null;

	/**
	 * The select location open modal state.
	 * @default false
	 */
	let openSelectLocation = false;

	/**
	 * The selected employee location coordinate.
	 */
	let selectedCoordinate = employee?.geoJson.geometry.coordinates;

	/**
	 * The location name of the employee.
	 */
	let locationName = employee
		? getLocationName(
				employee.geoJson.properties.wayName,
				employee.geoJson.properties.municipalityName,
			)
		: "";

	/**
	 * The form fields minimum and maximum value lengths.
	 */
	const formFieldsLengths = {
		username: {
			min: 5,
			max: 50,
		},
		firstName: {
			min: 0,
			max: 30,
		},
		lastName: {
			min: 0,
			max: 30,
		},
		dateOfBirth: {
			min: 10,
			max: 10,
		},
		phoneNumber: {
			min: 1,
			max: 20,
		},
	};

	/**
	 * Error messages of the form fields.
	 */
	let formErrorMessages = {
		username: "",
		firstName: "",
		lastName: "",
		dateOfBirth: "",
		phoneNumber: "",
		location: "",
		scheduleStart: "",
		scheduleEnd: "",
	};

	/**
	 * Validates form and sets error messages on the form fields
	 * if they contain errors.
	 * @param usernameValidity
	 * @param firstNameValidity
	 * @param lastNameValidity
	 * @param dateOfBirthValidity
	 * @param phoneNumberValidity
	 * @param locationValidity
	 * @param coordinate
	 */
	function validateForm(
		usernameValidity: ValidityState,
		firstNameValidity: ValidityState,
		lastNameValidity: ValidityState,
		dateOfBirthValidity: ValidityState,
		phoneNumberValidity: ValidityState,
		locationValidity: ValidityState,
		coordinate: number[] | undefined,
	): coordinate is number[] {
		// Username Validation.
		if (usernameValidity.valueMissing) {
			formErrorMessages.username = $t("error.valueMissing");
		} else if (usernameValidity.tooShort) {
			formErrorMessages.username = $t("error.tooShort", {
				minLength: formFieldsLengths.username.min,
			});
		} else if (usernameValidity.tooLong) {
			formErrorMessages.username = $t("error.tooLong", {
				maxLength: formFieldsLengths.username.max,
			});
		} else {
			formErrorMessages.username = "";
		}

		//FirstName Validation.
		if (firstNameValidity.valueMissing) {
			formErrorMessages.firstName = $t("error.valueMissing");
		} else if (firstNameValidity.tooShort) {
			formErrorMessages.firstName = $t("error.tooShort", {
				minLength: formFieldsLengths.firstName.min,
			});
		} else if (firstNameValidity.tooLong) {
			formErrorMessages.firstName = $t("error.tooLong", {
				maxLength: formFieldsLengths.firstName.max,
			});
		} else {
			formErrorMessages.firstName = "";
		}

		// LastName Validation.
		if (lastNameValidity.valueMissing) {
			formErrorMessages.lastName = $t("error.valueMissing");
		} else if (lastNameValidity.patternMismatch) {
			formErrorMessages.lastName = $t("error.patternMismatch");
		} else if (lastNameValidity.tooShort) {
			formErrorMessages.lastName = $t("error.tooShort", {
				minLength: formFieldsLengths.lastName.min,
			});
		} else if (lastNameValidity.tooLong) {
			formErrorMessages.lastName = $t("error.tooLong", {
				maxLength: formFieldsLengths.lastName.max,
			});
		} else {
			formErrorMessages.lastName = "";
		}

		// DateOfBirth Validation.
		if (dateOfBirthValidity.valueMissing) {
			formErrorMessages.dateOfBirth = $t("error.valueMissing");
		} else if (dateOfBirthValidity.patternMismatch) {
			formErrorMessages.dateOfBirth = $t("error.patternMismatch");
		} else if (dateOfBirthValidity.tooShort) {
			formErrorMessages.dateOfBirth = $t("error.tooShort", {
				minLength: formFieldsLengths.dateOfBirth.min,
			});
		} else if (dateOfBirthValidity.tooLong) {
			formErrorMessages.dateOfBirth = $t("error.tooLong", {
				maxLength: formFieldsLengths.dateOfBirth.max,
			});
		} else {
			formErrorMessages.dateOfBirth = "";
		}

		// PhoneNumber Validation.
		if (phoneNumberValidity.valueMissing) {
			formErrorMessages.phoneNumber = $t("error.valueMissing");
		} else if (phoneNumberValidity.patternMismatch) {
			formErrorMessages.phoneNumber = $t("error.patternMismatch");
		} else if (phoneNumberValidity.tooShort) {
			formErrorMessages.phoneNumber = $t("error.tooShort", {
				minLength: formFieldsLengths.phoneNumber.min,
			});
		} else if (phoneNumberValidity.tooLong) {
			formErrorMessages.phoneNumber = $t("error.tooLong", {
				maxLength: formFieldsLengths.phoneNumber.max,
			});
		} else {
			formErrorMessages.phoneNumber = "";
		}

		// Location Validation.
		if (locationValidity.valueMissing) {
			formErrorMessages.location = $t("error.valueMissing");
		} else {
			formErrorMessages.location = "";
		}

		return (
			!formErrorMessages.username &&
			!formErrorMessages.firstName &&
			!formErrorMessages.lastName &&
			!formErrorMessages.dateOfBirth &&
			!formErrorMessages.phoneNumber &&
			!formErrorMessages.location &&
			!!coordinate
		);
	}

	/**
	 * Handles the submit event of the form.
	 * @param e Submit event.
	 */
	function handleSubmit(e: SubmitEvent) {
		const form = e.currentTarget as HTMLFormElement;
		const formData = new FormData(form);

		const username = formData.get("username") ?? "";
		const firstName = formData.get("firstName") ?? "";
		const lastName = formData.get("lastName") ?? "";
		const dateOfBirth = formData.get("dateOfBirth") ?? "";
		const phoneNumber = formData.get("phoneNumber") ?? "";
		const location = formData.get("location") ?? "";
		const scheduleStart = formData.get("scheduleStart") ?? "";
		const scheduleEnd = formData.get("scheduleEnd") ?? "";

		// Check if all fields are strings.
		if (
			typeof username !== "string" ||
			typeof firstName !== "string" ||
			typeof lastName !== "string" ||
			typeof dateOfBirth !== "string" ||
			typeof phoneNumber !== "string" ||
			typeof location !== "string" ||
			typeof scheduleStart !== "string" ||
			typeof scheduleEnd !== "string"
		) {
			return;
		}

		const usernameInput = form.elements.namedItem(
			"username",
		) as HTMLInputElement;

		const firstNameInput = form.elements.namedItem(
			"firstName",
		) as HTMLInputElement;
		const lastNameInput = form.elements.namedItem(
			"lastName",
		) as HTMLInputElement;
		const dateOfBirthInput = form.elements.namedItem(
			"dateOfBirth",
		) as HTMLInputElement;
		const phoneNumberInput = form.elements.namedItem(
			"phoneNumber",
		) as HTMLInputElement;
		const locationInput = form.elements.namedItem(
			"location",
		) as HTMLInputElement;

		// Check if form is valid to prevent making a server request.
		if (
			!validateForm(
				usernameInput.validity,
				firstNameInput.validity,
				lastNameInput.validity,
				dateOfBirthInput.validity,
				phoneNumberInput.validity,
				locationInput.validity,
				selectedCoordinate,
			)
		) {
			return;
		}

		onSave(
			username,
			firstName,
			lastName,
			dateOfBirth,
			phoneNumber,
			{
				type: "Feature",
				geometry: {
					type: "Point",
					coordinates: selectedCoordinate,
				},
				properties: {},
			},
			scheduleStart,
			scheduleEnd,
		);
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
		<DetailsSection label={$t("personalInfo")}>
			<DetailsFields>
				<!-- FirstName -->
				<FormControl
					label={$t("employees.firstName")}
					error={!!formErrorMessages.firstName}
					helperText={formErrorMessages.firstName}
				>
					<Input
						required
						name="firstName"
						value={employee?.firstName}
						error={!!formErrorMessages.firstName}
						placeholder={$t("employees.firstName.placeholder")}
						minLength={formFieldsLengths.firstName.min}
						maxLength={formFieldsLengths.firstName.max}
					/>
				</FormControl>

				<!-- LastName -->
				<FormControl
					label={$t("employees.lastName")}
					error={!!formErrorMessages.lastName}
					helperText={formErrorMessages.lastName}
				>
					<Input
						required
						name="lastName"
						value={employee?.lastName}
						error={!!formErrorMessages.lastName}
						placeholder={$t("employees.firstName.placeholder")}
						minLength={formFieldsLengths.firstName.min}
						maxLength={formFieldsLengths.firstName.max}
					/>
				</FormControl>

				<!-- Username -->
				<FormControl
					label={$t("employees.username")}
					error={!!formErrorMessages.username}
					helperText={formErrorMessages.username}
				>
					<Input
						required
						name="username"
						value={employee?.username}
						error={!!formErrorMessages.username}
						placeholder={$t("employees.username.placeholder")}
						minLength={formFieldsLengths.username.min}
						maxLength={formFieldsLengths.username.max}
					/>
				</FormControl>

				<!-- Birthdate -->
				<FormControl
					label={$t("employees.birthdate")}
					error={!!formErrorMessages.dateOfBirth}
					helperText={formErrorMessages.dateOfBirth}
				>
					<Input
						required
						name="dateOfBirth"
						value={employee?.dateOfBirth}
						error={!!formErrorMessages.dateOfBirth}
						placeholder={$t("employees.birthdate.placeholder")}
						type="date"
						minLength={formFieldsLengths.dateOfBirth.min}
						maxLength={formFieldsLengths.dateOfBirth.max}
					/>
				</FormControl>

				<!-- PhoneNumber -->
				<FormControl
					label={$t("employees.phone")}
					error={!!formErrorMessages.phoneNumber}
					helperText={formErrorMessages.phoneNumber}
				>
					<Input
						required
						name="phoneNumber"
						value={employee?.phoneNumber}
						error={!!formErrorMessages.phoneNumber}
						placeholder={$t("employees.phone.placeholder")}
						minLength={formFieldsLengths.phoneNumber.min}
						maxLength={formFieldsLengths.phoneNumber.max}
					/>
				</FormControl>

				<!-- Location -->
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
		<SelectLocation
			open={openSelectLocation}
			coordinate={selectedCoordinate}
			onOpenChange={open => (openSelectLocation = open)}
			onSave={(coordinate, name) => {
				selectedCoordinate = convertToResourceProjection(coordinate);
				locationName = name;
			}}
		/>
		<DetailsSection label={$t("work")}>
			<DetailsFields>
				<!-- scheduleStart -->
				<FormControl
					label={$t("employees.scheduleStart")}
					error={!!formErrorMessages.scheduleStart}
					helperText={formErrorMessages.scheduleStart}
				>
					<Input
						required
						name="scheduleStart"
						value={employee?.scheduleStart}
						error={!!formErrorMessages.scheduleStart}
						placeholder={$t("employees.scheduleStart.placeholder")}
						type="time"
					/>
				</FormControl>

				<!-- scheduleEnd -->
				<FormControl
					label={$t("employees.scheduleEnd")}
					error={!!formErrorMessages.scheduleEnd}
					helperText={formErrorMessages.scheduleEnd}
				>
					<Input
						required
						name="scheduleEnd"
						value={employee?.scheduleEnd}
						error={!!formErrorMessages.scheduleEnd}
						placeholder={$t("employees.scheduleEnd.placeholder")}
						type="time"
					/>
				</FormControl>
			</DetailsFields>
		</DetailsSection>
	</DetailsContent>
</form>

<style>
	form {
		display: contents;
	}
</style>
