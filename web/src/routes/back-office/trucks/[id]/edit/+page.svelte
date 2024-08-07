<script lang="ts">
	import { goto } from "$app/navigation";
	import { page } from "$app/stores";
	import type { GeoJSONFeaturePoint } from "$domain/geojson";
	import type { Truck } from "$domain/truck";
	import ecomapHttpClient from "$lib/clients/ecomap/http";
	import Spinner from "$lib/components/Spinner.svelte";
	import { BackOfficeRoutes } from "$lib/constants/routes";
	import { getToastContext } from "$lib/contexts/toast";
	import { t } from "$lib/utils/i8n";

	import Card from "../../../components/Card.svelte";
	import TruckForm from "../../components/TruckForm.svelte";

	/**
	 * Truck ID.
	 */
	const id: string = $page.params.id;

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

		goto(`${BackOfficeRoutes.TRUCKS}/${id}`);
	}

	const truckPromise = fetchTruck();
</script>

<Card class="m-10 flex flex-col gap-10">
	{#await truckPromise}
		<Spinner class="flex h-full items-center justify-center" />
	{:then truck}
		<TruckForm
			{truck}
			isSubmitting={isSubmittingForm}
			back={`${BackOfficeRoutes.TRUCKS}/${truck.id}`}
			title={`${truck.make} ${truck.model}`}
			onSave={updateTruck}
		/>
	{:catch}
		<div class="flex h-1/2 flex-col items-center justify-center">
			<h2 class="text-2xl font-semibold">{$t("trucks.notFound.title")}</h2>
			<p>{$t("trucks.notFound.description")}</p>
		</div>
	{/await}
</Card>
