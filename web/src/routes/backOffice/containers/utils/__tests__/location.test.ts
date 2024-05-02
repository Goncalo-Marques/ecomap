import { describe, expect, it } from "vitest";
import { getContainerLocation } from "../location";

describe("Container utils", () => {
	const containerLocationCases = [
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

	it.each(containerLocationCases)(
		"should return $expectedLocation",
		({ wayName, municipalityName, expectedLocation }) => {
			expect(getContainerLocation(wayName, municipalityName)).toBe(
				expectedLocation,
			);
		},
	);
});
