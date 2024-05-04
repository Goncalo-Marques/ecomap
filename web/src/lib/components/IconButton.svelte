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

<button type="button" on:click={onClick}>
	<Icon name={icon} {size} />
</button>

<style>
	button {
		display: flex;
		justify-content: center;
		align-items: center;
		border-radius: 50%;
		padding: 0.25rem;

		&:hover {
			background-color: var(--gray-100);
		}
		&:active {
			background-color: var(--gray-200);
		}
	}
</style>
