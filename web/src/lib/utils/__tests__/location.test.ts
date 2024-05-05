import { describe, expect, it } from "vitest";
import { getLocationName } from "../location";

describe("Utils", () => {
	const locationCases = [
		{
			wayName: "Main Street",
			municipalityName: "London",
			expectedLocation: "Main Street, London",
		},
		{
			wayName: undefined,
			municipalityName: "London",
			expectedLocation: "Unknown way, London",
		},
		{
			wayName: "Main Street",
			municipalityName: undefined,
			expectedLocation: "Main Street",
		},
		{
			wayName: undefined,
			municipalityName: undefined,
			expectedLocation: "Unknown way",
		},
	];

	it.each(locationCases)(
		"should return correct location name",
		({ wayName, municipalityName, expectedLocation }) => {
			expect(getLocationName(wayName, municipalityName)).toBe(expectedLocation);
		},
	);
});
