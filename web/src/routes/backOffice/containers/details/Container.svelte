<script lang="ts">
	import type { Container } from "../../../../domain/container";
	import Spinner from "../../../../lib/components/Spinner.svelte";
	import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
	import { t } from "../../../../lib/utils/i8n";
	import Field from "../../../../lib/components/Field.svelte";
	import DetailsFields from "../../../../lib/components/details/DetailsFields.svelte";
	import DetailsSection from "../../../../lib/components/details/DetailsSection.svelte";
	import DetailsContent from "../../../../lib/components/details/DetailsContent.svelte";
	import DetailsHeader from "../../../../lib/components/details/DetailsHeader.svelte";
	import { formatDate } from "../../../../lib/utils/date";
	import { DateFormats } from "../../../../lib/constants/date";
	import Input from "../../../../lib/components/Input.svelte";
	import Map from "../../../../lib/components/map/Map.svelte";
	import OlMap from "ol/Map";
	import Select from "../../../../lib/components/Select.svelte";
	import Option from "../../../../lib/components/Option.svelte";
	import { categoryOptions } from "../../../../lib/constants/container";
	import FormControl from "../../../../lib/components/FormControl.svelte";
	import Button from "../../../../lib/components/Button.svelte";
	import { Link } from "svelte-routing";
	import SelectLocation from "./SelectLocation.svelte";

	/**
	 * Container ID.
	 */
	export let id: string;

	/**
	 * TODO.
	 */
	export let mode: "view" | "edit";

	let map: OlMap;

	let openSelectLocation = false;

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

{#await containerPromise}
	<div class="container-loading">
		<Spinner />
	</div>
{:then container}
	<DetailsHeader to="" title="TODO">
		{#if mode === "view"}
			<Link to={`${container.id}/edit`}>
				<Button startIcon="edit">Editar informação</Button>
			</Link>
		{:else}
			<Button variant="tertiary">Cancelar</Button>
			<Button startIcon="check">Guardar</Button>
		{/if}
	</DetailsHeader>
	<DetailsContent>
		<DetailsSection label={$t("generalInfo")}>
			<DetailsFields>
				{#if mode === "view"}
					<Field
						label={$t("containers.category")}
						value={$t(`containers.category.${container.category}`)}
					/>
					<Field label={$t("containers.location")} value="TODO" />
				{:else}
					<FormControl label={$t("containers.category")}>
						<Select name="category" value={container.category}>
							{#each categoryOptions as category}
								<Option value={category}>
									{$t(`containers.category.${category}`)}
								</Option>
							{/each}
						</Select>
					</FormControl>
					<FormControl label={$t("containers.location")}>
						<Input
							readonly
							name="location"
							placeholder={$t("containers.location")}
							endIcon="location_on"
							onClick={() => (openSelectLocation = true)}
						/>
					</FormControl>
				{/if}
			</DetailsFields>
		</DetailsSection>
		{#if mode === "view"}
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
		{:else}
			<DetailsSection class="container-map-preview" label={"Pré-visualização"}>
				<Map bind:map />
			</DetailsSection>
		{/if}
		<SelectLocation
			open={openSelectLocation}
			onClose={() => (openSelectLocation = false)}
		/>
	</DetailsContent>
{:catch}
	<div class="container-not-found">
		<h2>{$t("containers.notFound.title")}</h2>
		<p>{$t("containers.notFound.description")}</p>
	</div>
{/await}

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
