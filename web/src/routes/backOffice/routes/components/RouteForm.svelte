<script lang="ts">
	import Button from "../../../../lib/components/Button.svelte";
	import { t } from "../../../../lib/utils/i8n";
	import DetailsFields from "../../../../lib/components/details/DetailsFields.svelte";
	import DetailsSection from "../../../../lib/components/details/DetailsSection.svelte";
	import DetailsContent from "../../../../lib/components/details/DetailsContent.svelte";
	import FormControl from "../../../../lib/components/FormControl.svelte";
	import { Link } from "svelte-routing";
	import DetailsHeader from "../../../../lib/components/details/DetailsHeader.svelte";
	import { getLocationName } from "../../../../lib/utils/location";
	import type { Route } from "../../../../domain/route";
	import { LANDFILL_ICON_SRC } from "../../../../lib/constants/map";
	import LocationInput from "../../../../lib/components/LocationInput.svelte";
	import Select from "../../../../lib/components/Select.svelte";
	import { onMount } from "svelte";
	import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
	import { getBatchPaginatedResponse } from "../../../../lib/utils/request";
	import type { Truck } from "../../../../domain/truck";
	import type { Warehouse } from "../../../../domain/warehouse";
	import Option from "../../../../lib/components/Option.svelte";
	import Input from "../../../../lib/components/Input.svelte";
	import SelectContainers from "./SelectContainers.svelte";
	import DriversTable from "./DriversTable.svelte";
	import CollectorsTable from "./CollectorsTable.svelte";

	/**
	 * The back route.
	 */
	export let back: string;

	/**
	 * The title in the form.
	 */
	export let title: string;

	/**
	 * Callback fired when save action is triggered.
	 */
	export let onSave: (
		name: string,
		departureWarehouseId: string,
		arrivalWarehouseId: string,
		truckId: string,
		containerIds: string[],
	) => void;

	/**
	 * Route data.
	 * @default null
	 */
	export let route: Route | null = null;

	/**
	 * The select containers open modal state.
	 * @default false
	 */
	let openSelectContainers = false;

	/**
	 * TODO.
	 */
	let containers = {
		added: [],
		deleted: [],
	};

	const options: { truck: Truck[]; warehouse: Warehouse[] } = {
		truck: [],
		warehouse: [],
	};

	/**
	 * Error messages of the form fields.
	 */
	let formErrorMessages = {
		name: "",
		departureWarehouse: "",
		arrivalWarehouse: "",
		truck: "",
		containers: "",
	};

	/**
	 * Validates the form and sets error messages on the form fields
	 * if they contain any errors.
	 * @param nameInput Name input field.
	 * @param truckInput Truck input field.
	 * @param departureWarehouseInput Departure warehouse input field.
	 * @param arrivalWarehouseInput Arrival warehouse input field.
	 * @param containersInput Containers input field.
	 * @returns `true` if form is valid, `false` otherwise.
	 */
	function validateForm(
		nameInput: HTMLInputElement,
		truckInput: HTMLInputElement,
		departureWarehouseInput: HTMLInputElement,
		arrivalWarehouseInput: HTMLInputElement,
		containersInput: HTMLInputElement,
	) {
		if (nameInput.validity.valueMissing) {
			formErrorMessages.name = $t("error.valueMissing");
		} else {
			formErrorMessages.name = "";
		}

		if (truckInput.validity.valueMissing) {
			formErrorMessages.truck = $t("error.valueMissing");
		} else {
			formErrorMessages.truck = "";
		}

		if (departureWarehouseInput.validity.valueMissing) {
			formErrorMessages.departureWarehouse = $t("error.valueMissing");
		} else {
			formErrorMessages.departureWarehouse = "";
		}

		if (arrivalWarehouseInput.validity.valueMissing) {
			formErrorMessages.arrivalWarehouse = $t("error.valueMissing");
		} else {
			formErrorMessages.arrivalWarehouse = "";
		}

		if (containersInput.validity.valueMissing) {
			formErrorMessages.containers = $t("error.valueMissing");
		} else {
			formErrorMessages.containers = "";
		}

		return (
			!formErrorMessages.name &&
			!formErrorMessages.truck &&
			!formErrorMessages.departureWarehouse &&
			!formErrorMessages.arrivalWarehouse &&
			!formErrorMessages.containers
		);
	}

	/**
	 * Handles the submit event of the form.
	 * @param e Submit event.
	 */
	function handleSubmit(
		e: Event & { currentTarget: EventTarget & HTMLFormElement },
	) {
		const form = e.currentTarget;
		const formElements = form.elements;
		const nameInput = formElements.namedItem("name") as HTMLInputElement;
		const truckInput = formElements.namedItem("truck") as HTMLInputElement;
		const departureWarehouseInput = formElements.namedItem(
			"departureWarehouse",
		) as HTMLInputElement;
		const arrivalWarehouseInput = formElements.namedItem(
			"arrivalWarehouse",
		) as HTMLInputElement;
		const containersInput = formElements.namedItem(
			"containers",
		) as HTMLInputElement;

		// Check if form is valid to prevent making a server request.
		if (
			!validateForm(
				nameInput,
				truckInput,
				departureWarehouseInput,
				arrivalWarehouseInput,
				containersInput,
			)
		) {
			return;
		}
	}

	/**
	 * TODO.
	 */
	async function getOptions() {
		const [trucksRes, warehousesRes] = await Promise.allSettled([
			getBatchPaginatedResponse(async (limit, offset) => {
				const res = await ecomapHttpClient.GET("/trucks", {
					params: { query: { limit, offset } },
				});

				if (res.error) {
					return { total: 0, items: [] };
				}

				return { total: res.data.total, items: res.data.trucks };
			}),
			getBatchPaginatedResponse(async (limit, offset) => {
				const res = await ecomapHttpClient.GET("/warehouses", {
					params: { query: { limit, offset } },
				});

				if (res.error) {
					return { total: 0, items: [] };
				}

				return { total: res.data.total, items: res.data.warehouses };
			}),
		]);

		if (trucksRes.status === "fulfilled") {
			options.truck = trucksRes.value;
		}

		if (warehousesRes.status === "fulfilled") {
			options.warehouse = warehousesRes.value;
		}
	}

	onMount(() => {
		getOptions();
	});
