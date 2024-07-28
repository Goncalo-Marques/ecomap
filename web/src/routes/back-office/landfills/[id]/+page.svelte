<script lang="ts">
	import { goto } from "$app/navigation";
	import { page } from "$app/stores";
	import type { Landfill } from "$domain/landfill";
	import ecomapHttpClient from "$lib/clients/ecomap/http";
	import Button from "$lib/components/Button.svelte";
	import DetailsContent from "$lib/components/details/DetailsContent.svelte";
	import DetailsFields from "$lib/components/details/DetailsFields.svelte";
	import DetailsHeader from "$lib/components/details/DetailsHeader.svelte";
	import DetailsSection from "$lib/components/details/DetailsSection.svelte";
	import Field from "$lib/components/Field.svelte";
	import Spinner from "$lib/components/Spinner.svelte";
	import { DateFormats } from "$lib/constants/date";
	import { BackOfficeRoutes } from "$lib/constants/routes";
	import { getToastContext } from "$lib/contexts/toast";
	import { formatDate } from "$lib/utils/date";
	import { t } from "$lib/utils/i8n";
	import { getLocationName } from "$lib/utils/location";

	import Card from "../../components/Card.svelte";

	/**
	 * Landfill ID.
	 */
	const id: string = $page.params.id;

	/**
	 * Toast context.
	 */
	const toast = getToastContext();

	/**
	 * Fetches landfill data.
	 */
	async function fetchLandfill(): Promise<Landfill> {
		const res = await ecomapHttpClient.GET("/landfills/{landfillId}", {
			params: { path: { landfillId: id } },
		});

		if (res.error) {
			throw new Error("Failed to retrieve landfill details");
		}

		return res.data;
	}

	/**
	 * Deletes the landfill displayed on the page.
	 */
	async function deleteLandfill() {
		const res = await ecomapHttpClient.DELETE("/landfills/{landfillId}", {
			params: {
				path: {
					landfillId: id,
				},
			},
		});

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
			title: $t("landfills.delete.success"),
			description: undefined,
		});

		goto(BackOfficeRoutes.LANDFILLS);
	}

	const landfillPromise = fetchLandfill();
</script>

<Card class="m-10 flex flex-col gap-10">
	{#await landfillPromise}
		<Spinner class="flex h-full items-center justify-center" />
	{:then landfill}
		{@const locationName = getLocationName(
			landfill.geoJson.properties.wayName,
			landfill.geoJson.properties.municipalityName,
		)}
		<DetailsHeader href={BackOfficeRoutes.LANDFILLS} title={locationName}>
			<Button
				startIcon="delete"
				actionType="danger"
				variant="secondary"
				onClick={deleteLandfill}
			/>
			<a
				href={`${BackOfficeRoutes.LANDFILLS}/${landfill.id}/map`}
				class="contents"
			>
				<Button variant="secondary" startIcon="map">
					{$t("map")}
				</Button>
			</a>
			<a
				href={`${BackOfficeRoutes.LANDFILLS}/${landfill.id}/edit`}
				class="contents"
			>
				<Button startIcon="edit">{$t("editInfo")}</Button>
			</a>
		</DetailsHeader>
		<DetailsContent>
			<DetailsSection label={$t("generalInfo")}>
				<DetailsFields>
					<Field label={$t("location")} value={locationName} />
				</DetailsFields>
			</DetailsSection>
			<DetailsSection label={$t("additionalInfo")}>
				<DetailsFields>
					<Field
						label={$t("createdAt")}
						value={formatDate(landfill.createdAt, DateFormats.shortDateTime)}
					/>
					<Field
						label={$t("modifiedAt")}
						value={formatDate(landfill.modifiedAt, DateFormats.shortDateTime)}
					/>
				</DetailsFields>
			</DetailsSection>
		</DetailsContent>
	{:catch}
		<div class="flex h-1/2 flex-col items-center justify-center">
			<h2 class="text-2xl font-semibold">{$t("landfills.notFound.title")}</h2>
			<p>{$t("landfills.notFound.description")}</p>
		</div>
	{/await}
</Card>
