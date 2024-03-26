import { derived, writable } from "svelte/store";

/**
 * Href store.
 */
const href = writable(window.location.href);

const originalPushState = history.pushState;
const originalReplaceState = history.replaceState;

/**
 * Updates href store with the current location href.
 */
function updateHref() {
	href.set(window.location.href);
}

// Monkey patch pushState and replaceState events to update href store when either one is triggered.
history.pushState = function (...args) {
	originalPushState.apply(this, args);
	updateHref();
};
history.replaceState = function (...args) {
	originalReplaceState.apply(this, args);
	updateHref();
};

window.addEventListener("popstate", updateHref);
window.addEventListener("hashchange", updateHref);

/**
 * URL store.
 */
const url = {
	subscribe: derived(href, $href => new URL($href)).subscribe,
};

export default url;