</script>

<form novalidate on:submit|preventDefault={handleSubmit}>
	<DetailsHeader to={back} {title}>
		<Link to={back} style="display:contents">
			<Button variant="tertiary">{$t("cancel")}</Button>
		</Link>
		<Button type="submit" startIcon="check">{$t("save")}</Button>
	</DetailsHeader>
	<DetailsContent>
		<DetailsSection label={$t("generalInfo")}>
			<DetailsFields>
				<FormControl
					label={$t("routes.name")}
					error={!!formErrorMessages.name}
					helperText={formErrorMessages.name}
				>
					<Input
						required
						name="name"
						maxLength={50}
						error={!!formErrorMessages.name}
						placeholder={$t("routes.name.placeholder")}
					/>
				</FormControl>

				<FormControl
					label={$t("truck")}
					error={!!formErrorMessages.truck}
					helperText={formErrorMessages.truck}
				>
					<Select
						required
						name="truck"
						error={!!formErrorMessages.truck}
						placeholder={$t("routes.truck.placeholder")}
						value={route?.truck.id}
					>
						{#each options.truck as truck}
							<Option value={truck.id}>
								{`${truck.make} ${truck.model} (${truck.licensePlate})`}
							</Option>
						{/each}
					</Select>
				</FormControl>

				<FormControl
					label={$t("departure")}
					error={!!formErrorMessages.departureWarehouse}
					helperText={formErrorMessages.departureWarehouse}
				>
					<Select
						required
						name="truck"
						error={!!formErrorMessages.departureWarehouse}
						placeholder={$t("routes.departure.placeholder")}
						value={route?.departureWarehouse.id}
					>
						{#each options.warehouse as warehouse}
							<Option value={warehouse.id}>
								{getLocationName(
									warehouse.geoJson.properties.wayName,
									warehouse.geoJson.properties.municipalityName,
								)}
							</Option>
						{/each}
					</Select>
				</FormControl>

				<FormControl
					label={$t("arrival")}
					error={!!formErrorMessages.arrivalWarehouse}
					helperText={formErrorMessages.arrivalWarehouse}
				>
					<Select
						required
						name="truck"
						error={!!formErrorMessages.arrivalWarehouse}
						placeholder={$t("routes.arrival.placeholder")}
						value={route?.arrivalWarehouse.id}
					>
						{#each options.warehouse as warehouse}
							<Option value={warehouse.id}>
								{getLocationName(
									warehouse.geoJson.properties.wayName,
									warehouse.geoJson.properties.municipalityName,
								)}
							</Option>
						{/each}
					</Select>
				</FormControl>

				<FormControl
					label={$t("containers")}
					error={!!formErrorMessages.containers}
					helperText={formErrorMessages.containers}
				>
					<LocationInput
						required
						name="containers"
						placeholder={$t("routes.containers.placeholder")}
						error={!!formErrorMessages.containers}
						onClick={() => (openSelectContainers = true)}
					/>
				</FormControl>
			</DetailsFields>
		</DetailsSection>
		<DetailsSection label={$t("routes.employees.role.drivers")}>
			<DriversTable />
		</DetailsSection>
		<DetailsSection label={$t("routes.employees.role.collectors")}>
			<CollectorsTable />
		</DetailsSection>
	</DetailsContent>
	<SelectContainers
		open={openSelectContainers}
		onOpenChange={open => (openSelectContainers = open)}
		onSave={() => {}}
	/>
</form>

<style>
	form {
		display: contents;
	}
</style>
