import type { EmployeesFilters, PaginatedEmployees } from "$domain/employees";
import ecomapHttpClient from "$lib/clients/ecomap/http";
import { DEFAULT_PAGE_SIZE } from "$lib/constants/pagination";
import { BackOfficeRoutes } from "$lib/constants/routes";
import { createTableStore } from "$lib/stores/table";

/**
 * The search parameter names for each filter of the employees table.
 */
const FILTERS_PARAMS_NAMES: Record<keyof EmployeesFilters, string> = {
	pageIndex: "pageIndex",
	username: "username",
	sort: "sort",
	order: "order",
	role: "role",
};

/**
 * The initial data of the employees table.
 */
const initialData: PaginatedEmployees = {
	employees: [],
	total: 0,
};

/**
 * The initial filters of the employees table.
 */
export const initialFilters: EmployeesFilters = {
	pageIndex: 0,
	username: "",
	sort: "createdAt",
	order: "desc",
	role: undefined,
};

/**
 * Maps URL search params to employees filters.
 * @param searchParams URL search params.
 * @returns Employees filters.
 */
function searchParamsToFilters(
	searchParams: URLSearchParams,
): EmployeesFilters {
	let pageIndex = initialFilters.pageIndex;
	let username = initialFilters.username;
	let role = initialFilters.role;
	let sort = initialFilters.sort;
	let order = initialFilters.order;

	const pageIndexParam = Number(
		searchParams.get(FILTERS_PARAMS_NAMES.pageIndex),
	);
	const nameParam = searchParams.get(FILTERS_PARAMS_NAMES.username);
	const roleParam = searchParams.get(FILTERS_PARAMS_NAMES.role);
	const sortParam = searchParams.get(FILTERS_PARAMS_NAMES.sort);
	const orderParam = searchParams.get(FILTERS_PARAMS_NAMES.order);

	// Update page index when it's a valid number.
	if (!Number.isNaN(pageIndexParam)) {
		pageIndex = pageIndexParam;
	}

	// Update username when it's a non empty value.
	if (nameParam) {
		username = nameParam;
	}

	// Update role when it's a valid role.
	switch (roleParam) {
		case "wasteOperator":
		case "manager":
			role = roleParam;
	}

	// Update sort when sortParam is username.
	switch (sortParam) {
		case "username":
		case "firstName":
		case "lastName":
		case "createdAt":
			sort = sortParam;
	}

	// Update order when it's a valid sorting direction.
	switch (orderParam) {
		case "asc":
		case "desc":
			order = orderParam;
	}

	return {
		pageIndex,
		username,
		sort,
		role,
		order,
	};
}

/**
 * Maps filters of the employees table to URL search params.
 * @param filters employees filters.
 * @returns URL search params.
 */
function filtersToSearchParams(filters: EmployeesFilters): URLSearchParams {
	const { pageIndex, username, role, sort, order } = filters;

	const searchParams = new URLSearchParams(window.location.search);
	searchParams.set(FILTERS_PARAMS_NAMES.pageIndex, pageIndex.toString());

	if (username) {
		searchParams.set(FILTERS_PARAMS_NAMES.username, username);
	} else {
		searchParams.delete(FILTERS_PARAMS_NAMES.username);
	}

	if (role) {
		searchParams.set(FILTERS_PARAMS_NAMES.role, role);
	} else {
		searchParams.delete(FILTERS_PARAMS_NAMES.role);
	}

	if (sort) {
		searchParams.set(FILTERS_PARAMS_NAMES.sort, sort);
	} else {
		searchParams.delete(FILTERS_PARAMS_NAMES.sort);
	}

	if (order) {
		searchParams.set(FILTERS_PARAMS_NAMES.order, order);
	} else {
		searchParams.delete(FILTERS_PARAMS_NAMES.order);
	}

	return searchParams;
}

/**
 * Retrieves employees to be displayed in the employees table.
 * @param filters Employees filters.
 * @returns Employees.
 */
async function getEmployees(
	filters: EmployeesFilters,
): Promise<PaginatedEmployees> {
	const { pageIndex, username, role, sort, order } = filters;

	const res = await ecomapHttpClient.GET("/employees", {
		params: {
			query: {
				offset: pageIndex * DEFAULT_PAGE_SIZE,
				limit: DEFAULT_PAGE_SIZE,
				sort,
				order,
				role,
				username: username || undefined,
			},
		},
	});

	if (res.error) {
		return { total: 0, employees: [] };
	}

	return res.data;
}

const EmployeesStore = createTableStore(
	BackOfficeRoutes.EMPLOYEES,
	initialData,
	filtersToSearchParams,
	searchParamsToFilters,
	getEmployees,
);

export default EmployeesStore;
