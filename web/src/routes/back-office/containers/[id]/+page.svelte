<script lang="ts">
	import { goto } from "$app/navigation";
	import { page } from "$app/stores";
	import type { Container } from "$domain/container";
	import ecomapHttpClient from "$lib/clients/ecomap/http";
	import Button from "$lib/components/Button.svelte";
	import DetailsContent from "$lib/components/details/DetailsContent.svelte";
	import DetailsFields from "$lib/components/details/DetailsFields.svelte";
	import DetailsHeader from "$lib/components/details/DetailsHeader.svelte";
	import DetailsSection from "$lib/components/details/DetailsSection.svelte";
	import Field from "$lib/components/Field.svelte";
	import Spinner from "$lib/components/Spinner.svelte";
	import { DateFormats } from "$lib/constants/date";
	import { getToastContext } from "$lib/contexts/toast";
	import { formatDate } from "$lib/utils/date";
	import { t } from "$lib/utils/i8n";
	import { getLocationName } from "$lib/utils/location";

	import { BackOfficeRoutes } from "../../../constants/routes";
	import Card from "../../components/Card.svelte";

	/**
	 * Container ID.
	 */
	const id: string = $page.params.id;

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

		goto(BackOfficeRoutes.CONTAINERS);
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
		<DetailsHeader href={BackOfficeRoutes.CONTAINERS} title={locationName}>
			<Button
				startIcon="delete"
				actionType="danger"
				variant="secondary"
				onClick={deleteContainer}
			/>
			<a
				href={`${BackOfficeRoutes.CONTAINERS}/${container.id}/map`}
				class="contents"
			>
				<Button variant="secondary" startIcon="map">
					{$t("map")}
				</Button>
			</a>
			<a
				href={`${BackOfficeRoutes.CONTAINERS}/${container.id}/edit`}
				class="contents"
			>
				<Button startIcon="edit">{$t("editInfo")}</Button>
			</a>
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
