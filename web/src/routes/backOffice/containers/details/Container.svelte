<script lang="ts">
	import { Link, navigate } from "svelte-routing";
	import type { Container } from "../../../../domain/container";
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
	import { getLocationName } from "../../../../lib/utils/location";
	import { BackOfficeRoutes } from "../../../constants/routes";
	import { getToastContext } from "../../../../lib/contexts/toast";

	/**
	 * Container ID.
	 */
	export let id: string;

	/**
	 * Toast context.
	 */
	const toast = getToastContext();

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
	 * Deletes the container displayed on the page.
	 */
	async function deleteContainer() {
		const res = await ecomapHttpClient.DELETE("/containers/{containerId}", {
			params: {
				path: {
					containerId: id,
				},
			},
		});

		if (res.error) {
			if (res.error.code === "conflict") {
				toast.show({
					type: "error",
					title: $t("containers.delete.conflict.title"),
					description: $t("containers.delete.conflict.description"),
				});
			} else {
				toast.show({
					type: "error",
					title: $t("error.unexpected.title"),
					description: $t("error.unexpected.description"),
				});
			}

			return;
		}

		toast.show({
			type: "success",
			title: $t("containers.delete.success"),
			description: undefined,
		});

		navigate(BackOfficeRoutes.CONTAINERS);
	}

	const containerPromise = fetchContainer();
</script>

<Card class="m-10 flex flex-col gap-10">
	{#await containerPromise}
		<Spinner class="flex h-full items-center justify-center" />
	{:then container}
		{@const locationName = getLocationName(
			container.geoJson.properties.wayName,
			container.geoJson.properties.municipalityName,
		)}
		<DetailsHeader to="" title={locationName}>
			<Button
				startIcon="delete"
				actionType="danger"
				variant="secondary"
				onClick={deleteContainer}
			/>
			<Link to={`${container.id}/map`} class="contents">
				<Button variant="secondary" startIcon="map">
					{$t("map")}
				</Button>
			</Link>
			<Link to={`${container.id}/edit`} class="contents">
				<Button startIcon="edit">{$t("editInfo")}</Button>
			</Link>
		</DetailsHeader>
		<DetailsContent>
			<DetailsSection label={$t("generalInfo")}>
				<DetailsFields>
					<Field label={$t("location")} value={locationName} />
					<Field
						label={$t("containers.category")}
						value={$t(`containers.category.${container.category}`)}
					/>
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
		<div class="flex h-1/2 flex-col items-center justify-center">
			<h2 class="text-2xl font-semibold">{$t("containers.notFound.title")}</h2>
			<p>{$t("containers.notFound.description")}</p>
		</div>
	{/await}
</Card>
