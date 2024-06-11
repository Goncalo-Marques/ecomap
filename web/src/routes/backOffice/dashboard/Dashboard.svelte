<script lang="ts">
	import ecomapHttpClient from "../../../lib/clients/ecomap/http";
	import { t } from "../../../lib/utils/i8n";
	import { getBatchPaginatedResponse } from "../../../lib/utils/request";
	import ContainersAdded from "./partials/ContainersAdded.svelte";
	import ContainersByCategory from "./partials/ContainersByCategory.svelte";
	import ContainersByMunicipality from "./partials/ContainersByMunicipality.svelte";
	import TruckAmount from "./partials/TruckAmount.svelte";
	import WarehouseAmount from "./partials/WarehouseAmount.svelte";
	import ActiveEmployees from "./partials/ActiveEmployees.svelte";

	/**
	 * Retrieves the containers to be displayed in the charts.
	 * @returns Containers.
	 */
	async function getContainers() {
		const containers = await getBatchPaginatedResponse(
			async (limit, offset) => {
				const res = await ecomapHttpClient.GET("/containers", {
					params: {
						query: {
							offset,
							limit,
							sort: "createdAt",
							order: "asc",
						},
					},
				});

				if (res.error) {
					return { total: 0, items: [] };
				}

				return { total: res.data.total, items: res.data.containers };
			},
		);

		return containers;
	}

	let containersPromise = getContainers();
</script>

<main class="page-layout">
	<h1>{$t("dashboard")}</h1>
	<div class="dashboard-content">
		<ActiveEmployees />
		<WarehouseAmount />
		<TruckAmount />
		<ContainersAdded {containersPromise} />
		<ContainersByCategory {containersPromise} />
		<ContainersByMunicipality {containersPromise} />
	</div>
</main>

<style>
	h1 {
		font: var(--text-2xl-semibold);
	}

	.dashboard-content {
		display: grid;
		grid-template-columns: 1fr 1fr 1fr;
		grid-template-areas:
			"activeEmployees warehouseAmount truckAmount"
			"containersAdded containersAdded containersByCategory"
			"containersByMunicipality containersByMunicipality containersByMunicipality";
		gap: 1rem;
	}
</style>
