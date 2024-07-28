import { get } from "svelte/store";

import { t } from "./i8n";

/**
 * Retrieves the location name of a resource given its way and municipality name.
 * @param wayName Way name.
 * @param municipalityName Municipality name.
 * @returns Location name.
 */
export function getLocationName(
	wayName: string | undefined,
	municipalityName: string | undefined,
) {
	const translation = get(t);
	const unknownWay = translation("location.unknownWay");

	const way = wayName ?? unknownWay;

	let locationName = way;

	if (municipalityName) {
		locationName += `, ${municipalityName}`;
	}

	return locationName;
}
