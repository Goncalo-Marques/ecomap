import type { EmployeeRoles } from "../../../../domain/employees";

/**
 * Indicates if given input is a valid employee role.
 * @param input
 */
export function isValidEmployeeRole(input: string): input is EmployeeRoles {
	switch (input) {
		case "wasteOperator":
		case "manager":
			return true;
	}

	return false;
}
