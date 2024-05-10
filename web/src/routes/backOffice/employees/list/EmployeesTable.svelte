<script lang="ts">
	import type {
		Employee,
		EmployeeSortableFields,
	} from "../../../../domain/employees";
	import Table from "../../../../lib/components/table/Table.svelte";
	import type {
		Columns,
		SortingDirection,
	} from "../../../../lib/components/table/types";
	import { DEFAULT_PAGE_SIZE } from "../../../../lib/constants/pagination";
	import { t } from "../../../../lib/utils/i8n";
	import { rolesOptions } from "../constants/roles";
	import employeesStore from "./employeesStore";
	import { formatDate } from "../../../../lib/utils/date";
	import { DateFormats } from "../../../../lib/constants/date";

	const { loading, data, filters } = employeesStore;

	const columns: Columns<Employee> = [
		{
			type: "accessor",
			field: "username",
			header: $t("employees.username"),
			enableSorting: true,
			enableFiltering: false,
			cell(username) {
				return username;
			},
		},
		{
			type: "accessor",
			field: "firstName",
			header: $t("employees.firstName"),
			enableSorting: true,
			enableFiltering: false,
			cell(firstName) {
				return firstName;
			},
		},
		{
			type: "accessor",
			field: "lastName",
			header: $t("employees.lastName"),
			enableSorting: true,
			enableFiltering: false,
			cell(lastName) {
				return lastName;
			},
		},
		{
			type: "accessor",
			field: "role",
			header: $t("employees.role"),
			enableSorting: false,
			enableFiltering: true,
			filterOptions: rolesOptions.map(role => {
				return {
					value: role,
					label: $t(`employees.role.${role}`),
				};
			}),
			filterInitialValue: $filters.role,
			cell(role) {
				return $t(`employees.role.${role}`);
			},
			onFilterChange(role) {
				filters.update(filters => {
					return {
						...filters,
						pageIndex: 0,
						role,
					};
				});
			},
		},
		{
			type: "accessor",
			field: "scheduleEnd",
			header: $t("employees.schedule"),
			enableSorting: false,
			enableFiltering: false,
			cell(scheduleStart, row) {
				// The date 1970-01-01, date when Unix time started, is used here to format scheduleStart and scheduleEnd. This is needed because the formatDate() helper uses Intl.DateTimeFormat to format dates, without this date the helper will not work.
				return `${formatDate(`1970-01-01 ${scheduleStart}`, DateFormats.shortTime)} - ${formatDate(`1970-01-01 ${row.scheduleEnd}`, DateFormats.shortTime)}`;
			},
		},
	];

	/**
	 * Handles changes of the employees table pages.
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

	/**
	 * Handles changes of employees username sorting order.
	 * @param field Field to be sorted.
	 * @param direction Direction in which the field is sorted.
	 */
	function handleSortingChange(
		field: EmployeeSortableFields,
		direction: SortingDirection,
	) {
		filters.update(store => {
			return {
				...store,
				order: direction,
				sort: field,
			};
		});
	}
</script>

<Table
	{columns}
	loading={$loading}
	rows={$data.employees}
	sortingField={$filters.sort}
	sortingOrder={$filters.order}
	onSortingChange={handleSortingChange}
	pagination={{
		name: $t("employees.title").toLowerCase(),
		pageIndex: $filters.pageIndex,
		pageSize: DEFAULT_PAGE_SIZE,
		total: $data.total,
		onPageChange: handlePageChange,
	}}
/>
