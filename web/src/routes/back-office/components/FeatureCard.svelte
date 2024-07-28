<script lang="ts">
	import Button from "$lib/components/Button.svelte";
	import Icon from "$lib/components/Icon.svelte";

	/**
	 * The icon to be displayed.
	 */
	export let icon: string;

	/**
	 * The title of the card.
	 */
	export let title: string;

	/**
	 * The resource href.
	 */
	export let resourceHref: string;
</script>

<article class="flex flex-col gap-2 rounded border border-gray-300 p-4">
	<div class="flex items-center gap-2">
		<div class="flex flex-1 items-center gap-1">
			<Icon name={icon} />
			<h3 class="flex-1">{title}</h3>
		</div>
		<a href={resourceHref} class="contents">
			<Button startIcon="open_in_new" variant="tertiary" />
		</a>
	</div>
	<slot />
</article>
