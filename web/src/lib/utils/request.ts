interface PaginatedResponse<T> {
	total: number;
	items: T[];
}

/**
 * Retrieves a batch of all paginated items for a given request.
 * @param request Single request to be made.
 * @param [limit=100] Pagination limit.
 * @returns Paginated items.
 */
export async function getBatchPaginatedResponse<T>(
	request: (limit: number, offset: number) => Promise<PaginatedResponse<T>>,
	limit: number = 100,
): Promise<T[]> {
	const items: T[] = [];

	const firstResponse = await request(limit, 0);
	if (!firstResponse.total) {
		return items;
	}

	items.push(...firstResponse.items);

	const remainingRequestsAmount = Math.ceil(firstResponse.total / limit) - 1;
	const remainingRequestsPromises: Promise<PaginatedResponse<T>>[] = [];

	for (let i = 1; i <= remainingRequestsAmount; i++) {
		remainingRequestsPromises.push(request(limit, limit * i));
	}

	const responses = await Promise.allSettled(remainingRequestsPromises);

	for (const response of responses) {
		if (response.status === "fulfilled") {
			items.push(...response.value.items);
		}
	}

	return items;
}
