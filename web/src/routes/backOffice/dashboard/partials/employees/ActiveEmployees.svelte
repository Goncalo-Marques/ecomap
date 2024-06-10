<script lang="ts">
	import ecomapHttpClient from "../../../../../lib/clients/ecomap/http";
	import Spinner from "../../../../../lib/components/Spinner.svelte";
	import { getBatchPaginatedResponse } from "../../../../../lib/utils/request";
	import Card from "../../../components/Card.svelte";
	import { getActiveEmployees } from "./utils";

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

<Card element="article">
	<h2>Colaboradores ativos</h2>
	<p>
		{#await activeEmployeeAmountPromise}
			<Spinner />
		{:then activeEmployeeAmount}
			{activeEmployeeAmount}
		{:catch}
			Error
		{/await}
	</p>
</Card>

<style>
	:global(.active-employees-card) {
		grid-area: activeEmployees;
	}

	h2 {
		font: var(--text-base-regular);
		color: var(--gray-500);
	}

	p {
		font: var(--text-xl-semibold);
	}
</style>
