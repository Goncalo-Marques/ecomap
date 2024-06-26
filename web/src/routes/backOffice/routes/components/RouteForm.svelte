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
	import LocationInput from "../../../../lib/components/LocationInput.svelte";
	import Select from "../../../../lib/components/Select.svelte";
	import ecomapHttpClient from "../../../../lib/clients/ecomap/http";
	import { getBatchPaginatedResponse } from "../../../../lib/utils/request";
	import type { Truck } from "../../../../domain/truck";
	import type { Warehouse } from "../../../../domain/warehouse";
	import Option from "../../../../lib/components/Option.svelte";
	import Input from "../../../../lib/components/Input.svelte";
	import SelectContainers from "./SelectContainers.svelte";
	import OperatorsTable from "./OperatorsTable.svelte";
	import type { Employee } from "../../../../domain/employees";
	import type {
		Container,
		SelectedRouteContainersIds,
	} from "../../../../domain/container";
	import type {
		RouteEmployee,
		SelectedRouteEmployees,
	} from "../../../../domain/routeEmployee";
	import { getToastContext } from "../../../../lib/contexts/toast";

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
		truck: Truck,
		containersIds: SelectedRouteContainersIds,
		routeEmployees: SelectedRouteEmployees,
	) => void;

	/**
	 * Route data.
	 * @default null
	 */
	export let route: Route | null = null;

	/**
	 * Indicates if form is being submitted.
	 */
	export let isSubmitting: boolean;

	/**
	 * Toast context.
	 */
	const toast = getToastContext();

	/**
	 * The select containers open modal state.
	 * @default false
	 */
	let openSelectContainers = false;

	/**
	 * The containers map with the original, added and deleted containers for the route.
	 */
	let containersMap: {
		original: Container[];
		added: Container[];
		deleted: Container[];
	} = {
		original: [],
		added: [],
		deleted: [],
	};

	/**
	 * The route employees.
	 */
	let routeEmployees: RouteEmployee[] = [];

	/**
	 * The selected drivers for the route.
	 */
	let selectedDrivers: Employee[] = [];

	/**
	 * The selected collectors for the route.
	 */
	let selectedCollectors: Employee[] = [];

	/**
	 * The promise with the truck options for the truck select input.
	 */
	let truckOptionsPromise: Promise<Truck[]>;

	/**
	 * The promise with the warehouse options for the truck select input.
	 */
	let warehouseOptionsPromise: Promise<Warehouse[]>;

	/**
	 * The loading state for the waste operators.
	 * @default false
	 */
	let loadingWasteOperators: boolean = false;

	/**
	 * The waste operators.
	 */
	let wasteOperators: Employee[];

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
		truckInput: HTMLSelectElement,
		departureWarehouseInput: HTMLSelectElement,
		arrivalWarehouseInput: HTMLSelectElement,
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
	async function handleSubmit(
		e: Event & { currentTarget: EventTarget & HTMLFormElement },
	) {
		const form = e.currentTarget;
		const formElements = form.elements;
		const nameInput = formElements.namedItem("name");
		const truckInput = formElements.namedItem("truck");
		const departureWarehouseInput =
			formElements.namedItem("departureWarehouse");
		const arrivalWarehouseInput = formElements.namedItem("arrivalWarehouse");
		const containersInput = formElements.namedItem("containers");

		if (
			!(nameInput instanceof HTMLInputElement) ||
			!(truckInput instanceof HTMLSelectElement) ||
			!(departureWarehouseInput instanceof HTMLSelectElement) ||
			!(arrivalWarehouseInput instanceof HTMLSelectElement) ||
			!(containersInput instanceof HTMLInputElement)
		) {
			throw new Error("Form elements are not inputs");
		}

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

		let selectedTruck: Truck;
		try {
			selectedTruck = JSON.parse(truckInput.value);

			if (
				selectedTruck.personCapacity <
				selectedDrivers.length + selectedCollectors.length
			) {
				toast.show({
					type: "error",
					title: $t("routes.truck.error.personCapacityExceeded.title"),
					description: $t(
						"routes.truck.error.personCapacityExceeded.description",
					),
				});
				return;
			}
		} catch {
			return;
		}

		// Set the added and deleted container IDs.
		const containersIds = {
			added: containersMap.added.map(container => container.id),
			deleted: containersMap.deleted.map(container => container.id),
		};

		const selectedRouteEmployees: SelectedRouteEmployees = {
			added: [],
			deleted: [],
		};

		// Add all new drivers.
		for (const selectedDriver of selectedDrivers) {
			if (routeEmployees.every(operator => operator.id !== selectedDriver.id)) {
				selectedRouteEmployees.added.push({
					id: selectedDriver.id,
					routeRole: "driver",
				});
			}
		}

		// Add all removed drivers.
		for (const routeEmployee of routeEmployees.filter(
			routeOperator => routeOperator.routeRole === "driver",
		)) {
			if (
				selectedDrivers.every(
					selectedDriver => selectedDriver.id !== routeEmployee.id,
				)
			) {
				selectedRouteEmployees.deleted.push({
					id: routeEmployee.id,
					routeRole: "driver",
				});
			}
		}

		// Add all new collectors.
		for (const selectedCollector of selectedCollectors) {
			if (
				routeEmployees.every(
					routeOperator => routeOperator.id !== selectedCollector.id,
				)
			) {
				selectedRouteEmployees.added.push({
					id: selectedCollector.id,
					routeRole: "collector",
				});
			}
		}

		// Add all removed collectors.
		for (const routeOperator of routeEmployees.filter(
			routeOperator => routeOperator.routeRole === "collector",
		)) {
			if (
				selectedCollectors.every(
					selectedCollector => selectedCollector.id !== routeOperator.id,
				)
			) {
				selectedRouteEmployees.deleted.push({
					id: routeOperator.id,
					routeRole: "collector",
				});
			}
		}

		onSave(
			nameInput.value,
			departureWarehouseInput.value,
			arrivalWarehouseInput.value,
			selectedTruck,
			containersIds,
			selectedRouteEmployees,
		);
	}

	/**
	 * Retrieves the value displayed in the container input.
	 * @param originalAmount The number of containers on the route.
	 * @param addedAmount The number of containers added to the route.
	 * @param deletedAmount The number of containers deleted from the route.
	 * @returns Container input value.
	 */
	function getContainersInputValue(
		originalAmount: number,
		addedAmount: number,
		deletedAmount: number,
	) {
		const containerAmount = originalAmount + addedAmount - deletedAmount;

		return `${containerAmount} ${$t(containerAmount === 1 ? "container" : "containers").toLowerCase()}`;
	}

	/**
	 * Retrieves the route operators from the route ID.
	 * @param id Route ID.
	 */
	async function getRouteOperators(id: string) {
		routeEmployees = await getBatchPaginatedResponse(async (limit, offset) => {
			const res = await ecomapHttpClient.GET("/routes/{routeId}/employees", {
				params: { path: { routeId: id }, query: { limit, offset } },
			});

			if (res.error) {
				return { total: 0, items: [] };
			}

			return { total: res.data.total, items: res.data.employees };
		});

		selectedDrivers = routeEmployees.filter(
			routeOperator => routeOperator.routeRole === "driver",
		);
		selectedCollectors = routeEmployees.filter(
			routeOperator => routeOperator.routeRole === "collector",
		);
	}

	/**
	 * Retrieves the route containers from the route ID.
	 * @param id Route ID.
	 */
	async function getRouteContainers(id: string) {
		const routeContainers = await getBatchPaginatedResponse(
			async (limit, offset) => {
				const res = await ecomapHttpClient.GET("/routes/{routeId}/containers", {
					params: { path: { routeId: id }, query: { limit, offset } },
				});

				if (res.error) {
					return { total: 0, items: [] };
				}

				return { total: res.data.total, items: res.data.containers };
			},
		);

		containersMap.original = routeContainers;
	}

	/**
	 * Retrieves truck options.
	 * @returns Truck options.
	 */
	async function getTruckOptions() {
		return getBatchPaginatedResponse(async (limit, offset) => {
			const res = await ecomapHttpClient.GET("/trucks", {
				params: { query: { limit, offset } },
			});

			if (res.error) {
				return { total: 0, items: [] };
			}

			return { total: res.data.total, items: res.data.trucks };
		});
	}

	/**
	 * Retrieves warehouse options.
	 * @returns Warehouse options.
	 */
	async function getWarehouseOptions() {
		return getBatchPaginatedResponse(async (limit, offset) => {
			const res = await ecomapHttpClient.GET("/warehouses", {
				params: { query: { limit, offset } },
			});

			if (res.error) {
				return { total: 0, items: [] };
			}

			return { total: res.data.total, items: res.data.warehouses };
		});
	}

	/**
	 * Fetch waste operators.
	 */
	async function fetchWasteOperators() {
		loadingWasteOperators = true;

		wasteOperators = await getBatchPaginatedResponse(async (limit, offset) => {
			const res = await ecomapHttpClient.GET("/employees", {
				params: {
					query: {
						limit,
						offset,
						role: "wasteOperator",
					},
				},
			});

			if (res.error) {
				return { total: 0, items: [] };
			}

			return { total: res.data.total, items: res.data.employees };
		});

		loadingWasteOperators = false;
	}

	truckOptionsPromise = getTruckOptions();
	warehouseOptionsPromise = getWarehouseOptions();
	fetchWasteOperators();

	// If a route is defined, fetch its containers and operators.
	if (route) {
		getRouteContainers(route.id);
		getRouteOperators(route.id);
	}
