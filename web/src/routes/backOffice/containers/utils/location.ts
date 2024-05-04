import { get } from "svelte/store";
import { t } from "../../../../lib/utils/i8n";
import type { Coordinate } from "ol/coordinate";

/**
 * Retrieves the location of a container given its way and municipality name.
 * @param wayName Container way name.
 * @param municipalityName Container municipality name.
 * @returns Container location.
 */
export function getContainerLocation(
	wayName: string | undefined,
	municipalityName: string | undefined,
) {
	const translation = get(t);
	const unknownWay = translation("containers.location.unknownWay");

	const way = wayName ?? unknownWay;

	let locationName = way;

	if (municipalityName) {
		locationName += `, ${municipalityName}`;
	}

	return locationName;
}

/**
 * Formats the container location coordinate.
 * @param coordinate Container coordinate.
 */
export function formatContainerCoordinate(coordinate: Coordinate | undefined) {
	if (!coordinate) {
		return;
	}

	const [lon, lat] = coordinate;

	return `${lat}, ${lon}`;
}
