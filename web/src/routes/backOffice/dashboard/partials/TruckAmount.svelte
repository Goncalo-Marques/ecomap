<script lang="ts">
	import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
	import Spinner from "../../../../lib/components/Spinner.svelte";
	import { t } from "../../../../lib/utils/i8n";
	import Card from "../../components/Card.svelte";

	/**
	 * Retrieves the number of trucks.
	 * @returns Number of trucks.
	 */
	async function getTruckAmount() {
		const res = await ecomapHttpClient.GET("/trucks");

		if (res.error) {
			throw new Error("Failed to retrieve truck amount");
		}

		return res.data.total;
	}

	const truckAmountPromise = getTruckAmount();
</script>

<Card element="article" class="truck-amount-card">
	<h2>{$t("dashboard.amountOfTrucks")}</h2>
	<p>
		{#await truckAmountPromise}
			<Spinner />
		{:then truckAmount}
			{truckAmount}
		{:catch}
			{$t("error.unexpected.title")}
		{/await}
	</p>
</Card>

<style>
	:global(.truck-amount-card) {
		grid-area: truckAmount;
	}

	h2 {
		font: var(--text-base-regular);
		color: var(--gray-500);
	}

	p {
		font: var(--text-xl-semibold);
	}
</style>