</script>

<form novalidate on:submit|preventDefault={handleSubmit}>
	<DetailsHeader to={back} {title}>
		<Link to={back} style="display:contents">
			<Button variant="tertiary">{$t("cancel")}</Button>
		</Link>
		<Button type="submit" disabled={isSubmitting} startIcon="check">
			{$t("save")}
		</Button>
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
						value={route?.name}
						error={!!formErrorMessages.name}
						placeholder={$t("routes.name.placeholder")}
					/>
				</FormControl>

				<FormControl
					label={$t("truck")}
					error={!!formErrorMessages.truck}
					helperText={formErrorMessages.truck}
				>
					{#await truckOptionsPromise}
						<Select placeholder={$t("routes.truck.placeholder")} />
					{:then truckOptions}
						<Select
							required
							name="truck"
							error={!!formErrorMessages.truck}
							placeholder={$t("routes.truck.placeholder")}
							value={route ? JSON.stringify(route.truck) : undefined}
						>
							{#each truckOptions as truck}
								<Option value={JSON.stringify(truck)}>
									{`${truck.make} ${truck.model} (${truck.licensePlate}, ${$t("personCapacity").toLowerCase()}: ${truck.personCapacity})`}
								</Option>
							{/each}
						</Select>
					{/await}
				</FormControl>

				<FormControl
					label={$t("departure")}
					error={!!formErrorMessages.departureWarehouse}
					helperText={formErrorMessages.departureWarehouse}
				>
					{#await warehouseOptionsPromise}
						<Select placeholder={$t("routes.departure.placeholder")} />
					{:then warehouseOptions}
						<Select
							required
							name="departureWarehouse"
							error={!!formErrorMessages.departureWarehouse}
							placeholder={$t("routes.departure.placeholder")}
							value={route?.departureWarehouse.id}
						>
							{#each warehouseOptions as warehouse}
								<Option value={warehouse.id}>
									{getLocationName(
										warehouse.geoJson.properties.wayName,
										warehouse.geoJson.properties.municipalityName,
									)}
								</Option>
							{/each}
						</Select>
					{/await}
				</FormControl>

				<FormControl
					label={$t("arrival")}
					error={!!formErrorMessages.arrivalWarehouse}
					helperText={formErrorMessages.arrivalWarehouse}
				>
					{#await warehouseOptionsPromise}
						<Select placeholder={$t("routes.arrival.placeholder")} />
					{:then warehouseOptions}
						<Select
							required
							name="arrivalWarehouse"
							error={!!formErrorMessages.arrivalWarehouse}
							placeholder={$t("routes.arrival.placeholder")}
							value={route?.arrivalWarehouse.id}
						>
							{#each warehouseOptions as warehouse}
								<Option value={warehouse.id}>
									{getLocationName(
										warehouse.geoJson.properties.wayName,
										warehouse.geoJson.properties.municipalityName,
									)}
								</Option>
							{/each}
						</Select>
					{/await}
				</FormControl>

				<FormControl
					label={$t("containers")}
					error={!!formErrorMessages.containers}
					helperText={formErrorMessages.containers}
				>
					<LocationInput
						required
						value={getContainersInputValue(
							containersMap.original.length,
							containersMap.added.length,
							containersMap.deleted.length,
						)}
						name="containers"
						placeholder={$t("routes.containers.placeholder")}
						error={!!formErrorMessages.containers}
						onClick={() => (openSelectContainers = true)}
					/>
				</FormControl>
			</DetailsFields>
		</DetailsSection>
		<DetailsSection
			class="drivers-collectors"
			label={$t("routes.employees.role.drivers")}
		>
			<OperatorsTable
				operators={wasteOperators}
				loading={loadingWasteOperators}
				bind:selectedOperators={selectedDrivers}
				disabledOperators={selectedCollectors}
			/>
		</DetailsSection>
		<DetailsSection
			class="drivers-collectors"
			label={$t("routes.employees.role.collectors")}
		>
			<OperatorsTable
				operators={wasteOperators}
				loading={loadingWasteOperators}
				bind:selectedOperators={selectedCollectors}
				disabledOperators={selectedDrivers}
			/>
		</DetailsSection>
	</DetailsContent>
	<SelectContainers
		routeId={route?.id}
		open={openSelectContainers}
		onOpenChange={open => (openSelectContainers = open)}
		onSave={(addedContainers, deletedContainers) => {
			containersMap.added = addedContainers;
			containersMap.deleted = deletedContainers;
		}}
	/>
</form>

<style>
	form {
		display: contents;
	}

	:global(.drivers-collectors) {
		flex: 1;
	}
</style>
