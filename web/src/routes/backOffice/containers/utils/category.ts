import type { ContainerCategory } from "../../../../domain/container";

/**
 * Indicates if a given input is a valid container category.
 * @param input Input to be validated.
 */
export function isValidContainerCategory(
	input: string,
): input is ContainerCategory {
	switch (input) {
		case "general":
		case "paper":
		case "plastic":
		case "metal":
		case "glass":
		case "organic":
		case "hazardous":
			return true;
	}

	return false;
}
