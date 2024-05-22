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
	import type { Employee, EmployeeRoles } from "../../../../domain/employees";
	import SelectLocation from "../../../../lib/components/SelectLocation.svelte";
	import { convertToResourceProjection } from "../../../../lib/utils/map";
	import { isValidEmployeeRole } from "../utils/employee";
	import Select from "../../../../lib/components/Select.svelte";
	import { rolesOptions } from "../constants/roles";
	import Option from "../../../../lib/components/Option.svelte";

	/**
	 * The back route.
	 */
	export let back: string;

	/**
	 * The title in the form.
	 */
	export let title: string;

	/**
	 * Set form as a create form, to create a new employee.\
	 * @default false
	 */
	export let createForm: boolean = false;

	/**
	 * Callback fired when save action is triggered.
	 */
	export let onSave: onSaveType | onSaveCreateType;

	/**
	 * Type of callback to create new employee.
	 */
	type onSaveCreateType = (
		username: string,
		password: string,
		firstName: string,
		lastName: string,
		role: EmployeeRoles,
		dateOfBirth: string,
		phoneNumber: string,
		location: GeoJSONFeaturePoint,
		scheduleStart: string,
		scheduleEnd: string,
	) => void;

	/**
	 * Type of callback to update employee.
	 */
	type onSaveType = (
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
			min: 0,
			max: 30,
		},
		password: {
			min: 0,
			max: 30,
		},
		firstName: {
			min: 0,
			max: 30,
		},
		lastName: {
			min: 0,
			max: 30,
		},
		phoneNumber: {
			min: 0,
			max: 11,
		},
	};

	/**
	 * Error messages of the form fields.
	 */
	let formErrorMessages = {
		username: "",
		password: "",
		firstName: "",
		lastName: "",
		role: "",
		dateOfBirth: "",
		phoneNumber: "",
		location: "",
		scheduleStart: "",
		scheduleEnd: "",
	};

	/**
	 * Validates form and sets error messages on the form fields if they contain errors.
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
		passwordValidity: ValidityState,
		firstNameValidity: ValidityState,
		lastNameValidity: ValidityState,
		roleValidity: ValidityState,
		dateOfBirthValidity: ValidityState,
		phoneNumberValidity: ValidityState,
		locationValidity: ValidityState,
		scheduleStart: ValidityState,
		scheduleEnd: ValidityState,
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

		// Password Validation.
		if (passwordValidity.valueMissing) {
			formErrorMessages.password = $t("error.valueMissing");
		} else if (passwordValidity.patternMismatch) {
			formErrorMessages.phoneNumber = $t("error.patternMismatch");
		} else if (passwordValidity.tooShort) {
			formErrorMessages.password = $t("error.tooShort", {
				minLength: formFieldsLengths.password.min,
			});
		} else if (passwordValidity.tooLong) {
			formErrorMessages.password = $t("error.tooLong", {
				maxLength: formFieldsLengths.password.max,
			});
		} else {
			formErrorMessages.password = "";
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

		//  Role Validation.
		if (roleValidity.valueMissing) {
			formErrorMessages.role = $t("error.valueMissing");
		} else {
			formErrorMessages.role = "";
		}

		// DateOfBirth Validation.
		if (dateOfBirthValidity.valueMissing) {
			formErrorMessages.dateOfBirth = $t("error.valueMissing");
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

		// ScheduleStart Validation.
		if (scheduleStart.valueMissing) {
			formErrorMessages.scheduleStart = $t("error.valueMissing");
		} else {
			formErrorMessages.scheduleStart = "";
		}

		// ScheduleEnd Validation.
		if (scheduleEnd.valueMissing) {
			formErrorMessages.scheduleEnd = $t("error.valueMissing");
		} else {
			formErrorMessages.scheduleEnd = "";
		}

		// Location Validation.
		if (locationValidity.valueMissing) {
			formErrorMessages.location = $t("error.valueMissing");
		} else {
			formErrorMessages.location = "";
		}

		return (
			!formErrorMessages.username &&
			!formErrorMessages.password &&
			!formErrorMessages.firstName &&
			!formErrorMessages.lastName &&
			!formErrorMessages.role &&
			!formErrorMessages.dateOfBirth &&
			!formErrorMessages.phoneNumber &&
			!formErrorMessages.scheduleStart &&
			!formErrorMessages.scheduleEnd &&
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
		const password = formData.get("password") ?? "";
		const firstName = formData.get("firstName") ?? "";
		const lastName = formData.get("lastName") ?? "";
		const role = formData.get("role") ?? "";
		const dateOfBirth = formData.get("dateOfBirth") ?? "";
		const phoneNumber = formData.get("phoneNumber") ?? "";
		const location = formData.get("location") ?? "";
		const scheduleStart = formData.get("scheduleStart") ?? "";
		const scheduleEnd = formData.get("scheduleEnd") ?? "";

		// Check if all fields are strings.
		if (
			typeof username !== "string" ||
			typeof password !== "string" ||
			typeof firstName !== "string" ||
			typeof lastName !== "string" ||
			typeof role !== "string" ||
			typeof dateOfBirth !== "string" ||
			typeof phoneNumber !== "string" ||
			typeof location !== "string" ||
			typeof scheduleStart !== "string" ||
			typeof scheduleEnd !== "string"
		) {
			return;
		}

		if (!isValidEmployeeRole(role)) {
			return;
		}

		const usernameInput = form.elements.namedItem(
			"username",
		) as HTMLInputElement;

		const passwordInput = form.elements.namedItem(
			"password",
		) as HTMLInputElement;

		const firstNameInput = form.elements.namedItem(
			"firstName",
		) as HTMLInputElement;

		const lastNameInput = form.elements.namedItem(
			"lastName",
		) as HTMLInputElement;

		const roleInput = form.elements.namedItem("role") as HTMLInputElement;

		const dateOfBirthInput = form.elements.namedItem(
			"dateOfBirth",
		) as HTMLInputElement;

		const phoneNumberInput = form.elements.namedItem(
			"phoneNumber",
		) as HTMLInputElement;

		const locationInput = form.elements.namedItem(
			"location",
		) as HTMLInputElement;

		const scheduleStartInput = form.elements.namedItem(
			"scheduleStart",
		) as HTMLInputElement;

		const scheduleEndInput = form.elements.namedItem(
			"scheduleEnd",
		) as HTMLInputElement;

		// Check if form is valid to prevent making a server request.
		if (
			!validateForm(
				usernameInput.validity,
				passwordInput.validity,
				firstNameInput.validity,
				lastNameInput.validity,
				roleInput.validity,
				dateOfBirthInput.validity,
				phoneNumberInput.validity,
				locationInput.validity,
				scheduleStartInput.validity,
				scheduleEndInput.validity,
				selectedCoordinate,
			)
		) {
			return;
		}

		if (createForm) {
			(onSave as onSaveCreateType)(
				username,
				password,
				firstName,
				lastName,
				role,
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
			return;
		}

		(onSave as onSaveType)(
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
						placeholder={$t("employees.lastName.placeholder")}
						minLength={formFieldsLengths.lastName.min}
						maxLength={formFieldsLengths.lastName.max}
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

				{#if createForm}
					<!-- Password -->
					<!-- TODO pattern to password -->
					<FormControl
						label={$t("employees.password")}
						error={!!formErrorMessages.password}
						helperText={formErrorMessages.password}
					>
						<Input
							required
							name="password"
							type="password"
							pattern={``}
							error={!!formErrorMessages.password}
							placeholder={$t("employees.password.placeholder")}
							minLength={formFieldsLengths.password.min}
							maxLength={formFieldsLengths.password.max}
						/>
					</FormControl>
				{/if}

				<!-- dateOfBirth -->
				<FormControl
					label={$t("employees.dateOfBirth")}
					error={!!formErrorMessages.dateOfBirth}
					helperText={formErrorMessages.dateOfBirth}
				>
					<Input
						required
						name="dateOfBirth"
						value={employee?.dateOfBirth}
						error={!!formErrorMessages.dateOfBirth}
						placeholder={$t("employees.dateOfBirth.placeholder")}
						type="date"
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
						pattern={`^[0-9]{3}[ ]?[0-9]{3}[ ]?[0-9]{3}$`}
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
						step="1"
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
						step="1"
						value={employee?.scheduleEnd}
						error={!!formErrorMessages.scheduleEnd}
						placeholder={$t("employees.scheduleEnd.placeholder")}
						type="time"
					/>
				</FormControl>

				{#if createForm}
					<!-- Role -->
					<FormControl
						label={$t("employees.role")}
						error={!!formErrorMessages.role}
						helperText={formErrorMessages.role}
					>
						<Select
							required
							name="role"
							error={!!formErrorMessages.role}
							placeholder={$t("employees.role.placeholder")}
						>
							{#each rolesOptions as role}
								<Option value={role}>
									{$t(`employees.role.${role}`)}
								</Option>
							{/each}
						</Select>
					</FormControl>
				{/if}
			</DetailsFields>
		</DetailsSection>
	</DetailsContent>
</form>

<style>
	form {
		display: contents;
	}
</style>
