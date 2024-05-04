<script lang="ts">
	import { Link, navigate } from "svelte-routing";
	import type {
		Container,
		ContainerCategory,
	} from "../../../../domain/container";
	import Spinner from "../../../../lib/components/Spinner.svelte";
	import Button from "../../../../lib/components/Button.svelte";
	import Card from "../../components/Card.svelte";
	import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
	import { t } from "../../../../lib/utils/i8n";
	import Field from "../../../../lib/components/Field.svelte";
	import DetailsFields from "../../../../lib/components/details/DetailsFields.svelte";
	import DetailsSection from "../../../../lib/components/details/DetailsSection.svelte";
	import DetailsContent from "../../../../lib/components/details/DetailsContent.svelte";
	import DetailsHeader from "../../../../lib/components/details/DetailsHeader.svelte";
	import { formatDate } from "../../../../lib/utils/date";
	import { DateFormats } from "../../../../lib/constants/date";
	import { getContainerLocation } from "../utils/location";

	/**
	 * Container ID.
	 */
	export let id: string;

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

	const containerPromise = fetchContainer();
</script>

<Card class="page-layout">
	{#await containerPromise}
		<div class="container-loading">
			<Spinner />
		</div>
	{:then container}
		{@const locationName = getContainerLocation(
			container.geoJson.properties.wayName,
			container.geoJson.properties.municipalityName,
		)}
		<DetailsHeader to="" title={locationName}>
			<Link to={`${container.id}/map`}>
				<Button variant="secondary" startIcon="map">
					{$t("sidebar.map")}
				</Button>
			</Link>
			<Link to={`${container.id}/edit`}>
				<Button startIcon="edit">Editar informação</Button>
			</Link>
		</DetailsHeader>
		<DetailsContent>
			<DetailsSection label={$t("generalInfo")}>
				<DetailsFields>
					<Field
						label={$t("containers.category")}
						value={$t(`containers.category.${container.category}`)}
					/>
					<Field label={$t("containers.location")} value={locationName} />
				</DetailsFields>
			</DetailsSection>
			<DetailsSection label={$t("additionalInfo")}>
				<DetailsFields>
					<Field
						label={$t("createdAt")}
						value={formatDate(container.createdAt, DateFormats.shortDateTime)}
					/>
					<Field
						label={$t("modifiedAt")}
						value={formatDate(container.modifiedAt, DateFormats.shortDateTime)}
					/>
				</DetailsFields>
			</DetailsSection>
		</DetailsContent>
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
</style>
