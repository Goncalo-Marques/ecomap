<script lang="ts">
	import { navigate } from "svelte-routing";
	import TruckForm from "../components/TruckForm.svelte";
	import { BackOfficeRoutes } from "../../../constants/routes";
	import { t } from "../../../../lib/utils/i8n";
	import Spinner from "../../../../lib/components/Spinner.svelte";
	import Card from "../../components/Card.svelte";
	import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
	import type { GeoJSONFeaturePoint } from "../../../../domain/geojson";
	import type { Truck } from "../../../../domain/truck";
	import { getToastContext } from "../../../../lib/contexts/toast";

	/**
	 * Truck ID.
	 */
	export let id: string;

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
	 * Fetches truck data.
	 */
	async function fetchTruck(): Promise<Truck> {
		const res = await ecomapHttpClient.GET("/trucks/{truckId}", {
			params: { path: { truckId: id } },
		});

		if (res.error) {
			throw new Error("Failed to retrieve truck details");
		}

		return res.data;
	}

	/**
	 * Updates a truck with a given truck make, model, license plate, person capacity and location.
	 * @param make Truck make.
	 * @param model Truck model.
	 * @param licensePlate Truck license plate.
	 * @param personCapacity Truck person capacity.
	 * @param location Truck location.
	 */
	async function updateTruck(
		make: string,
		model: string,
		licensePlate: string,
		personCapacity: number,
		location: GeoJSONFeaturePoint,
	) {
		isSubmittingForm = true;

		const res = await ecomapHttpClient.PATCH("/trucks/{truckId}", {
			params: {
				path: {
					truckId: id,
				},
			},
			body: {
				make,
				model,
				licensePlate,
				personCapacity,
				geoJson: location,
			},
		});

		isSubmittingForm = false;

		if (res.error) {
			toast.show({
				type: "error",
				title: $t("error.unexpected.title"),
				description: $t("error.unexpected.description"),
			});

			return;
		}

		toast.show({
			type: "success",
			title: $t("trucks.update.success"),
			description: undefined,
		});

		navigate(`${BackOfficeRoutes.TRUCKS}/${id}`);
	}

	const truckPromise = fetchTruck();
</script>

<Card class="page-layout">
	{#await truckPromise}
		<div class="truck-loading">
			<Spinner />
		</div>
	{:then truck}
		<TruckForm
			{truck}
			isSubmitting={isSubmittingForm}
			back={truck.id}
			title={`${truck.make} ${truck.model}`}
			onSave={updateTruck}
		/>
	{:catch}
		<div class="truck-not-found">
			<h2>{$t("trucks.notFound.title")}</h2>
			<p>{$t("trucks.notFound.description")}</p>
		</div>
	{/await}
</Card>

<style>
	.truck-loading {
		height: 100%;
		display: flex;
		justify-content: center;
		align-items: center;
	}

	.truck-not-found {
		height: 50%;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
	}
</style>
