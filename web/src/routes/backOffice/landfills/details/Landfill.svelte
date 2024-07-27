<script lang="ts">
	import { Link, navigate } from "svelte-routing";
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
	import type { Landfill } from "../../../../domain/landfill";

	/**
	 * Landfill ID.
	 */
	export let id: string;

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

		navigate(BackOfficeRoutes.LANDFILLS);
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
		<DetailsHeader to="" title={locationName}>
			<Button
				startIcon="delete"
				actionType="danger"
				variant="secondary"
				onClick={deleteLandfill}
			/>
			<Link to={`${landfill.id}/map`} class="contents">
				<Button variant="secondary" startIcon="map">
					{$t("map")}
				</Button>
			</Link>
			<Link to={`${landfill.id}/edit`} class="contents">
				<Button startIcon="edit">{$t("editInfo")}</Button>
			</Link>
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
