<script lang="ts">
	import type { Employee } from "../../../domain/employees";
	import ecomapHttpClient from "../../../lib/clients/ecomap/http";
	import Spinner from "../../../lib/components/Spinner.svelte";
	import { t } from "../../../lib/utils/i8n";
	import { getBatchPaginatedResponse } from "../../../lib/utils/request";
	import Card from "../components/Card.svelte";

	/**
	 * Retrieves employees to be displayed in the employees table.
	 * @returns Active employees.
	 */
	async function getActiveEmployees(): Promise<number> {
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

		// let activeEmployees = 0;

		// const today = new Date();
		// const unixTime = new Date(0);
		// unixTime.setHours(today.getHours());
		// unixTime.setMinutes(today.getMinutes());
		// unixTime.setSeconds(today.getSeconds());

		// const scheduleStart = new Date(0);
		// const scheduleEnd = new Date(0);

		// for (const employee of employees) {
		// 	const [scheduleStartHours, scheduleStartMinutes, scheduleStartSeconds] =
		// 		employee.scheduleStart.split(":").map(Number);
		// 	scheduleStart.setHours(scheduleStartHours);
		// 	scheduleStart.setMinutes(scheduleStartMinutes);
		// 	scheduleStart.setSeconds(scheduleStartSeconds);

		// 	const [scheduleEndHours, scheduleEndMinutes, scheduleEndSeconds] =
		// 		employee.scheduleEnd.split(":").map(Number);
		// 	scheduleEnd.setHours(scheduleEndHours);
		// 	scheduleEnd.setMinutes(scheduleEndMinutes);
		// 	scheduleEnd.setSeconds(scheduleEndSeconds);

		// 	if (unixTime >= scheduleStart && unixTime <= scheduleEnd) {
		// 		activeEmployees++;
		// 	}
		// }

		return employees.length;
	}

	const activeEmployeesPromise = getActiveEmployees();
</script>

<main class="page-layout">
	<h1>{$t("dashboard")}</h1>
	<div class="row">
		<Card>
			<h2>Colaboradores ativos</h2>
			{#await activeEmployeesPromise}
				<Spinner />
			{:then activeEmployees}
				<p>{activeEmployees}</p>
			{:catch}
				<p>Erro</p>
			{/await}
		</Card>
		<Card>
			<h2>Colaboradores ativos</h2>
			{#await activeEmployeesPromise}
				<Spinner />
			{:then activeEmployees}
				<p>{activeEmployees}</p>
			{:catch}
				<p>Erro</p>
			{/await}
		</Card>
	</div>
</main>

<style>
	.row {
		display: flex;
		gap: 1rem;

		& > * {
			flex: 1;
		}
	}
</style>
