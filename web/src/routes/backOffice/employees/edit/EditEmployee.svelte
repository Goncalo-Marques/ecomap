<script lang="ts">
	import { navigate } from "svelte-routing";
	import { BackOfficeRoutes } from "../../../constants/routes";
	import { t } from "../../../../lib/utils/i8n";
	import Spinner from "../../../../lib/components/Spinner.svelte";
	import Card from "../../components/Card.svelte";
	import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
	import type { GeoJSONFeaturePoint } from "../../../../domain/geojson";
	import type { Employee } from "../../../../domain/employees";
	import { getToastContext } from "../../../../lib/contexts/toast";
	import EmployeeForm from "../components/EmployeeForm.svelte";

	/**
	 * Employee ID.
	 */
	export let id: string;

	/**
	 * Toast context.
	 */
	const toast = getToastContext();

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
	) {
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

		toast.show({
			type: "success",
			title: $t("employees.update.success"),
			description: undefined,
		});

		navigate(`${BackOfficeRoutes.EMPLOYEES}/${id}`);
	}

	const employeePromise = fetchEmployee();
</script>

<Card class="page-layout">
	{#await employeePromise}
		<div class="employee-loading">
			<Spinner />
		</div>
	{:then employee}
		<EmployeeForm
			{employee}
			back={employee.id}
			title={`${employee.firstName} ${employee.lastName}`}
			onSave={updateEmployee}
			action="edit"
		/>
	{:catch}
		<div class="employee-not-found">
			<h2>{$t("employees.notFound.title")}</h2>
			<p>{$t("employees.notFound.description")}</p>
		</div>
	{/await}
</Card>

<style>
	.employee-loading {
		height: 100%;
		display: flex;
		justify-content: center;
		align-items: center;
	}

	.employee-not-found {
		height: 50%;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
	}
</style>
