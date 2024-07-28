<script lang="ts">
	import { goto } from "$app/navigation";
	import { page } from "$app/stores";
	import type { Employee } from "$domain/employees";
	import ecomapHttpClient from "$lib/clients/ecomap/http";
	import Button from "$lib/components/Button.svelte";
	import DetailsContent from "$lib/components/details/DetailsContent.svelte";
	import DetailsFields from "$lib/components/details/DetailsFields.svelte";
	import DetailsHeader from "$lib/components/details/DetailsHeader.svelte";
	import DetailsSection from "$lib/components/details/DetailsSection.svelte";
	import Field from "$lib/components/Field.svelte";
	import Spinner from "$lib/components/Spinner.svelte";
	import { DateFormats } from "$lib/constants/date";
	import { getToastContext } from "$lib/contexts/toast";
	import { isViewingSelf } from "$lib/utils/auth";
	import { formatDate, formatTime } from "$lib/utils/date";
	import { getSupportedLocaleName, locale, t } from "$lib/utils/i8n";
	import { getLocationName } from "$lib/utils/location";

	import { BackOfficeRoutes } from "../../../constants/routes";
	import UpdatePasswordModal from "../components/UpdatePasswordModal.svelte";

	/**
	 * Toast context.
	 */
	const toast = getToastContext();

	/**
	 * The update password open modal state.
	 * @default false
	 */
	let openUpdatePasswordModal = false;

	/**
	 * Fetches employee data.
	 */
	async function fetchEmployee(id: string): Promise<Employee> {
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
	async function deleteEmployee(id: string) {
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

		goto(BackOfficeRoutes.EMPLOYEES);
	}

	let employeePromise: Promise<Employee>;

	// Refetch employee information whenever the employee ID changes in the URL.
	$: employeePromise = fetchEmployee($page.params.id);
</script>

{#await employeePromise}
	<Spinner class="flex h-full items-center justify-center" />
{:then employee}
	{@const locationName = getLocationName(
		employee.geoJson.properties.wayName,
		employee.geoJson.properties.municipalityName,
	)}
	{@const isSelf = isViewingSelf(employee.id)}

	<DetailsHeader
		href={BackOfficeRoutes.EMPLOYEES}
		title={`${employee.firstName} ${employee.lastName}`}
	>
		{#if isSelf}
			<Button
				variant="secondary"
				onClick={() => (openUpdatePasswordModal = true)}
			>
				{$t("employees.updatePassword.title")}
			</Button>
		{:else}
			<Button
				startIcon="delete"
				actionType="danger"
				variant="secondary"
				onClick={() => deleteEmployee($page.params.id)}
			/>
		{/if}
		<a
			href={`${BackOfficeRoutes.EMPLOYEES}/${employee.id}/edit`}
			class="contents"
		>
			<Button startIcon="edit">{$t("editInfo")}</Button>
		</a>
	</DetailsHeader>
	<DetailsContent>
		<DetailsSection label={$t("personalInfo")}>
			<DetailsFields>
				<Field label={$t("employees.firstName")} value={employee.firstName} />
				<Field label={$t("employees.lastName")} value={employee.lastName} />
				<Field label={$t("employees.username")} value={employee.username} />
				<Field
					label={$t("employees.dateOfBirth")}
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
		{#if isSelf}
			<DetailsSection label={$t("preferences")}>
				<DetailsFields>
					<Field
						label={$t("language")}
						value={getSupportedLocaleName($locale)}
					/>
				</DetailsFields>
			</DetailsSection>
		{/if}
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
	<UpdatePasswordModal
		{employee}
		open={openUpdatePasswordModal}
		onOpenChange={open => (openUpdatePasswordModal = open)}
	/>
{:catch}
	<div class="flex h-1/2 flex-col items-center justify-center">
		<h2 class="text-2xl font-semibold">{$t("employees.notFound.title")}</h2>
		<p>{$t("employees.notFound.description")}</p>
	</div>
{/await}
