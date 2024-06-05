import type { components } from "../../api/ecomap/http";
import type { SortingDirection } from "../lib/components/table/types";

/**
 * Employee.
 */
export type Employee = components["schemas"]["Employee"];

/**
 * Sortable fields of a employees.
 */
export type EmployeeSortableFields = NonNullable<
	components["parameters"]["EmployeeSortQueryParam"]
>;

/**
 * Employee roles.
 */
export type EmployeeRole = Employee["role"];

/**
 * Paginated employees.
 */
export type PaginatedEmployees = components["schemas"]["EmployeesPaginated"];

/**
 * Filters of employees.
 */
export interface EmployeesFilters {
	pageIndex: number;
	username: string;
	sort: EmployeeSortableFields;
	order: SortingDirection;
	role?: EmployeeRole;
}
