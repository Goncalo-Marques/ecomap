import type { Coordinate } from "ol/coordinate";
import { transform } from "ol/proj";

import { OL_PROJECTION, RESOURCE_PROJECTION } from "$lib/constants/map";

/**
 * Converts coordinate in {@link OL_PROJECTION} projection to the resource geometry projection.
 * @param coordinate Coordinate in {@link OL_PROJECTION} projection.
 * @returns Coordinate in resource projection system.
 */
export function convertToResourceProjection(coordinate: Coordinate) {
	return transform(coordinate, OL_PROJECTION, RESOURCE_PROJECTION);
}

/**
 * Converts coordinate in {@link RESOURCE_PROJECTION} projection to the map projection.
 * @param coordinate Coordinate in {@link RESOURCE_PROJECTION} projection.
 * @returns Coordinate in map projection system.
 */
export function convertToMapProjection(coordinate: Coordinate) {
	return transform(coordinate, RESOURCE_PROJECTION, OL_PROJECTION);
}
