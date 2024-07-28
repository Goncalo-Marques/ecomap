<script lang="ts">
	import { goto } from "$app/navigation";
	import { page } from "$app/stores";
	import type { Employee } from "$domain/employees";
	import type { GeoJSONFeaturePoint } from "$domain/geojson";
	import ecomapHttpClient from "$lib/clients/ecomap/http";
	import Spinner from "$lib/components/Spinner.svelte";
	import { BackOfficeRoutes } from "$lib/constants/routes";
	import { getToastContext } from "$lib/contexts/toast";
	import { isViewingSelf } from "$lib/utils/auth";
	import { locale, t } from "$lib/utils/i8n";

	import EmployeeForm from "../../components/EmployeeForm.svelte";

	/**
	 * Employee ID.
	 */
	const id: string = $page.params.id;

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
	 * Fetches employee data.
	 */
	async function fetchEmployee(): Promise<Employee> {
		const res = await ecomapHttpClient.GET("/employees/{employeeId}", {
			params: { path: { employeeId: id } },
		});

		if (res.error) {
			throw new Error("Failed to retrieve employee details");
		}

		return res.data;
	}

	/**
	 * Updates employee with given username, firstName, lastName, dateOfBirth, phoneNumber, location, scheduleStart and scheduleEnd.
	 * @param username Employee username.
	 * @param firstName Employee firstName.
	 * @param lastName Employee lastName.
	 * @param dateOfBirth Employee dateOfBirth.
	 * @param phoneNumber Employee phoneNumber.
	 * @param location Employee location.
	 * @param scheduleStart Employee scheduleStart.
	 * @param scheduleEnd Employee scheduleEnd.
	 * @param selectedLocale Selected locale for the application.
	 */
	async function updateEmployee(
		username: string,
		firstName: string,
		lastName: string,
		dateOfBirth: string,
		phoneNumber: string,
		location: GeoJSONFeaturePoint,
		scheduleStart: string,
		scheduleEnd: string,
		selectedLocale: string,
	) {
		isSubmittingForm = true;

		// Adding seconds to times. Necessary because API receives times with seconds.
		scheduleStart += ":00";
		scheduleEnd += ":00";

		const res = await ecomapHttpClient.PATCH("/employees/{employeeId}", {
			params: {
				path: {
					employeeId: id,
				},
			},
			body: {
				username,
				firstName,
				lastName,
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
				case "conflict":
					toast.show({
						type: "error",
						title: $t("employees.update.conflict.title"),
						description: $t("employees.update.conflict.description"),
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

		// If the employee is viewing themselves, update the application locale.
		if (isViewingSelf(id)) {
			locale.set(selectedLocale);
		}

		toast.show({
			type: "success",
			title: $t("employees.update.success"),
			description: undefined,
		});

		goto(`${BackOfficeRoutes.EMPLOYEES}/${id}`);
	}

	const employeePromise = fetchEmployee();
</script>

{#await employeePromise}
	<Spinner class="flex h-full items-center justify-center" />
{:then employee}
	<EmployeeForm
		{employee}
		action="edit"
		title={`${employee.firstName} ${employee.lastName}`}
		isSubmitting={isSubmittingForm}
		back={`${BackOfficeRoutes.EMPLOYEES}/${employee.id}`}
		onSave={updateEmployee}
	/>
{:catch}
	<div class="flex h-1/2 flex-col items-center justify-center">
		<h2>{$t("employees.notFound.title")}</h2>
		<p>{$t("employees.notFound.description")}</p>
	</div>
{/await}
