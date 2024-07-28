import { describe, expect, it } from "vitest";

import { getVisiblePages } from "../utils";

describe("Table utils", () => {
	describe("Pagination", () => {
		const pages = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9];

		const cases = [
			{
				pageIndex: 0,
				expectedVisiblePages: [0, 1, 2, 3, 4],
			},
			{
				pageIndex: 3,
				expectedVisiblePages: [1, 2, 3, 4, 5],
			},
			{
				pageIndex: 9,
				expectedVisiblePages: [5, 6, 7, 8, 9],
			},
		];

		it.each(cases)(
			"should return visible pages",
			({ pageIndex, expectedVisiblePages }) => {
				expect(getVisiblePages(pages, pageIndex)).toEqual(expectedVisiblePages);
			},
		);
	});
});
