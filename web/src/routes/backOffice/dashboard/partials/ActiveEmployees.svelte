<script lang="ts">
	import type { Employee } from "../../../../domain/employees";
	import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
	import Spinner from "../../../../lib/components/Spinner.svelte";
	import { t } from "../../../../lib/utils/i8n";
	import { getBatchPaginatedResponse } from "../../../../lib/utils/request";
	import Card from "../../components/Card.svelte";

	/**
	 * Retrieves list of active employees.
	 * @param employees Employees.
	 * @returns Active employees.
	 */
	function getActiveEmployees(employees: Employee[]): Employee[] {
		const today = new Date();

		const currentDate = new Date(0);
		currentDate.setHours(today.getHours());
		currentDate.setMinutes(today.getMinutes());
		currentDate.setSeconds(today.getSeconds());

		const activeEmployees: Employee[] = [];

		for (const employee of employees) {
			const scheduleStart = new Date(0);
			const scheduleEnd = new Date(0);

			// Set start time of employee schedule.
			const [scheduleStartHours, scheduleStartMinutes, scheduleStartSeconds] =
				employee.scheduleStart.split(":").map(Number);
			scheduleStart.setHours(scheduleStartHours);
			scheduleStart.setMinutes(scheduleStartMinutes);
			scheduleStart.setSeconds(scheduleStartSeconds);

			// Set end time of employee schedule.
			const [scheduleEndHours, scheduleEndMinutes, scheduleEndSeconds] =
				employee.scheduleEnd.split(":").map(Number);
			scheduleEnd.setHours(scheduleEndHours);
			scheduleEnd.setMinutes(scheduleEndMinutes);
			scheduleEnd.setSeconds(scheduleEndSeconds);

			// If the end time of the employee's schedule is before the start time
			// of the schedule, it means that the employee is working at night.
			if (scheduleEndHours < scheduleStartHours) {
				scheduleEnd.setDate(scheduleEnd.getDate() + 1);
			}

			// Check if today's time is between the employee's scheduled time.
			if (currentDate >= scheduleStart && currentDate <= scheduleEnd) {
				activeEmployees.push(employee);
			}
		}

		return activeEmployees;
	}

	/**
	 * Retrieves the number of active employees.
	 * @returns Number of active employees.
	 */
	async function getActiveEmployeeAmount() {
		const employees = await getBatchPaginatedResponse(async (limit, offset) => {
			const res = await ecomapHttpClient.GET("/employees", {
				params: {
					query: {
						offset,
						limit,
					},
				},
			});

			if (res.error) {
				return { total: 0, items: [] };
			}

			return { total: res.data.total, items: res.data.employees };
		});

		return getActiveEmployees(employees).length;
	}

	const activeEmployeeAmountPromise = getActiveEmployeeAmount();
</script>

<Card element="article" class="[grid-area:activeEmployees]">
	<h2 class="text-gray-500">{$t("dashboard.activeEmployees")}</h2>
	<p class="text-xl font-semibold">
		{#await activeEmployeeAmountPromise}
			<Spinner />
		{:then activeEmployeeAmount}
			{activeEmployeeAmount}
		{:catch}
			{$t("error.unexpected.title")}
		{/await}
	</p>
</Card>
