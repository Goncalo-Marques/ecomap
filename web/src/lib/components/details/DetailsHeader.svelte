<script lang="ts">
	import Icon from "../Icon.svelte";

	/**
	 * The title displayed in the header.
	 */
	export let title: string;

	/**
	 * The route to navigate to when the back button is clicked.
	 */
	export let href: string;
</script>

<header class="flex flex-wrap items-center gap-4">
	<a {href} class="contents text-green-700">
		<Icon name="arrow_back" size="medium" />
	</a>
	<h1 class="flex-1 text-2xl font-semibold">{title}</h1>
	<div class="flex items-center gap-4">
		<slot />
	</div>
</header>
