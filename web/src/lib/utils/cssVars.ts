const root = document.querySelector(":root");

if (!root) {
	throw new Error(":root element not found");
}

const rootStyle = getComputedStyle(root);

/**
 * Retrieves a CSS variable value.
 * @param variable CSS variable.
 * @returns CSS variable value.
 */
export function getCssVariable(variable: string) {
	return rootStyle.getPropertyValue(variable);
}
