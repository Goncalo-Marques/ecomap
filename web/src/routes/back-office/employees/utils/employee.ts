import type { EmployeeRole } from "$domain/employees";

/**
 * Indicates if given input is a valid employee role.
 * @param input Input to be validated.
 * @returns True if input is a valid role, false otherwise.
 */
export function isValidEmployeeRole(input: string): input is EmployeeRole {
	switch (input) {
		case "wasteOperator":
		case "manager":
			return true;
	}

	return false;
}
