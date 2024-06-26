<script lang="ts">
	import { navigate } from "svelte-routing";
	import type {
		Container,
		ContainerCategory,
	} from "../../../../domain/container";
	import Spinner from "../../../../lib/components/Spinner.svelte";
	import Card from "../../components/Card.svelte";
	import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
	import { t } from "../../../../lib/utils/i8n";
	import type { GeoJSONFeaturePoint } from "../../../../domain/geojson";
	import ContainerForm from "../components/ContainerForm.svelte";
	import { BackOfficeRoutes } from "../../../constants/routes";
	import { getToastContext } from "../../../../lib/contexts/toast";
	import { getLocationName } from "../../../../lib/utils/location";

	/**
	 * Container ID.
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
	 * Fetches container data.
	 */
	async function fetchContainer(): Promise<Container> {
		const res = await ecomapHttpClient.GET("/containers/{containerId}", {
			params: { path: { containerId: id } },
		});

		if (res.error) {
			throw new Error("Failed to retrieve container details");
		}

		return res.data;
	}

	/**
	 * Updates a container with a given category and location.
	 * @param category Container category.
	 * @param location Container location.
	 */
	async function updateContainer(
		category: ContainerCategory,
		location: GeoJSONFeaturePoint,
	) {
		isSubmittingForm = true;

		const res = await ecomapHttpClient.PATCH("/containers/{containerId}", {
			params: {
				path: {
					containerId: id,
				},
			},
			body: {
				category,
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
			title: $t("containers.update.success"),
			description: undefined,
		});

		navigate(`${BackOfficeRoutes.CONTAINERS}/${id}`);
	}

	const containerPromise = fetchContainer();
</script>

<Card class="page-layout">
	{#await containerPromise}
		<div class="container-loading">
			<Spinner />
		</div>
	{:then container}
		{@const locationName = getLocationName(
			container.geoJson.properties.wayName,
			container.geoJson.properties.municipalityName,
		)}
		<ContainerForm
			{container}
			isSubmitting={isSubmittingForm}
			back={container.id}
			title={locationName}
			onSave={updateContainer}
		/>
	{:catch}
		<div class="container-not-found">
			<h2>{$t("containers.notFound.title")}</h2>
			<p>{$t("containers.notFound.description")}</p>
		</div>
	{/await}
</Card>

<style>
	.container-loading {
		height: 100%;
		display: flex;
		justify-content: center;
		align-items: center;
	}

	.container-not-found {
		height: 50%;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
	}

	:global(.container-map-preview) {
		flex: 1;
	}
</style>
