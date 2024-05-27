<script lang="ts">
	import { navigate } from "svelte-routing";
	import Card from "../../components/Card.svelte";
	import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
	import { t } from "../../../../lib/utils/i8n";
	import type { GeoJSONFeaturePoint } from "../../../../domain/geojson";
	import { BackOfficeRoutes } from "../../../constants/routes";
	import { getToastContext } from "../../../../lib/contexts/toast";
	import EmployeeForm from "../components/EmployeeForm.svelte";
	import type { EmployeeRoles } from "../../../../domain/employees";

	/**
	 * Toast context.
	 */
	const toast = getToastContext();

	/**
	 * Creates a new employee with given username, password, firstName, lastName, role, dateOfBirth, phoneNumber, location, scheduleStart and scheduleEnd.
	 *
	 * @param username
	 * @param password
	 * @param firstName
	 * @param lastName
	 * @param role
	 * @param dateOfBirth
	 * @param phoneNumber
	 * @param location
	 * @param scheduleStart
	 * @param scheduleEnd
	 */
	async function createEmployee(
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
	) {
		// Adding seconds to times. Necessary because api receives times with seconds.
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

		if (res.error) {
			switch (res.error.code) {
				case "bad_request":
					toast.show({
						type: "error",
						title: $t(
							"employees.updatePassword.error.passwordConstraints.title",
						),
						description: $t(
							"employees.updatePassword.error.passwordConstraints.description",
						),
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

		navigate(`${BackOfficeRoutes.EMPLOYEES}/${res.data.id}`);
	}
</script>

<Card class="page-layout">
	<EmployeeForm
		back=""
		title={$t("employees.create.title")}
		onSave={createEmployee}
		createForm
	/>
</Card>
