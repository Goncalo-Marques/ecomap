<script lang="ts">
	import { goto } from "$app/navigation";
	import type { EmployeeRole } from "$domain/employees";
	import type { GeoJSONFeaturePoint } from "$domain/geojson";
	import ecomapHttpClient from "$lib/clients/ecomap/http";
	import { BackOfficeRoutes } from "$lib/constants/routes";
	import { getToastContext } from "$lib/contexts/toast";
	import { t } from "$lib/utils/i8n";

	import EmployeeForm from "../components/EmployeeForm.svelte";

	/**
	 * Toast context.
	 */
	const toast = getToastContext();

	/**
	 * Indicates if form is being submitted.
	 * @default false
	 */
	let isSubmittingForm: boolean = false;

	/**
	 * Creates a new employee with a given username, password, firstName, lastName, role, dateOfBirth, phoneNumber, location, scheduleStart and scheduleEnd.
	 * @param username Employee username.
	 * @param password Employee password.
	 * @param firstName Employee firstName.
	 * @param lastName Employee lastName.
	 * @param role Employee role.
	 * @param dateOfBirth Employee dateOfBirth.
	 * @param phoneNumber Employee phoneNumber.
	 * @param location Employee location.
	 * @param scheduleStart Employee scheduleStart.
	 * @param scheduleEnd Employee scheduleEnd.
	 */
	async function createEmployee(
		username: string,
		password: string,
		firstName: string,
		lastName: string,
		role: EmployeeRole,
		dateOfBirth: string,
		phoneNumber: string,
		location: GeoJSONFeaturePoint,
		scheduleStart: string,
		scheduleEnd: string,
	) {
		isSubmittingForm = true;

		// Adding seconds to times. Necessary because API receives times with seconds.
		scheduleStart += ":00";
		scheduleEnd += ":00";

		const res = await ecomapHttpClient.POST("/employees", {
			body: {
				username,
				password,
				firstName,
				lastName,
				role,
				dateOfBirth,
				phoneNumber,
				geoJson: location,
				scheduleStart,
				scheduleEnd,
			},
		});

		isSubmittingForm = false;

		if (res.error) {
			switch (res.error.code) {
				case "bad_request":
					toast.show({
						type: "error",
						title: $t("employees.password.error.passwordConstraints.title"),
						description: $t(
							"employees.password.error.passwordConstraints.description",
						),
					});
					break;
				case "conflict":
					toast.show({
						type: "error",
						title: $t("employees.create.conflict.title"),
						description: $t("employees.create.conflict.description"),
					});
					break;

				default:
					toast.show({
						type: "error",
						title: $t("error.unexpected.title"),
						description: $t("error.unexpected.description"),
					});
					break;
			}

			return;
		}

		toast.show({
			type: "success",
			title: $t("employees.create.success"),
			description: undefined,
		});

		goto(`${BackOfficeRoutes.EMPLOYEES}/${res.data.id}`);
	}
</script>

<EmployeeForm
	back={BackOfficeRoutes.EMPLOYEES}
	action="create"
	title={$t("employees.create.title")}
	isSubmitting={isSubmittingForm}
	onSave={createEmployee}
/>
