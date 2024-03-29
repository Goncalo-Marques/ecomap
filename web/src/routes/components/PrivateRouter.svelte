<script lang="ts">
	import { onMount } from "svelte";
	import { writable } from "svelte/store";
	import { Router } from "svelte-routing";
	import { CommonRoutes } from "../constants/routes";
	import { isAuthenticated } from "../../lib/utils/auth";
	import url from "../../lib/utils/url";

	let isUserAuthenticated = writable(false);

	/**
	 * Verifies if user is authenticated to access a private route.
	 * @returns `true` if user is authenticated. Otherwise, returns `false`.
	 */
	function verifyAuthentication(): boolean {
		// Redirect to sign in page if user is not authenticated.
		const authenticated = isAuthenticated();

		isUserAuthenticated.set(authenticated);

		return authenticated;
	}

	/**
	 * Callback fired when URL changes.
	 * Verifies if the user is authenticated and if not, redirects to sign in page.
	 */
	function onUrlChange() {
		if (!verifyAuthentication()) {
			location.replace(CommonRoutes.SIGN_IN);
		}
	}

	onMount(() => {
		// Verifies if the user is authenticated to access a route inside the private router.
		verifyAuthentication();

		// URL subscription to watch for route changes.
		const unsubscribe = url.subscribe(onUrlChange);

		return function () {
			unsubscribe();
		};
	});
</script>

<Router>
	{#if $isUserAuthenticated}
		<slot />
	{/if}
</Router>
