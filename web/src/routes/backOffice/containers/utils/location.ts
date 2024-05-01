import { get } from "svelte/store";
import { t } from "../../../../lib/utils/i8n";

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
	const unknownMunicipality = translation(
		"containers.location.unknownMunicipality",
	);

	const way = wayName ?? unknownWay;

	let municipality = unknownMunicipality;
	if (municipalityName) {
		municipality = municipalityName;
	}

	return `${way}, ${municipality}`;
}
