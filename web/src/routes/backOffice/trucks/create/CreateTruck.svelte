<script lang="ts">
	import { navigate } from "svelte-routing";
	import Card from "../../components/Card.svelte";
	import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
	import { t } from "../../../../lib/utils/i8n";
	import type { GeoJSONFeaturePoint } from "../../../../domain/geojson";
	import TruckForm from "../components/TruckForm.svelte";
	import { BackOfficeRoutes } from "../../../constants/routes";
	import { getToastContext } from "../../../../lib/contexts/toast";

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
	 * Creates a truck with a given truck make, model, license plate, person capacity and location.
	 * @param make Truck make.
	 * @param model Truck model.
	 * @param licensePlate Truck license plate.
	 * @param personCapacity Truck person capacity.
	 * @param location Truck location.
	 */
	async function createTruck(
		make: string,
		model: string,
		licensePlate: string,
		personCapacity: number,
		location: GeoJSONFeaturePoint,
	) {
		isSubmittingForm = true;

		const res = await ecomapHttpClient.POST("/trucks", {
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
			title: $t("trucks.create.success"),
			description: undefined,
		});

		navigate(`${BackOfficeRoutes.TRUCKS}/${res.data.id}`);
	}
</script>

<Card class="m-10 flex flex-col gap-10">
	<TruckForm
		back=""
		title={$t("trucks.create.title")}
		isSubmitting={isSubmittingForm}
		onSave={createTruck}
	/>
</Card>
