import { pushState } from "$app/navigation";

/**
 * Updates the current URL with new search params.
 * @param searchParams Search params to be set in the URL.
 */
export function updateSearchParams(searchParams: URLSearchParams) {
	pushState(`${location.pathname}?${searchParams}`, {});
}

/**
 * Clears search params from the current URL.
 */
export function clearSearchParams() {
	pushState(location.pathname, {});
}
