<script lang="ts">
	import type { MouseEventHandler } from "svelte/elements";

	import type { IconSize } from "../../domain/components/icon";
	import Icon from "./Icon.svelte";

	/**
	 * The name of the icon to be used.
	 */
	export let icon: string;

	/**
	 * The icon size.
	 */
	export let size: IconSize = "medium";

	/**
	 * Callback fired when button is clicked.
	 * @default null
	 */
	export let onClick: MouseEventHandler<HTMLButtonElement> | null = null;
</script>

<button
	type="button"
	class="flex items-center justify-center rounded-full p-1 hover:bg-gray-100 active:bg-gray-200"
	on:click={onClick}
>
	<Icon name={icon} {size} />
</button>
