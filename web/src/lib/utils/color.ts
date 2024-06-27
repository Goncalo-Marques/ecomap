/**
 * Retrieves a color with an opacity.
 * @param hex Color hexadecimal value.
 * @param opacity Value between 0 and 1.
 * @returns Hexadecimal color with an opacity.
 */
export function getColorWithOpacity(hex: string, opacity: number) {
	return `${hex}${Math.floor(opacity * 255).toString(16)}`;
}
