/**
 * Updates the current URL with new search params.
 * @param searchParams Search params to be set in the URL.
 */
export function updateSearchParams(searchParams: URLSearchParams) {
	history.pushState(
		null,
		"",
		`${location.origin}${location.pathname}?${searchParams}`,
	);
}
