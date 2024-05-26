<script lang="ts">
	import type { RouteEmployee } from "../../../../../domain/routeEmployee";
	import Table from "../../../../../lib/components/table/Table.svelte";
	import type { Columns } from "../../../../../lib/components/table/types";
	import { DEFAULT_PAGE_SIZE } from "../../../../../lib/constants/pagination";
	import { t } from "../../../../../lib/utils/i8n";
	import { formatTime } from "../../../../../lib/utils/date";
	import type { ComponentProps } from "svelte";
	import TableDetailsAction from "../../../../../lib/components/table/TableDetailsAction.svelte";
	import createRouteEmployeesStore from "./routeEmployeesStore";
	import { ROUTE_EMPLOYEES_ROLES } from "../../constants/routeEmployee";
	import { BackOfficeRoutes } from "../../../../constants/routes";

	/**
	 * Route ID.
	 */
	export let routeId: string;

	const { loading, data, filters } = createRouteEmployeesStore(routeId);

	const columns: Columns<RouteEmployee> = [
		{
			type: "accessor",
			field: "firstName",
			header: $t("employees.firstName"),
			enableSorting: false,
			enableFiltering: false,
			cell(firstName) {
				return firstName;
			},
		},
		{
			type: "accessor",
			field: "lastName",
			header: $t("employees.lastName"),
			enableSorting: false,
			enableFiltering: false,
			cell(lastName) {
				return lastName;
			},
		},
		{
			type: "accessor",
			field: "routeRole",
			header: $t("routes.employees.role"),
			enableSorting: false,
			enableFiltering: true,
			filterOptions: ROUTE_EMPLOYEES_ROLES.map(routeRole => {
				return {
					value: routeRole,
					label: $t(`routes.employees.role.${routeRole}`),
				};
			}),
			filterInitialValue: $filters.routeRole,
			cell(routeRole) {
				return $t(`routes.employees.role.${routeRole}`);
			},
			onFilterChange(routeRole) {
				filters.update(filters => {
					return {
						...filters,
						pageIndex: 0,
						routeRole,
					};
				});
			},
		},
		{
			type: "accessor",
			field: "scheduleStart",
			header: $t("employees.schedule"),
			enableSorting: false,
			enableFiltering: false,
			cell(scheduleStart, row) {
				return `${formatTime(scheduleStart)} - ${formatTime(row.scheduleEnd)}`;
			},
		},
		{
			type: "display",
			header: "",
			align: "center",
			size: 120,
			cell(row) {
				const props: ComponentProps<TableDetailsAction> = {
					id: `${BackOfficeRoutes.EMPLOYEES}/${row.id}`,
				};

				return {
					component: TableDetailsAction,
					props,
				};
			},
		},
	];

	/**
	 * Handles changes of the route employees table pages.
	 * @param pageIndex New page index.
	 */
	function handlePageChange(pageIndex: number) {
		filters.update(store => {
			return {
				...store,
				pageIndex,
			};
		});
	}
</script>

<Table
	{columns}
	loading={$loading}
	rows={$data.employees}
	pagination={{
		name: $t("operators").toLowerCase(),
		pageIndex: $filters.pageIndex,
		pageSize: DEFAULT_PAGE_SIZE,
		total: $data.total,
		onPageChange: handlePageChange,
	}}
/>
