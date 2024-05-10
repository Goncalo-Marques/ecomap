<script lang="ts">
	import Card from "../../components/Card.svelte";
	import Spinner from "../../../../lib/components/Spinner.svelte";
	import { navigate } from "svelte-routing";
	import Button from "../../../../lib/components/Button.svelte";
	import { t } from "../../../../lib/utils/i8n";
	import DetailsHeader from "../../../../lib/components/details/DetailsHeader.svelte";
	import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
	import type { Employee } from "../../../../domain/employees";
	import DetailsContent from "../../../../lib/components/details/DetailsContent.svelte";
	import DetailsSection from "../../../../lib/components/details/DetailsSection.svelte";
	import { getLocationName } from "../../../../lib/utils/location";
	import DetailsFields from "../../../../lib/components/details/DetailsFields.svelte";
	import Field from "../../../../lib/components/Field.svelte";
	import { formatDate, formatTime } from "../../../../lib/utils/date";
	import { DateFormats } from "../../../../lib/constants/date";
	import { getToastContext } from "../../../../lib/contexts/toast";
	import { BackOfficeRoutes } from "../../../constants/routes";

	/**
	 * Employee ID.
	 */
	export let id: string;

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
	 * Deletes the employee displayed on the page.
	 */
	async function deleteEmployee() {
		const res = await ecomapHttpClient.DELETE("/employees/{employeeId}", {
			params: {
				path: {
					employeeId: id,
				},
			},
		});

		if (res.error) {
			if (res.error.code === "conflict") {
				toast.show({
					type: "error",
					title: $t("employees.delete.conflict.title"),
					description: $t("employees.delete.conflict.description"),
				});
			} else {
				toast.show({
					type: "error",
					title: $t("error.unexpected.title"),
					description: $t("error.unexpected.description"),
				});
			}

			return;
		}

		toast.show({
			type: "success",
			title: $t("employees.delete.success"),
			description: undefined,
		});

		navigate(BackOfficeRoutes.EMPLOYEES);
	}

	const employeePromise = fetchEmployee();
</script>

<Card class="page-layout">
	{#await employeePromise}
		<div class="employee-loading">
			<Spinner />
		</div>
	{:then employee}
		{@const locationName = getLocationName(
			employee.geoJson.properties.wayName,
			employee.geoJson.properties.municipalityName,
		)}

		<DetailsHeader to="" title={`${employee.firstName} ${employee.lastName}`}>
			<Button
				startIcon="delete"
				actionType="danger"
				variant="secondary"
				onClick={deleteEmployee}
			/>
		</DetailsHeader>
		<DetailsContent>
			<DetailsSection label={$t("personalInfo")}>
				<DetailsFields>
					<Field label={$t("employees.firstName")} value={employee.firstName} />
					<Field label={$t("employees.lastName")} value={employee.lastName} />
					<Field label={$t("employees.username")} value={employee.username} />
					<Field
						label={$t("employees.birthdate")}
						value={formatDate(employee.dateOfBirth, DateFormats.shortDate)}
					/>
					<Field label={$t("employees.phone")} value={employee.phoneNumber} />
					<Field label={$t("employees.location")} value={locationName} />
				</DetailsFields>
			</DetailsSection>
			<DetailsSection label={$t("work")}>
				<DetailsFields>
					<Field
						label={$t("employees.scheduleStart")}
						value={formatTime(employee.scheduleStart)}
					/>
					<Field
						label={$t("employees.scheduleEnd")}
						value={formatTime(employee.scheduleEnd)}
					/>
					<Field
						label={$t("employees.role")}
						value={$t(`employees.role.${employee.role}`)}
					/>
				</DetailsFields>
			</DetailsSection>
			<DetailsSection label={$t("additionalInfo")}>
				<DetailsFields>
					<Field
						label={$t("createdAt")}
						value={formatDate(employee.createdAt, DateFormats.shortDateTime)}
					/>
					<Field
						label={$t("modifiedAt")}
						value={formatDate(employee.modifiedAt, DateFormats.shortDateTime)}
					/>
				</DetailsFields>
			</DetailsSection>
		</DetailsContent>
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
